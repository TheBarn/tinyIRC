package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/TheBarn/tinyIRC/utils"
)

const (
	intro = `Welcome to this tiny IRC client

Commands:
------------------------------
/nick <nickname>
/list
/join <#channel>
/leave
/who
/msg <nickname> <message>
------------------------------

First Enter your nickname using the /nick command:
`
)

type user struct {
	nick    string
	channel string
}

func printPrompt(user *user) {
	fmt.Printf("\r\033[K")
	if user.channel != "" {
		fmt.Printf("%s ", user.channel)
	}
	if user.nick != "" {
		fmt.Printf("%s ", user.nick)
	}
	fmt.Printf("> ")
}

func printMsg(user *user, msg string) {
	fmt.Printf("\r\033[K\n\033[2A")
	fmt.Printf(msg)
	fmt.Printf("\n\n")
	printPrompt(user)
}

func handleServerMessage(msg string, user *user) {
	args := strings.Fields(msg)
	switch args[0] {
	case "/msg":
		message := msg[5:]
		printMsg(user, message)
	case "/nick":
		if len(args) == 2 {
			user.nick = args[1]
			printPrompt(user)
		}
	case "/join":
		if len(args) == 2 {
			user.channel = args[1]
			printPrompt(user)
		}
	case "/leave":
		user.channel = ""
		printPrompt(user)
	case "/warning":
		message := "\033[0;31m" + msg[9:] + "\033[0m"
		printMsg(user, message)
	}
}

func getServerMessages(conn net.Conn, user *user) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		handleServerMessage(scanner.Text(), user)
	}
}

func pingServer(conn net.Conn, user *user) {
	for {
		time.Sleep(time.Second)
		err := utils.SendBytes(conn, "")
		if err != nil {
			printMsg(user, "Server is down")
			os.Exit(1)
		}
	}
}

func handleInput(user *user, conn net.Conn, input string) error {
	//	if input != "" && input[0] != '/' {
	//		printMsg(user, input)
	//	}
	err := utils.SendBytes(conn, input)
	return err
}

func launchPrompt(conn net.Conn) {
	user := user{}
	go getServerMessages(conn, &user)
	go pingServer(conn, &user)
	fmt.Printf(intro)
	printPrompt(&user)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		fmt.Printf("\r\033[K\033[1A\033[K\n")
		err := handleInput(&user, conn, input)
		if err != nil {
			printMsg(&user, "Server is down")
			os.Exit(1)
		}
		printPrompt(&user)
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
