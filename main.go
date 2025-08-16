package main

import (
	"log"
	"net"
	"time"
)

// receive req and reply response
func handleConnection(conn net.Conn) {
	defer conn.Close()
	var buf []byte = make([]byte, 4096)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			log.Fatal(err)
			return
		}

		req := string(buf[:n])
		log.Println("Received request: ", req)
		time.Sleep(1 * time.Millisecond)
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 12\r\n\r\nHello, world"))
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening at port 3000")

	for {
		// conn == socket == communication channel
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("handle conn from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
