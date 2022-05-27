package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	client := NewClient()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("JogoDaVelha> ")
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		words := strings.Split(line, " ")
		if line == "" || len(words) == 0 {
			continue
		}
		params := words[1:]

		switch words[0] {
		// USER
		case "new":
			handleError(client.handleNew(params))
		case "in":
			handleError(client.handleIn(params))
		case "pass":
			handleError(client.handlePass(params))
		case "out":
			handleError(client.handleOut(params))
		case "l":
			handleError(client.handleL(params))
		case "halloffame":
			handleError(client.handleHalloffame(params))

		// GAME
		case "call":
			handleError(client.handleCall(params))
		case "play":
			handleError(client.handlePlay(params))
			client.handleTableChanged()
		case "delay":
			handleError(client.handleDelay(params))
		case "over":
			handleError(client.handleOver(params))

		// OTHER
		case "bye":
			fmt.Printf("Fechando o programa...")
			os.Exit(0)
		case "help":
			PrintCommandExplanation()
		default:
			fmt.Printf("'%s' não é um comando conhecido.\n", words[0])
		}
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
