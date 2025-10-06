package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spaghetti-lover/multithread-redis/internal/server"
)

func main() {
	var wg sync.WaitGroup

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(2)

	//Run single threaded server with epoll/kqueue
	//go server.RunIoMultiplexingServer(&wg)

	//Run multi-threaded server with epoll/kqueue
	s := server.NewServer()
	go s.Start(&wg)
	//go s.StartMultiListeners(&wg)

	go server.WaitForSignal(&wg, sigChan, s)
	wg.Wait()

	log.Println("Graceful shutdown complete")
}
