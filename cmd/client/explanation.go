package main

import "fmt"

func PrintCommandExplanation() {
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
