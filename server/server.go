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

type server struct {
	channels []*channel
}

type channel struct {
	name  string
	users []*user
}

type user struct {
	conn    net.Conn
	nick    string
	channel *channel
}

func newServer() server {
	server := server{}
	channelNames := []string{"#welcome", "#IRChelp", "#golang"}
	for _, channelName := range channelNames {
		server.channels = append(server.channels, &channel{name: channelName})
	}
	return server
}

func handleChanMessage(user *user, msg string) {
	fmt.Println("Message:", msg)
}

func checkNickname(nick string) bool {
	//TOCHECK no null string!!
	if len(nick) > 9 {
		return false
	}
	_, err := regexp.MatchString(`^[a-zA-Z0-9_]$`, nick)
	if err != nil {
		return false
	}
	return true
}

func pickChannel(server *server, channelName string) (*channel, bool) {
	for _, channel := range server.channels {
		if channel.name == channelName {
			return channel, true
		}
	}
	return nil, false
}

func handleCommand(server *server, user *user, cmd string) {
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
			utils.SendBytes(user.conn, "usage: /nick nickname")
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
	case "/join":
		if len(args) != 2 {
			utils.SendBytes(user.conn, "usage: /join #channel")
			return
		}
		channelName := args[1]
		if channelName[0] != '#' {
			utils.SendBytes(user.conn, "channel names start by '#'")
			return
		}
		channel, ok := pickChannel(server, channelName)
		if !ok {
			utils.SendBytes(user.conn, "no channel by this name, try /list")
			return
		}
		if user.nick == "" {
			utils.SendBytes(user.conn, "you should pick a nickname first")
			return
		}
		user.channel = channel
		channel.users = append(channel.users, user)
		utils.SendBytes(user.conn, "/join "+channelName)
		utils.SendBytes(user.conn, fmt.Sprintf("welcome in channel %s!", channelName))
	case "/list":
		channelNames := []string{}
		for _, channel := range server.channels {
			channelNames = append(channelNames, channel.name)
		}
		utils.SendBytes(user.conn, fmt.Sprintf("%v", channelNames))
	case "/leave":
		if user.channel == nil {
			utils.SendBytes(user.conn, "you do not belong to any channel")
			return
		}
		channelName := user.channel.name
		user.channel = nil
		utils.SendBytes(user.conn, "/leave")
		utils.SendBytes(user.conn, "your left channel "+channelName)
	case "/who":
		if user.channel == nil {
			utils.SendBytes(user.conn, "you do not belong to any channel")
			return
		}
		channel, ok := pickChannel(server, user.channel.name)
		if !ok {
			fmt.Println("error during /who")
		}
		nicknames := []string{}
		for _, user := range channel.users {
			nicknames = append(nicknames, user.nick)
		}
		utils.SendBytes(user.conn, fmt.Sprintf("%v", nicknames))
	}
}

func readCommand(server *server, user *user) {
	scanner := bufio.NewScanner(user.conn)
	for {
		if ok := scanner.Scan(); !ok {
			break
		}
		go handleCommand(server, user, string(scanner.Text()))
	}
}

func handleRequest(server *server, conn net.Conn) {
	fmt.Printf("Serving %v\n", conn.RemoteAddr())
	user := user{conn: conn}
	readCommand(server, &user)
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
	server := newServer()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			os.Exit(1)
		}
		go handleRequest(&server, conn)
	}
}
