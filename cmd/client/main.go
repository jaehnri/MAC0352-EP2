package main

import (
	"bufio"
	services "ep2/internal/client/services"
	"fmt"
	"os"
	"strings"
)

func main() {
	client := services.NewClientService()
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
			handleError(client.HandleNew(params))
		case "in":
			handleError(client.HandleIn(params))
		case "pass":
			handleError(client.HandlePass(params))
		case "out":
			handleError(client.HandleOut(params))
		case "l":
			handleError(client.HandleL(params))
		case "halloffame":
			handleError(client.HandleHalloffame(params))

		// GAME
		case "call":
			handleError(client.HandleCall(params))
		case "play":
			handleError(client.HandlePlay(params))
			client.HandleTableChanged()
		case "delay":
			handleError(client.HandleDelay(params))
		case "over":
			handleError(client.HandleOver(params))

		// OTHER
		case "bye":
			fmt.Printf("Fechando o programa...")
			os.Exit(0)
		case "help":
			printCommandExplanation()
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

func printCommandExplanation() {
	fmt.Println("new <usuario> <senha>: cria um novo usuário")
	fmt.Println("pass <senha antiga> <senha nova>: muda a senha do usuário")
	fmt.Println("in <usuario> <senha>: usuário entra no servidor")
	fmt.Println("halloffame: informa a tabela de pontuação de todos os usuários registrados no sistema")
	fmt.Println("l: lista todos os usuários conectados no momento e se estão ocupados em uma partida ou não")
	fmt.Println("call <oponente>: convida um oponente para jogar. Ele pode aceitar ou não")
	fmt.Println("play <linha> <coluna>: envia a jogada")
	fmt.Println("delay: durante uma partida, informa os 3 últimos valores de latência que foram medidos para o cliente do oponente")
	fmt.Println("over: encerra uma partida antes da hora")
	fmt.Println("out: desloga")
	fmt.Println("bye: finaliza a execução do cliente e retorna para o shell do sistema operaciona")
	fmt.Println("help: mostra os comandos existentes")
}
