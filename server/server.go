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
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
		}
		connInput := string(netData)
		fmt.Println(connInput)
		if connInput == "STOP" {
			break
		}
	}
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
