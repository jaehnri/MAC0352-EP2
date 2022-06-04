package services

import (
	"bufio"
	"ep2/internal/client/conn"
	"ep2/internal/client/domain/game"
	"ep2/internal/server/services"
	"ep2/pkg/config"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type stateStruct struct {
	// login
	isLogged bool
	username string
	// connection
	inGame          bool
	game            *game.Game
	oponentConn     *conn.OponentConnection
	oponentUsername string
}

type ClientChannels struct {
	OponentCommands chan string
	NewOponentConn  chan *conn.OponentConnection
	quitOponentConn chan string
}

type ClientService struct {
	state      *stateStruct
	serverConn *conn.ServerConnection
	StdIn      *bufio.Scanner
	Channels   *ClientChannels
}

func NewClientService(serverConn *conn.ServerConnection) *ClientService {
	c := &ClientService{
		state: &stateStruct{
			isLogged: false,
			inGame:   false,
		},
		Channels: &ClientChannels{
			OponentCommands: make(chan string),
			NewOponentConn:  make(chan *conn.OponentConnection),
			quitOponentConn: make(chan string),
		},
		serverConn: serverConn,
		StdIn:      bufio.NewScanner(os.Stdin),
	}
	go c.receiveHeartbeats()
	go c.sendHeartbeats()
	go c.acceptOponentConn()
	return c
}

// /////////////////////////////////////////////////////////////////////
// USER
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) HandleNew(params []string) error {
	err := c.serverConn.CreateUser(params[0], params[1])
	if err != nil {
		return err
	}
	fmt.Println("Usuário criado com sucesso")
	return nil
}

func (c *ClientService) HandleIn(params []string) error {
	if c.state.isLogged {
		return errors.New("você já está logado, faça logout para trocar de usuário")
	}
	username := params[0]
	err := c.serverConn.Login(username, params[1])
	if err != nil {
		return err
	}
	c.state.isLogged = true
	c.state.username = username
	fmt.Printf("Você está logado como '%s'\n", username)
	return nil
}

func (c *ClientService) HandlePass(params []string) error {
	if !c.state.isLogged {
		return errors.New("você não está logado")
	}
	err := c.serverConn.ChangePassword(c.state.username, params[0], params[1])
	if err != nil {
		return err
	}
	fmt.Println("Sua senha foi alterada.")
	return nil
}

func (c *ClientService) HandleOut(params []string) error {
	if c.state.inGame {
		return errors.New("você está em um jogo")
	}
	if !c.state.isLogged {
		return errors.New("você não está logado")
	}
	c.serverConn.Logout(c.state.username)
	c.state.isLogged = false
	return nil
}

func (c *ClientService) HandleL(params []string) error {
	fmt.Println("Usuários conectados:")
	users, err := c.serverConn.OnlineUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		fmt.Printf("• %s (%s)", user.Username, user.State)
	}
	return nil
}

func (c *ClientService) HandleHalloffame(params []string) error {
	fmt.Println("Usuários conectados:")
	users, err := c.serverConn.AllUsers()
	if err != nil {
		return err
	}
	for i, user := range users {
		fmt.Printf("%d. %s (%d pts)", i, user.Username, user.Points)
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////
// CALL
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) HandleCall(params []string) error {
	if c.state.inGame {
		return errors.New("você já está jogando")
	}
	oponentName := params[0]
	if oponentName == c.state.username {
		return errors.New("você não pode jogar contra si mesmo")
	}
	oponent, err := c.serverConn.GetUser(oponentName)
	if err != nil {
		return err
	}
	if oponent.State != services.Available {
		return fmt.Errorf("o usuário '%s' não está disponível", oponent.Username)
	}

	// connect
	c.state.oponentConn, err = conn.ConnectToClient(oponent.Address, config.ClientPort)
	if err != nil {
		return err
	}

	// send username
	err = c.state.oponentConn.SendUsername(c.state.username)
	if err != nil {
		return err
	}

	// read acceptance response (agree or not)
	fmt.Println("Aguardando a resposta do oponente...")
	accepted, err := c.state.oponentConn.ReadGameAcceptance()
	if err != nil {
		return err
	}

	// if not accepted, close the connection
	if !accepted {
		fmt.Println("O oponente rejeitou o jogo.")
		c.state.oponentConn.Disconnect()
		return nil
	}

	// start the game
	fmt.Println("O oponente aceitou o jogo.")
	c.serverConn.SendStartedGame(c.state.username, oponent.Username)
	c.startGame(game.X, oponent.Username)
	return nil
}

func (c *ClientService) HandleCallRequest(newOponentConn *conn.OponentConnection) error {
	oponentUsername, err := newOponentConn.Read()
	if err != nil {
		return err
	}

	var accepted bool
	if c.state.inGame {
		accepted = false
	} else {
		fmt.Printf("Você deseja jogar com %s?[s/n] ", oponentUsername)
		for {
			c.StdIn.Scan()
			text := strings.ToLower(c.StdIn.Text()[0:1])
			accepted = (text == "s")
			if text == "s" || text == "n" {
				fmt.Println("Resposta desconhecida.")
			} else {
				break
			}
		}
	}

	if accepted {
		c.state.oponentConn = newOponentConn
		c.state.oponentConn.SendAcceptGame()
		c.startGame(game.O, oponentUsername)
	} else {
		newOponentConn.SendRejectGame()
		newOponentConn.Disconnect()
	}
	return nil
}

func (c *ClientService) acceptOponentConn() {
	for {
		oponentConn := conn.WaitForOponentConnection()
		c.Channels.NewOponentConn <- oponentConn
	}
}

func (c *ClientService) startGame(userSymbol string, oponentUsername string) {
	c.state.inGame = true
	c.state.oponentUsername = oponentUsername
	c.state.game = game.NewGame(userSymbol)
	go c.listenOponent()
}

// /////////////////////////////////////////////////////////////////////
// PLAY
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) HandlePlay(params []string) error {
	if !c.state.inGame {
		return errors.New("você não está em um jogo")
	}
	i, erri := strconv.ParseInt(params[0], 10, 32)
	j, errj := strconv.ParseInt(params[1], 10, 32)
	if erri != nil || errj != nil {
		return errors.New("posição inválida")
	}
	err := c.state.game.Play(int(i), int(j))
	if err != nil {
		return err
	}
	c.handleTableChanged(true)
	fmt.Printf("Você colocou %s em (%d,%d).\n", c.state.game.User, i, j)
	c.state.oponentConn.SendPlay(int(i), int(j))
	return nil
}

