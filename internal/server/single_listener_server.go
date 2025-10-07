package server

import (
	"log"
	"net"
	"sync"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
)

func (s *Server) StartSingleListener(wg *sync.WaitGroup) {
	defer wg.Done()
	// Start all I/O handler event loops
	for _, handler := range s.ioHandlers {
		go handler.Run()
	}

	// Set up listener socket
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to acccept connection: %v", err)
			continue
		}

		// forward the new connection to an I/O handler in a round-robin manner
		handler := s.ioHandlers[s.nextIOHandler%s.numIOHandlers]
		s.nextIOHandler++

		if err := handler.AddConn(conn); err != nil {
			log.Printf("Failed to add connection to I/O handler %d: %v", handler.id, err)
			conn.Close()
		}
	}
}
