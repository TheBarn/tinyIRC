package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/TheBarn/tinyIRC/utils"
)

func main() {
	port := utils.ChoosePort()
	conn, err := net.Dial("tcp4", "localhost:"+port)
	if err != nil {
		fmt.Println("Error dialing:", err)
		os.Exit(1)
	}
	defer conn.Close()
	conn.Write([]byte("HEY FROM SERVER\n"))
	time.Sleep(5 * time.Second)
}
