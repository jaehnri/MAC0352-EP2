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
	scanner := bufio.NewScanner(os.Stdin)
	router := client.NewRouter()
	for {
		fmt.Printf("JogoDaVelha> ")
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		handleError(router.Route(line, clientService))
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
