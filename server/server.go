package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/TheBarn/tinyIRC/utils"
)

func handleRequest(conn net.Conn) {
	fmt.Printf("Serving %v\n", conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)
	for {
		if ok := scanner.Scan(); !ok {
			break
		}
		fmt.Println(scanner.Text())
	}
	fmt.Println("scanning ended")
	conn.Close()
}

func main() {
	port := utils.ChoosePort()
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
