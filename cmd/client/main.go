package main

import (
	"bufio"
	"ep2/internal/client"
	"ep2/internal/client/conn"
	"ep2/internal/client/services"
	"fmt"
	"os"
	"strings"
)

func main() {
	serverConn, err := conn.TcpConnectToServer("", 0) // TODO
	if err != nil {
		handleError(err)
		os.Exit(1)
	}
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
