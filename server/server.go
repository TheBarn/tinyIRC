package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/TheBarn/tinyIRC/utils"
)

type user struct {
	conn net.Conn
	nick string
}

func handleChanMessage(user *user, msg string) {
	fmt.Println("Message:", msg)
}

func checkNickname(nick string) bool {
	if len(nick) > 9 {
		return false
	}
	_, err := regexp.MatchString(`^[a-zA-Z0-9_]$`, nick)
	if err != nil {
		return false
	}
	return true
}

func handleCommand(user *user, cmd string) {
	if cmd == "" {
		return
	}
	if cmd[0] != '/' {
		handleChanMessage(user, cmd)
		return
	}
	args := strings.Fields(cmd)
	switch args[0] {
	case "/nick":
		if len(args) != 2 {
			utils.SendBytes(user.conn, "command /nick takes one nickname as argument")
			return
		}
		nick := args[1]
		ok := checkNickname(nick)
		if !ok {
			utils.SendBytes(user.conn, "nickname should have only 9 characters in [a-zA-Z0-9_]")
			return
		}
		user.nick = nick
		utils.SendBytes(user.conn, "/nick "+nick)
		utils.SendBytes(user.conn, "your nickame was changed to "+nick)
	}
}

func readCommand(user *user) {
	fmt.Println("READ", user)
	scanner := bufio.NewScanner(user.conn)
	for {
		if ok := scanner.Scan(); !ok {
			break
		}
		go handleCommand(user, string(scanner.Text()))
	}
}

func handleRequest(conn net.Conn) {
	fmt.Printf("Serving %v\n", conn.RemoteAddr())
	user := user{conn: conn}
	readCommand(&user)
	fmt.Println("close connection", user)
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
