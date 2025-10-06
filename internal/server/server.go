package server

import (
	"context"
	"hash/fnv"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
	"github.com/spaghetti-lover/multithread-redis/internal/constant"
	"github.com/spaghetti-lover/multithread-redis/internal/core"
	"github.com/spaghetti-lover/multithread-redis/internal/core/iomux"
)

type Server struct {
	workers       []*core.Worker
	ioHandlers    []*IOHandler
	numWorkers    int
	numIOHandlers int

	// add listener to close it on shutdown
	listener net.Listener

	// For round-robin assigment of new connection to I/O handlers
	nextIOHandler int
}

func (s *Server) getPartitionID(key string) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	return int(hasher.Sum32()) % s.numWorkers
}

func (s *Server) dispatch(task *core.Task) {
	// Commands like PING etc., don't have a key.
	// We can send them to any worker.
	var key string
	if len(task.Command.Args) > 0 {
		key = task.Command.Args[0]
	}
	workerID := s.getPartitionID(key)
	s.workers[workerID].TaskCh <- task
}

func NewServer() *Server {
	numCores := runtime.NumCPU()
	numIOHandlers := numCores / 3
	numWorkers := numCores - numIOHandlers
	log.Printf("Initializing server with %d workers and %d io handler\n", numWorkers, numIOHandlers)

	s := &Server{
		workers:       make([]*core.Worker, numWorkers),
		ioHandlers:    make([]*IOHandler, numIOHandlers),
		numWorkers:    numWorkers,
		numIOHandlers: numIOHandlers,
	}

	for i := 0; i < numWorkers; i++ {
		s.workers[i] = core.NewWorker(i, 1024)
		s.workers[i].Start(context.Background())
	}

	for i := 0; i < numIOHandlers; i++ {
		handler, err := NewIOHandler(i, s)
		if err != nil {
			log.Fatalf("Failed to create I/O handler %d: %v", i, err)
		}
		s.ioHandlers[i] = handler
	}

	return s
}

func (s *Server) Stop() {
	log.Println("Stopping server and all workers...")

	atomic.StoreInt32(&serverStatus, constant.ServerStatusShutdown)

	// Close listener to stop Accept() accepting new connections
	if s.listener != nil {
		s.listener.Close()
	}

	for i, worker := range s.workers {
		if worker != nil {
			log.Printf("Stopping worker %d", i)
			worker.Stop()
		}
	}

	for i, handler := range s.ioHandlers {
		if handler != nil {
			log.Printf("Stopping I/O handler %d", i)
			handler.Stop()
		}
	}
}

// func response(fd int, respCh <-chan []byte) {
// 	res := <-respCh
// 	syscall.Write(fd, res)
// }

func (s *Server) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// Start I/O Handler event loops
	for _, handler := range s.ioHandlers {
		go handler.Run()
	}

	// Setup listener socket
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}

	s.listener = listener
	defer listener.Close()

	log.Printf("Server is listenning on %s", config.Port)

	for {
		if atomic.LoadInt32(&serverStatus) == constant.ServerStatusShutdown {
			return
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			return
		}

		// get the file descriptor of the connection
		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			log.Println("Accepted connection is not a TCP connection")
			conn.Close()
			continue
		}
		connFile, err := tcpConn.File()
		if err != nil {
			log.Printf("Failed to get file from TCP connection: %v", err)
			conn.Close()
			continue
		}
		defer connFile.Close()
		connFd := int(connFile.Fd())

		// forward the new connection to an I/O handler in a round-robin manner
		handler := s.ioHandlers[s.nextIOHandler%s.numIOHandlers]
		s.nextIOHandler++

		if err := handler.AddConn(connFd); err != nil {
			log.Printf("Failed to add connection fd %d to I/O handler %d: %v", connFd, handler.id, err)
			_ = syscall.Close(connFd)
		}
	}
}

var serverStatus int32 = constant.ServerStatusIdle

func readCommand(fd int) (*core.Command, error) {
	var buf = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}
	return core.ParseCmd(buf)
}

// func respond(data string, fd int) error {
// 	if _, err := syscall.Write(fd, []byte(data)); err != nil {
// 		return err
// 	}
// 	return nil
// }

