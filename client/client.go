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
	if user.channel != "" {
		fmt.Printf("%s ", user.channel)
	}
	if user.nick != "" {
		fmt.Printf("%s ", user.nick)
	}
	fmt.Printf("> ")
}

func printServerMsg(msg string) {
	fmt.Println("\n\033[0;31m" + msg + "\033[0m")
}

func handleServerMessage(msg string, user *user) {
	args := strings.Fields(msg)
	switch args[0] {
	case "/nick":
		if len(args) == 2 {
			user.nick = args[1]
		}
	case "/join":
		if len(args) == 2 {
			user.channel = args[1]
		}
	case "/leave":
		user.channel = ""
	default:
		printServerMsg(msg)
		printPrompt(user)
	}
}

func getServerMessages(conn net.Conn, user *user) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		handleServerMessage(scanner.Text(), user)
	}
}

func pingServer(conn net.Conn) {
	for {
		time.Sleep(time.Second)
		err := utils.SendBytes(conn, "")
		if err != nil {
			printServerMsg("Server is down")
			os.Exit(1)
		}
	}
}

func launchPrompt(conn net.Conn) {
	user := user{}
	fmt.Printf(intro)
	printPrompt(&user)
	scanner := bufio.NewScanner(os.Stdin)
	go getServerMessages(conn, &user)
	go pingServer(conn)
	for scanner.Scan() {
		printPrompt(&user)
		err := utils.SendBytes(conn, scanner.Text())
		if err != nil {
			printServerMsg("Server is down")
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
