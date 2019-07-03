package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/TheBarn/tinyIRC/utils"
)

func launchPrompt(conn net.Conn) {
	fmt.Println("Welcome to the tiny IRC client")
	fmt.Printf("------------------------------\n\n")
	fmt.Printf("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Printf("> ")
		conn.Write([]byte(scanner.Text() + "\n"))
	}
}

func main() {
	port := utils.ChoosePort()
	conn, err := net.Dial("tcp4", "localhost:"+port)
	if err != nil {
		fmt.Println("Error dialing:", err)
		os.Exit(1)
	}
	defer conn.Close()
	launchPrompt(conn)
}
