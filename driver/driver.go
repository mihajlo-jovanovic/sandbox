package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const WORKERS = 10

type logData struct {
	elapsed  int64
	request  string
	response string
}

type logWriter struct {
	w io.Writer
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Fprint(writer.w, strconv.Itoa(int(time.Now().UnixNano()/1000000))+"|"+string(bytes))
}

func startLogger(logStream chan logData) {
	filename := "log.out"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	requestLogger := log.New(os.Stdout, "", log.LstdFlags)
	requestLogger.SetFlags(0)
	requestLogger.SetOutput(&logWriter{file})
	for {
		select {
		case ld, ok := <-logStream:
			if !ok {
				log.Println("no remaining log data for", filename)
				return
			}
			requestLogger.Printf("%v|%v|%v\n", ld.elapsed, ld.request, ld.response)
		}
	}
}

func main() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8888", time.Microsecond*1000)
	if err != nil {
		log.Fatal("DialTimeout Error", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(WORKERS)
	lines := make(chan string, 100)
	for i := 0; i < 100; i++ {
		lines <- "hi there"
	}

	logStream := make(chan logData)
	go startLogger(logStream)

	for i := 0; i < WORKERS; i++ {
		go doWork(conn, &wg, i, lines, logStream)

	}
	wg.Wait()
}

func doWork(conn net.Conn, wg *sync.WaitGroup, threadId int, lines chan string, logStream chan logData) {
	defer wg.Done()
	for {
		select {
		case payload, ok := <-lines:
			if !ok {
				log.Println("returning...")
				return
			}
			start := time.Now()
			buffer := []byte(payload + "\n")
			_, err := conn.Write(buffer)
			if err != nil {
				log.Printf("Error writing buffer: %v", err)
				return
			}
			//fmt.Printf("thread %d wrote %d bytes\n", n, threadId)
			// read response
			response := make([]byte, 8)
			_, err = conn.Read(response)
			if err != nil {
				log.Printf("Error while reading response: %v\n", err)
			}
			elapsed := time.Since(start)
			fmt.Printf("%s Thread %d received: %s\n", elapsed, threadId, response)
			logStream <- logData{elapsed.Nanoseconds() / 1000000, payload, string(response)}
		}
	}
}
