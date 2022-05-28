package main

import (
	"bufio"
	"ep2/internal/client"
	"ep2/internal/client/services"
	"fmt"
	"os"
	"strings"
)

func main() {
	clientService := services.NewClientService()
	router := client.NewRouter()

	scanner := bufio.NewScanner(os.Stdin)
	readTerminal := make(chan string)
	go concurrentlyReadLine(scanner, readTerminal)

	for {
		fmt.Printf("JogoDaVelha> ")
		select {
		case line := <-readTerminal:
			handleError(router.Route(line, clientService))
		case line := <-clientService.AlternateListenTo():
			handleError(router.Route(line, clientService))
		}
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
