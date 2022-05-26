package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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
		case "call", "new", "pass", "in", "halloffame", "l", "play", "delay", "over", "out", "bye":
			fmt.Printf("Esse comando ainda não foi programado. Vindo em breve...\n")
		default:
			fmt.Printf("'%s' não é um comando conhecido.\n", words[0])
		}
	}
}
