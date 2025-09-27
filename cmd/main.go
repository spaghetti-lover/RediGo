package main

import (
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
	go server.RunIoMultiplexingServer(&wg)
	go server.WaitForSignal(&wg, sigChan)
	wg.Wait()
}
