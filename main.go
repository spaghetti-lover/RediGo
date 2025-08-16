package main

import (
	"log"
	"net"
	"time"
)

// element in the queue
type Job struct {
	conn net.Conn
}

// thread in the pool
type Worker struct {
	id       int
	jobQueue chan Job
}

type Pool struct {
	// queue
	jobQueue chan Job
	workers  []*Worker
}

func NewPool(n int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers:  make([]*Worker, n),
	}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue)
		p.workers[i] = worker
		worker.Start()
	}
}

func (p *Pool) AddJob(conn net.Conn) {
	p.jobQueue <- Job{
		conn: conn,
	}
}

func NewWorker(id int, jobQueue chan Job) *Worker {
	return &Worker{
		id:       id,
		jobQueue: jobQueue,
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobQueue {
			log.Printf("worker %d is processing from %s", w.id, job.conn.RemoteAddr())
			handleConnection(job.conn)
		}
	}()
}

// receive req and reply response
func handleConnection(conn net.Conn) {
	defer conn.Close()
	var buf []byte = make([]byte, 4096)

	for {
		// Blocking call
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

		// Pretend to receive request
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

	defer listener.Close()

	pool := NewPool(2)
	pool.Start()
	for {
		// conn == socket == communication channel
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		pool.AddJob(conn)
	}
}
