package main

import (
	"ep2/internal/server/health"
	"ep2/internal/server/servers"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	tcpServer := servers.NewTCPServer()
	go tcpServer.StartTCPServer()

	udpServer := servers.NewUDPServer()
	go udpServer.StartUDPServer()

	heartbeatTcpServer := health.NewHeartbeatTCPServer()
	go heartbeatTcpServer.StartHeartbeatTCPServer()

	// listens for system signals to gracefully shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	switch <-signalChannel {
	case os.Interrupt:
		log.Print("Foi recebido um SIGINT, finalizando servidor...")
	case syscall.SIGTERM:
		log.Print("Foi recebido um SIGTERM, finalizando servidor...")
	}
}
