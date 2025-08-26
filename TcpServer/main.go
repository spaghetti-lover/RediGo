package main

import (
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	// read data from client
	log.Println(conn.RemoteAddr())
	var buf []byte = make([]byte, 1000)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	// process
	time.Sleep(time.Second * 10)
	// reply
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, world\r\n"))
	conn.Close()
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		// conn == socket == dedicated communication channel
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// create a new goroutine to handle the connection
		go handleConnection(conn)
	}
}
