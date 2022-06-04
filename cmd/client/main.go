package main

import (
	"bufio"
	"ep2/cmd/client/mocks"
	"ep2/internal/client"
	"ep2/internal/client/conn"
	"ep2/internal/client/services"
	"ep2/pkg/config"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	serverConn := createServerConnection()
	clientService := services.NewClientService(serverConn)
	router := client.NewRouter()

	readTerminal := make(chan string)
	go concurrentlyReadLine(clientService.StdIn, readTerminal)

	for {
		fmt.Printf("JogoDaVelha> ")
		select {
		case line := <-readTerminal:
			handleError(router.Route(line, clientService))
		case line := <-clientService.Channels.OponentCommands:
			handleError(router.Route(line, clientService))
		case newConn := <-clientService.Channels.NewOponentConn:
			handleError(clientService.HandleCallRequest(newConn))
		}
	}
}

const (
	tcp = "tcp"
	udp = "udp"
)

func createServerConnection() conn.IServerConnection {
	if len(os.Args) < 3 {
		handleError(errors.New("ERRO: formato esperado é: ./client <server-ip> <server-port> <conn-type>"))
		os.Exit(1)
	}
	serverIp := os.Args[1]
	serverPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		handleError(errors.New("porta inválida"))
		os.Exit(1)
	}

	connType := tcp
	if len(os.Args) >= 4 {
		connType = os.Args[3]
	}

	var serverConn conn.IServerConnection
	switch connType {
	case tcp:
		serverConn, err = conn.TcpConnectToServer(serverIp, serverPort)
	case udp:
		serverConn, err = conn.UdpConnectToServer(serverIp, serverPort)
	case "mock":
		if len(os.Args) < 6 {
			handleError(errors.New("run: ./client localhost 0 mock <ip-listen> <ip-connect>"))
			os.Exit(1)
		}
		config.ClientPortListen, _ = strconv.Atoi(os.Args[4])
		config.ClientPortConnect, _ = strconv.Atoi(os.Args[5])
		serverConn = mocks.NewServerConnection()
	default:
		handleError(errors.New("tipo de conexao desconhecida"))
		os.Exit(1)
	}

	if err != nil {
		handleError(err)
		os.Exit(1)
	}
	return serverConn
}

func concurrentlyReadLine(scanner *bufio.Scanner, read chan string) {
	for {
		scanner.Scan()
		read <- strings.TrimSpace(scanner.Text())
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Erro: " + err.Error())
	}
}