func (c *ClientService) HandlePlayed(params []string) error {
	if !c.state.inGame {
		return nil
	}
	i, erri := strconv.ParseInt(params[0], 10, 32)
	j, errj := strconv.ParseInt(params[1], 10, 32)
	if erri != nil || errj != nil {
		return errors.New("posição inválida")
	}
	err := c.state.game.OponentPlayed(int(i), int(j))
	if err != nil {
		return err
	}
	fmt.Printf("O oponente jogou em (%d,%d).\n", i, j)
	c.handleTableChanged(false)
	return nil
}

func (c *ClientService) listenOponent() {
	readFromOponent := make(chan string)
	go func() {
		for {
			str, _ := c.state.oponentConn.Read()
			readFromOponent <- str
		}
	}()

	for {
		select {
		case str := <-readFromOponent:
			c.Channels.OponentCommands <- str
		case <-c.Channels.quitOponentConn:
			return
		}
	}
}

func (c *ClientService) handleTableChanged(userPlayed bool) {
	c.state.game.PrintTable()
	gameState := c.state.game.State()
	if gameState == game.Playing {
		return
	}

	switch gameState {
	case game.Won:
		fmt.Println("Você ganhou!")
	case game.Draw:
		fmt.Println("Deu velha...")
	case game.Lost:
		fmt.Println("Você perdeu...")
	}

	if userPlayed {
		switch gameState {
		case game.Won:
			c.serverConn.SendWon(c.state.username, c.state.oponentUsername)
		case game.Draw:
			c.serverConn.SendDraw(c.state.username, c.state.oponentUsername)
		}
	}

	c.endGame()
}

// /////////////////////////////////////////////////////////////////////
// OVER
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) HandleOver(params []string) error {
	if !c.state.inGame {
		return errors.New("você não está em um jogo")
	}
	c.state.oponentConn.SendOver()
	c.serverConn.SendOver(c.state.username, c.state.oponentUsername)
	c.endGame()
	fmt.Println("Você se disconectou do jogo.")
	return nil
}

func (c *ClientService) HandleOvered(params []string) error {
	if !c.state.inGame {
		return nil
	}
	c.endGame()
	fmt.Println("O oponente se disconectou do jogo.")
	return nil
}

func (c *ClientService) endGame() {
	c.state.oponentConn.Disconnect()
	c.state.inGame = false
}

// /////////////////////////////////////////////////////////////////////
// MORE
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) HandleDelay(params []string) error {
	if !c.state.inGame {
		return errors.New("você não está em um jogo")
	}
	fmt.Println("A latência das últimas mensagens foram:")
	for _, latency := range c.state.oponentConn.Latency {
		fmt.Printf("- %d milisegundos\n", latency)
	}
	return nil
}

func (c *ClientService) HandleBye(params []string) error {
	if c.state.isLogged {
		return errors.New("você está logado, faça logout antes de sair")
	}
	c.serverConn.Disconnect()
	fmt.Printf("Fechando o programa...")
	os.Exit(0)
	return nil
}

func (c *ClientService) HandleHelp(params []string) error {
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
	return nil
}

// /////////////////////////////////////////////////////////////////////
// HEARTBEAT
// /////////////////////////////////////////////////////////////////////

const heartbeatPeriod = 5 * time.Second
const maximumHeartbeat = 3 * time.Minute

func (c *ClientService) sendHeartbeats() {
	for {
		time.Sleep(heartbeatPeriod)
		c.serverConn.SendHeartbeat(c.state.username)
	}
}

func (c *ClientService) receiveHeartbeats() {
	read := make(chan string)
	go c.readHeartbeats(read)
	for {
		select {
		case <-time.After(maximumHeartbeat):
			fmt.Println("Erro: O servidor não está disponível.")
			os.Exit(1)
		case <-read:
		}
	}
}
func (c *ClientService) readHeartbeats(read chan string) {
	for {
		str, err := c.serverConn.ReadHeartbeat()
		if err != nil {
			read <- str
		}
	}
}
