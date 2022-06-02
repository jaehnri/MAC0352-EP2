package main

import (
	"ep2/internal/server/health"
	"ep2/internal/server/servers"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	tcpServer := servers.NewTCPServer()
	go tcpServer.StartTCPServer()

	heartbeatTcpServer := health.NewHeartbeatTCPServer()
	go heartbeatTcpServer.StartHeartbeatTCPServer()

	// listens for system signals to gracefully shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	switch <-signalChannel {
	case os.Interrupt:
		fmt.Print("Foi recebido um SIGINT, finalizando servidor...")
	case syscall.SIGTERM:
		fmt.Print("Foi recebido um SIGTERM, finalizando servidor...")
	}
}
