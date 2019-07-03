package utils

import (
	"net"
)

//SendBytes send a '\n' terminated message through a TCP connection
func SendBytes(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg + "\n"))
	return err
}
