package main

import (
	"ep2/internal/servers"
)

func main() {
	tcpServer := servers.NewTCPServer()
	tcpServer.StartTCPServer()
}
