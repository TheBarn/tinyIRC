package main

import (
	"fmt"
	"net"
	"os"
)

const (
	defaultPort = "5555"
	buffSize    = 1024
)

func handleRequest(conn net.Conn) {
	fmt.Printf("Connection: %v\n", conn)
	buf := make([]byte, buffSize)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
	}
	fmt.Println(string(buf))
	conn.Write([]byte("Message received!"))
	conn.Close()
}

func main() {
	port := defaultPort
	args := os.Args
	switch len(args) {
	case 1:
		fmt.Printf("Using default port %v\n", defaultPort)
	case 2:
		port = args[1]
		fmt.Printf("Using port %v\n", port)
	case 3:
		fmt.Printf("Usage: ./tinyIRC [port]\n")
	}
	l, err := net.Listen("tcp4", "localhost:"+port)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}
