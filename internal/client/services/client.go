package services

import (
	"ep2/internal/client/conn"
	"ep2/internal/client/domain/game"
	"ep2/internal/server/services"
	"ep2/pkg/config"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type stateStruct struct {
	// login
	isLogged bool
	username string
	// connection
	inGame         bool
	game           *game.Game
	oponentConn    *conn.OponentConnection
	quitHearbeat   chan int
	oponentChannel chan string
}

type ClientService struct {
	state      *stateStruct
	serverConn *conn.ServerConnection
}

func NewClientService(serverConn *conn.ServerConnection) *ClientService {
	c := &ClientService{
		state: &stateStruct{
			isLogged:       false,
			inGame:         false,
			oponentChannel: make(chan string),
			quitHearbeat:   make(chan int),
		},
		serverConn: serverConn,
	}
	go c.receiveHeartbeats()
	go c.sendHeartbeats(c.state.quitHearbeat)
	return c
}

// /////////////////////////////////////////////////////////////////////
// USER
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) HandleNew(params []string) error {
	err := c.serverConn.Create(params[0], params[1])
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
	go c.listenOponent() // TODO: LISTEN TO NEW CONNECTIONS
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
	c.state.quitHearbeat <- 0
	c.serverConn.Logout(c.state.username)
	c.state.isLogged = false
	c.state.username = ""
	return nil
}
func (c *ClientService) HandleL(params []string) error {
	fmt.Println("Usuários conectados:")
	users, err := c.serverConn.Connected()
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
	users, err := c.serverConn.All()
	if err != nil {
		return err
	}
	for i, user := range users {
		fmt.Printf("%d. %s (%d pts)", i, user.Username, user.Points)
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////
// GAME
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) HandleCall(params []string) error {
	if !c.state.inGame {
		return errors.New("faça login antes de iniciar um jogo")
	}
	user, err := c.serverConn.Get(params[0])
	if err != nil {
		return err
	}
	if user.State != services.Available {
		return fmt.Errorf("o usuário '%s' não está disponível", user.Username)
	}
	c.state.oponentConn, err = conn.ConnectToClient(user.ConnectedIp, config.ClientPort)
	if err != nil {
		return err
	}
	go c.listenOponent()
	c.state.inGame = true
	c.state.game = game.NewGame(game.X) // TODO
	return nil
}
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
	c.handleTableChanged()
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
	c.handleTableChanged()
	return nil
}
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
func (c *ClientService) HandleOver(params []string) error {
	if !c.state.inGame {
		return errors.New("você não está em um jogo")
	}
	c.state.oponentConn.SendOver()
	c.state.oponentConn.Disconnect()
	c.state.oponentConn = nil
	c.state.inGame = false
	fmt.Println("Você se disconectou do jogo.")
	return nil
}
func (c *ClientService) HandleOvered(params []string) error {
	if !c.state.inGame {
		return nil
	}
	c.state.oponentConn.Disconnect()
	c.state.oponentConn = nil
	c.state.inGame = false
	fmt.Println("O oponente se disconectou do jogo.")
	return nil
}

func (c *ClientService) handleTableChanged() {
	c.state.game.PrintTable()
	switch c.state.game.State() {
	case game.Playing:
		return
	case game.Won:
		fmt.Println("Você ganhou!")
		c.serverConn.SendWon(c.state.username)
	case game.Draw:
		fmt.Println("Deu velha...")
		c.serverConn.SendDraw(c.state.username)
	case game.Lost:
		fmt.Println("Você perdeu...")
	}
	c.state.oponentConn.Disconnect()
	c.state.oponentConn = nil
}

// /////////////////////////////////////////////////////////////////////
// MORE
// /////////////////////////////////////////////////////////////////////

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
func (c *ClientService) HandleBye(params []string) error {
	if c.state.isLogged {
		return errors.New("você está logado, faça logout antes de sair")
	}
	fmt.Printf("Fechando o programa...")
	os.Exit(0)
	return nil
}

// /////////////////////////////////////////////////////////////////////
// ALTERNATE
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) AlternateListenTo() chan string {
	return c.state.oponentChannel
}

// /////////////////////////////////////////////////////////////////////
// CONCURRENCY
// /////////////////////////////////////////////////////////////////////

const heartbeatPeriod = 5 * time.Second
const maximumHeartbeat = 3 * time.Minute

func (c *ClientService) sendHeartbeats(quit chan int) {
	for {
		select {
		case <-time.After(heartbeatPeriod):
			fmt.Println("\nHeartbeat...") // TODO: remove
			c.serverConn.SendHeartbeat(c.state.username)
		case <-quit:
			return
		}
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

func (c *ClientService) listenOponent() {
	// TODO
}
