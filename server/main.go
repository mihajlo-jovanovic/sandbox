package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", "127.0.0.1:8888")
	req := make(chan string)
	signal := make(chan bool)
	for {
		conn, err := ln.Accept()
		//conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Fatal("Error binding to port", err)
		}
		go readRequests(conn, req, signal)
		go doWork(conn, req, signal)
	}
}

func readRequests(conn net.Conn, out chan string, s chan bool) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	done := false
	for {
		pay, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("finished...")
				done = true
				break
			}
			log.Println(err)
		}
		fmt.Printf("Got: %s\n", pay)
		out <- pay
	}
	if done {
		s <- true
	}
}

func doWork(conn net.Conn, in chan string, s chan bool) {
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	for {
		select {
		case signal, ok := <-s:
			if !ok {
				log.Println("returning...")
				return
			}
			if signal {
				return
			}
		case line, ok := <-in:
			if !ok {
				log.Println("returning...")
				return
			}
			fmt.Printf("Read from channel: %s\n", line)
			//time.Sleep(1 * time.Second)
			_, err := writer.WriteString("hi there")
			if err == nil {
				writer.Flush()
			}
		}
	}
}
