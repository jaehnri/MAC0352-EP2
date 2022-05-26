package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type userState struct {
	// login
	isLogged bool
	username string
	// connection
	isConnected bool
	connection  io.Reader
	// game Game
}

// TODO remove (this can be added to empty blocks to facilitate development)
func none() {}

func main() {
	state := userState{
		isLogged:    false,
		username:    "",
		isConnected: false,
		connection:  nil,
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("JogoDaVelha> ")
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		words := strings.Split(line, " ")
		if line == "" || len(words) == 0 {
			continue
		}

		switch words[0] {
		// USER
		case "new", "l", "halloffame":
			// TODO
		case "in":
			if state.isLogged {
				fmt.Println("Você já está logado. Faça logout para trocar de usuário.")
			} else {
				// TODO
				state.isLogged = true
			}
		case "pass":
			if !state.isLogged {
				fmt.Println("Você não está logado.")
			} else {
				// TODO
				none()
			}
		case "out":
			if state.isConnected {
				fmt.Println("Você está em um jogo.")
			} else if !state.isLogged {
				fmt.Println("Você não está logado.")
			} else {
				state.isLogged = false
				state.username = ""
			}

		// GAME
		case "call":
			if !state.isConnected {
				fmt.Println("Faça login antes de iniciar um jogo.")
			} else {
				// TODO
				state.isConnected = true
			}
		case "play":
			if !state.isConnected {
				fmt.Println("Você não está em um jogo.")
			} else {
				// TODO
				none()
			}
		case "delay":
			if !state.isConnected {
				fmt.Println("Você não está em um jogo.")
			} else {
				// TODO
				none()
			}
		case "over":
			if !state.isConnected {
				fmt.Println("Você não está em um jogo.")
			} else {
				// TODO
				state.isConnected = false
			}

		// OTHER
		case "bye":
			fmt.Printf("Fechando o programa...")
			os.Exit(0)
		default:
			fmt.Printf("'%s' não é um comando conhecido.\n", words[0])
		}
	}
}
