package main

import "github.com/spaghetti-lover/multithread-redis/internal/server"

func main() {
	server.RunIoMultiplexingServer()
}
