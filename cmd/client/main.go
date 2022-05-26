package main

import (
	"bufio"
	"ep2/internal/game"
	"ep2/internal/services"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type userState struct {
	// login
	isLogged bool
	username string
	// connection
	inGame bool
	conn   net.Conn
	game   game.Game
}

func main() {
	state := userState{
		isLogged: false,
		inGame:   false,
	}

	userService := services.NewUserService()
	gameService := services.NewGameService()

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
		case "new":
			{
				err := userService.Create(words[1], words[2])
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Usuário criado com sucesso")
				}
			}
		case "in":
			if state.isLogged {
				fmt.Println("Você já está logado. Faça logout para trocar de usuário.")
			} else {
				username := words[1]
				err := userService.Login(username, words[2])
				if err != nil {
					fmt.Println(err)
				} else {
					state.isLogged = true
					state.username = username
					fmt.Printf("Você está logado como '%s'\n", username)
				}
			}
		case "pass":
			if !state.isLogged {
				fmt.Println("Você não está logado.")
			} else {
				err := userService.ChangePassword(state.username, words[1], words[2])
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Sua senha foi alterada.")
				}
			}
		case "out":
			if state.inGame {
				fmt.Println("Você está em um jogo.")
			} else if !state.isLogged {
				fmt.Println("Você não está logado.")
			} else {
				userService.Logout(state.username)
				state.isLogged = false
				state.username = ""
			}
		case "l":
			fmt.Println("Usuários conectados:")
			for _, user := range userService.ListConnected() {
				fmt.Printf("• %s (%s)", user.Username, user.State)
			}
		case "halloffame":
			fmt.Println("Usuários conectados:")
			for i, user := range userService.ListConnected() {
				fmt.Printf("%d. %s (%d pts)", i, user.Username, user.Points)
			}

		// GAME
		case "call":
			if !state.inGame {
				fmt.Println("Faça login antes de iniciar um jogo.")
			} else {
				user, err := userService.Get(words[1])
				if err != nil {
					fmt.Println(err)
				} else {
					if user.State != services.Availale {
						fmt.Printf("O usuário '%s' não está disponível.", user.Username)
					} else {
						state.conn = gameService.Connect(user.ConnectedIp, user.ConnectedPort)
						state.inGame = true
						state.game = game.NewGame(game.X) // TODO
					}
				}
			}
		case "play":
			if !state.inGame {
				fmt.Println("Você não está em um jogo.")
			} else {
				i, erri := strconv.ParseInt(words[1], 10, 32)
				j, errj := strconv.ParseInt(words[2], 10, 32)
				if erri != nil || errj != nil {
					fmt.Println("Posição inválida.")
				} else {
					gameService.SendPlay(int(i), int(j))
					err := state.game.Play(int(i), int(j))
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Printf("Você colocou %s em (%d,%d).\n", state.game.User, i, j)
					}
				}
			}
		case "delay":
			if !state.inGame {
				fmt.Println("Você não está em um jogo.")
			} else {
				fmt.Printf("A latência é de %d millisegundos.\n", gameService.Delay)
			}
		case "over":
			if !state.inGame {
				fmt.Println("Você não está em um jogo.")
			} else {
				gameService.Disconnect(state.conn)
				state.inGame = false
				state.conn = nil
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
