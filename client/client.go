package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/TheBarn/tinyIRC/utils"
)

const (
	intro = `Welcome to this tiny IRC client

Commands:
------------------------------
/nick <nickname>
/list
/join <#channel>
/leave <#channel>
/who
/msg <nickname> <message>
------------------------------

First Enter your nickname using the /nick command:
`
)

func getServerMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		fmt.Printf("> ")
	}
}

func launchPrompt(conn net.Conn) {
	fmt.Printf(intro)
	fmt.Printf("> ")
	scanner := bufio.NewScanner(os.Stdin)
	go getServerMessages(conn)
	for scanner.Scan() {
		fmt.Printf("> ")
		err := utils.SendBytes(conn, scanner.Text())
		if err != nil {
			fmt.Println("Server is down.", err)
			os.Exit(1)
		}
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
