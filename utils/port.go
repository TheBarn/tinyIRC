package utils

import (
	"fmt"
	"os"
)

const (
	defaultPort = "6667"
)

//ChoosePort will return the chosen port in args or the default port
func ChoosePort() string {
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
	return port
}
