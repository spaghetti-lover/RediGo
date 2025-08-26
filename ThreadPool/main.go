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

// represent the thread in the pool
type Worker struct {
	id      int
	jobChan chan Job
}

// represent the thread pool
type Pool struct {
	jobQueue chan Job
	workers  []*Worker
}

// create a new worker
func NewWorker(id int, jobChan chan Job) *Worker {
	return &Worker{
		id:      id,
		jobChan: jobChan,
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobChan {
			log.Printf("Worker %d is handling job from %s", w.id, job.conn.RemoteAddr())
			handleConnection(job.conn)
		}
	}()
}

func NewPool(numOfWorker int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers:  make([]*Worker, numOfWorker),
	}
}

// push job to queue
func (p *Pool) AddJob(conn net.Conn) {
	p.jobQueue <- Job{conn: conn}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue)
		p.workers[i] = worker
		worker.Start()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1000)
	conn.Read(buf)
	time.Sleep(1 * time.Second)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nEngineer Pro\r\n"))
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// 1 pool with 2 threads
	pool := NewPool(2)
	pool.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//go handleConnection(conn)
		pool.AddJob(conn)
	}
}
