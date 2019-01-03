package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const WORKERS = 10

func main() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8888", time.Microsecond*1000)
	if err != nil {
		log.Fatal("DialTimeout Error", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(WORKERS)
	
	for i := 0; i < WORKERS; i++ {
		go doWork(conn, &wg, i)

	}
	wg.Wait()
}

func doWork(conn net.Conn, wg* sync.WaitGroup, threadId int) {
	defer wg.Done()

	start := time.Now()
	buffer := []byte("hi there" + "\n")
	n, err := conn.Write(buffer)
	if err != nil {
		log.Printf("Error writing buffer: %v", err)
		return
	}
	fmt.Printf("thread %d wrote %d bytes\n", n, threadId)
	// read response
	response := make([]byte, 8)
	n, err = conn.Read(response)
	if err != nil {
		log.Printf("Error while reading response: %v\n", err)
	}
	elapsed := time.Since(start)
	fmt.Printf("%s Thread %d received: %s\n", elapsed, threadId, response)
}
