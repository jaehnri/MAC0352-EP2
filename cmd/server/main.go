package main

import (
	"ep2/internal/server/servers"
)

func main() {
	tcpServer := servers.NewTCPServer()
	tcpServer.StartTCPServer()
}
