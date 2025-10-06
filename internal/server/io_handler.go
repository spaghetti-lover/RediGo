package server

import (
	"io"
	"log"
	"sync"
	"syscall"

	"github.com/spaghetti-lover/multithread-redis/internal/core"
	"github.com/spaghetti-lover/multithread-redis/internal/core/iomux"
)

// import (
// 	"sync"

// 	"github.com/spaghetti-lover/multithread-redis/internal/core/iomux"
// )

type IOHandler struct {
	id            int
	ioMultiplexer iomux.IOMultiplexer
	mu            sync.Mutex
	server        *Server
}

func NewIOHandler(id int, server *Server) (*IOHandler, error) {
	multiplexer, err := iomux.CreateIOMultiplexer()
	if err != nil {
		return nil, err
	}

	return &IOHandler{
		id:            id,
		ioMultiplexer: multiplexer,
		server:        server,
	}, nil
}

// Add connection to the handler's epoll monitoring list
func (h *IOHandler) AddConn(connFd int) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	log.Printf("I/O Handler %d is monitoring fd %d", h.id, connFd)
	return h.ioMultiplexer.Monitor(iomux.Event{
		Fd: connFd,
		Op: iomux.OpRead,
	})
}

func (h *IOHandler) Run() {
	log.Printf("I/O Handler %d started", h.id)
	for {
		events, err := h.ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		for _, event := range events {
			connFd := event.Fd
			cmd, err := readCommand(connFd)
			if err != nil {
				if err == io.EOF || err == syscall.ECONNRESET {
					log.Println("client disconnected")
					_ = syscall.Close(connFd)
					continue
				}
				log.Println("read error:", err)
				continue
			}

			replyCh := make(chan []byte, 1)
			task := &core.Task{
				Command: cmd,
				ReplyCh: replyCh,
			}
			h.server.dispatch(task)
			res := <-replyCh
			syscall.Write(connFd, res)
		}
	}
}

func (h *IOHandler) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.ioMultiplexer.Close(); err != nil {
		log.Printf("Error closing I/O handler %d: %v", h.id, err)
	} else {
		log.Printf("I/O handler %d closed successfully", h.id)
	}
}
