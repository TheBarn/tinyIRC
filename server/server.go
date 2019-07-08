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
	users    []*user
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

func handleChanMessage(server *server, user *user, msg string) {
	channel := user.channel
	if channel == nil {
		utils.SendBytes(user.conn, "/warning You are not registered on a channel")
		return
	}
	for _, usr := range server.users {
		if usr.channel == channel {
			utils.SendBytes(usr.conn, fmt.Sprintf("/msg %s %s: %s", channel.name, user.nick, msg))
		}
	}
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

func removeUserFromChannel(user *user) {
	channel := user.channel
	if channel == nil {
		return
	}
	chanUsers := channel.users
	for idx, usr := range chanUsers {
		if usr.conn == user.conn {
			chanUsers[idx] = chanUsers[len(chanUsers)-1]
			channel.users = chanUsers[:len(chanUsers)-1]
		}
	}
}

func removeUserFromServer(server *server, user *user) {
	for idx, usr := range server.users {
		if usr == user {
			server.users[idx] = server.users[len(server.users)-1]
			server.users = server.users[:len(server.users)-1]
		}
	}
}

func handleCommand(server *server, user *user, cmd string) {
	if cmd == "" {
		return
	}
	if cmd[0] != '/' {
		handleChanMessage(server, user, cmd)
		return
	}
	args := strings.Fields(cmd)
	switch args[0] {
	case "/nick":
		if len(args) != 2 {
			utils.SendBytes(user.conn, "/warning usage: /nick nickname")
			return
		}
		nick := args[1]
		ok := checkNickname(nick)
		if !ok {
			utils.SendBytes(user.conn, "/warning nickname should have only 9 characters in [a-zA-Z0-9_]")
			return
		}
		user.nick = nick
		utils.SendBytes(user.conn, "/nick "+nick)
		utils.SendBytes(user.conn, "/warning your nickame was changed to "+nick)
	case "/join":
		if len(args) != 2 {
			utils.SendBytes(user.conn, "/warning usage: /join #channel")
			return
		}
		channelName := args[1]
		if channelName[0] != '#' {
			utils.SendBytes(user.conn, "/warning channel names start by '#'")
			return
		}
		channel, ok := pickChannel(server, channelName)
		if !ok {
			utils.SendBytes(user.conn, "/warning no channel by this name, try /list")
			return
		}
		if user.nick == "" {
			utils.SendBytes(user.conn, "/warning you should pick a nickname first")
			return
		}
		user.channel = channel
		channel.users = append(channel.users, user)
		utils.SendBytes(user.conn, "/join "+channelName)
		utils.SendBytes(user.conn, fmt.Sprintf("/warning welcome in channel %s!", channelName))
		for _, usr := range server.users {
			if usr.channel == channel && usr != user {
				utils.SendBytes(usr.conn, fmt.Sprintf("/msg %s : %s has joined the channel", channel.name, user.nick))
			}
		}
	case "/list":
		channelNames := []string{}
		for _, channel := range server.channels {
			channelNames = append(channelNames, channel.name)
		}
		utils.SendBytes(user.conn, fmt.Sprintf("%v", channelNames))
	case "/leave":
		if user.channel == nil {
			utils.SendBytes(user.conn, "/warning you do not belong to any channel")
			return
		}
		removeUserFromChannel(user)
		channel := user.channel
		user.channel = nil
		utils.SendBytes(user.conn, "/leave")
		utils.SendBytes(user.conn, "/warning you left channel "+channel.name)
		for _, usr := range server.users {
			if usr.channel == channel {
				utils.SendBytes(usr.conn, fmt.Sprintf("/msg %s : %s has left the channel", channel.name, user.nick))
			}
		}
	case "/who":
		if user.channel == nil {
			utils.SendBytes(user.conn, "/warning you do not belong to any channel")
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
	server.users = append(server.users, &user)
	readCommand(server, &user)
	fmt.Println("close connection", user)
	removeUserFromServer(server, &user)
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
