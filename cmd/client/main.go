package main

import (
	"bufio"
	"ep2/internal/client"
	"ep2/internal/client/conn"
	"ep2/internal/client/services"
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

	scanner := bufio.NewScanner(os.Stdin)
	readTerminal := make(chan string)
	go concurrentlyReadLine(scanner, readTerminal)

	for {
		fmt.Printf("JogoDaVelha> ")
		var line string
		select {
		case line = <-readTerminal:
		case line = <-clientService.AlternateListenTo():
		}
		handleError(router.Route(line, clientService))
	}
}

const tcp = "tcp"
const udp = "udp"

func createServerConnection() *conn.ServerConnection {
	if len(os.Args) < 3 {
		handleError(errors.New("envie todos os argumentos"))
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

	var serverConn *conn.ServerConnection
	switch connType {
	case tcp:
		serverConn, err = conn.TcpConnectToServer(serverIp, serverPort)
	case udp:
		serverConn, err = conn.UdpConnectToServer(serverIp, serverPort)
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
		fmt.Println(err)
	}
}
