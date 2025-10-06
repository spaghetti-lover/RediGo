package server

import (
	"context"
	"log"
	"net"
	"sync"
	"syscall"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
	"golang.org/x/sys/unix"
)

func createReusablePortListener(network, addr string) (net.Listener, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var err error
			c.Control(func(fd uintptr) {
				err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			})
			return err
		}}
	return lc.Listen(context.Background(), network, addr)
}

func (s *Server) StartMultiListeners(wg *sync.WaitGroup) {
	defer wg.Done()

	//Start all I/O handler event loops
	for _, handler := range s.ioHandlers {
		go handler.Run()
	}

	for i := 0; i < config.ListenerNum; i++ {
		go func() {
			listener, err := createReusablePortListener(config.Protocol, config.Port)
			if err != nil {
				log.Fatal(err)
			}
			defer listener.Close()

			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("Failed to accept connections: %v", err)
					continue
				}

				tcpConn, ok := conn.(*net.TCPConn)
				if !ok {
					log.Println("Accepted connection is not a TCP connection")
					conn.Close()
					continue
				}

				connFile, err := tcpConn.File()
				if err != nil {
					log.Printf("Failed to get file desctiptor from TCP connection: %v", err)
					conn.Close()
					continue
				}
				connFd := int(connFile.Fd())

				// forward the new connection to an I/O handler in a round-robin manner
				handler := s.ioHandlers[s.nextIOHandler%s.numIOHandlers]
				s.nextIOHandler++

				if err := handler.AddConn(connFd); err != nil {
					log.Printf("Failed to add connection fd %d to I/O handler %d: %v", connFd, handler.id, err)
					_ = syscall.Close(connFd)

				}
			}
		}()
	}
}