func RunIoMultiplexingServer(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("starting an I/O Multiplexing TCP server on", config.Port)
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Get the file descriptor from the listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("listener is not a TCPListener")
	}
	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())

	// Create an ioMultiplexer instance (epoll in Linux, kqueue in MacOS)
	ioMultiplexer, err := iomux.CreateIOMultiplexer()
	if err != nil {
		log.Fatal(err)
	}
	defer ioMultiplexer.Close()

	// Monitor "read" events on the Server FD
	if err = ioMultiplexer.Monitor(iomux.Event{
		Fd: serverFd,
		Op: iomux.OpRead,
	}); err != nil {
		log.Fatal(err)
	}

	var events = make([]iomux.Event, config.MaxConnection)
	var lastActiveExpireExecTime = time.Now()

	for atomic.LoadInt32(&serverStatus) != constant.ServerStatusShutdown {
		// check last execution time and call if it is more than 100ms ago.
		if time.Now().After(lastActiveExpireExecTime.Add(constant.ActiveExpireFrequency)) {
			if !atomic.CompareAndSwapInt32(&serverStatus, constant.ServerStatusIdle, constant.ServerStatusRunning) {
				if atomic.LoadInt32(&serverStatus) == constant.ServerStatusShutdown {
					return
				}
			}
			core.ActiveDeleteExpiredKeys() //Busy
			atomic.SwapInt32(&serverStatus, constant.ServerStatusIdle)
			lastActiveExpireExecTime = time.Now() // Idle
		}
		// wait for file descriptors in the monitoring list to be ready for I/O
		// it is a blocking call.
		events, err = ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		if !atomic.CompareAndSwapInt32(&serverStatus, constant.ServerStatusIdle, constant.ServerStatusRunning) {
			if atomic.LoadInt32(&serverStatus) == constant.ServerStatusShutdown {
				return
			}
		}

		//Busy
		for i := 0; i < len(events); i++ {
			if events[i].Fd == serverFd {
				log.Printf("new client is trying to connect")
				// set up new connection
				connFd, _, err := syscall.Accept(serverFd)
				if err != nil {
					log.Println("err", err)
					continue
				}
				log.Printf("set up a new connection")
				// ask epoll to monitor this connection
				if err = ioMultiplexer.Monitor(iomux.Event{
					Fd: connFd,
					Op: iomux.OpRead,
				}); err != nil {
					log.Fatal(err)
				}
			} else {
				if atomic.LoadInt32(&serverStatus) == constant.ServerStatusShutdown {
					return
				}
				// handle data from an existing connection
				// read the command, execute it and write back the response.
				cmd, err := readCommand(events[i].Fd)
				if err != nil {
					if err == io.EOF || err == syscall.ECONNRESET {
						log.Println("client disconnected: ", err)

						err = ioMultiplexer.Unmonitor(iomux.Event{
							Fd: events[i].Fd,
							Op: iomux.OpRead,
						})
						if err != nil {
							log.Println("Can not unmonitor: ", err)
						}
						_ = syscall.Close(events[i].Fd)
						continue
					}
					log.Println("read error:", err)
					continue
				}
				if err = core.ExecuteAndResponse(cmd, events[i].Fd); err != nil {
					log.Println("err write: ", err)

					err = ioMultiplexer.Unmonitor(iomux.Event{
						Fd: events[i].Fd,
						Op: iomux.OpRead,
					})
					if err != nil {
						log.Println("Can not unmonitor: ", err)
					}
					_ = syscall.Close(events[i].Fd)
				}
			}
		}

		//Idle
		atomic.SwapInt32(&serverStatus, constant.ServerStatusIdle)
	}
}

func WaitForSignal(wg *sync.WaitGroup, sigChan chan os.Signal, server *Server) {
	defer wg.Done()
	<-sigChan

	log.Println("Signal received, initiating graceful shutdown...")

	// Stop server and all workers
	if server != nil {
		server.Stop()
	}

	// Set server status to shutdown
	atomic.StoreInt32(&serverStatus, constant.ServerStatusShutdown)
}
