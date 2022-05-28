package services

import (
	"ep2/internal/client/domain/game"
	"ep2/internal/client/repository"
	"ep2/internal/services"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type stateStruct struct {
	// login
	isLogged bool
	username string
	// connection
	inGame       bool
	conn         net.Conn
	game         game.Game
	quitHearbeat chan int
}

type ClientService struct {
	state               stateStruct
	userRepository      *repository.UserRepository
	gameRepository      *repository.GameRepository
	heartbeatRepository *repository.HeartbeatRepository
}

func NewClientService() *ClientService {
	return &ClientService{
		state: stateStruct{
			isLogged: false,
			inGame:   false,
		},
		userRepository:      repository.NewUserRepository(),
		gameRepository:      repository.NewGameRepository(),
		heartbeatRepository: repository.NewHeartbeatRepository(),
	}
}

// /////////////////////////////////////////////////////////////////////
// USER
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) HandleNew(params []string) error {
	err := c.userRepository.Create(params[0], params[1])
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
	err := c.userRepository.Login(username, params[1])
	if err != nil {
		return err
	}
	go c.listenCalled()
	c.state.quitHearbeat = make(chan int)
	go c.heartbeat(c.state.quitHearbeat)
	c.state.isLogged = true
	c.state.username = username
	fmt.Printf("Você está logado como '%s'\n", username)
	return nil
}
func (c *ClientService) HandlePass(params []string) error {
	if !c.state.isLogged {
		return errors.New("você não está logado")
	}
	err := c.userRepository.ChangePassword(c.state.username, params[0], params[1])
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
	c.userRepository.Logout(c.state.username)
	c.state.isLogged = false
	c.state.username = ""
	return nil
}
func (c *ClientService) HandleL(params []string) error {
	fmt.Println("Usuários conectados:")
	for _, user := range c.userRepository.Connected() {
		fmt.Printf("• %s (%s)", user.Username, user.State)
	}
	return nil
}
func (c *ClientService) HandleHalloffame(params []string) error {
	fmt.Println("Usuários conectados:")
	for i, user := range c.userRepository.All() {
		fmt.Printf("%d. %s (%d pts)", i, user.Username, user.Points)
	}
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
// GAME
// /////////////////////////////////////////////////////////////////////

func (c *ClientService) HandleCall(params []string) error {
	if !c.state.inGame {
		return errors.New("faça login antes de iniciar um jogo")
	}
	user, err := c.userRepository.Get(params[0])
	if err != nil {
		return err
	}
	if user.State != services.Available {
		return fmt.Errorf("o usuário '%s' não está disponível", user.Username)
	}
	c.state.conn, err = c.gameRepository.Connect(user.ConnectedIp, user.ConnectedPort)
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
	c.gameRepository.SendPlay(int(i), int(j))
	fmt.Printf("Você colocou %s em (%d,%d).\n", c.state.game.User, i, j)
	return nil
}
func (c *ClientService) HandleDelay(params []string) error {
	if !c.state.inGame {
		return errors.New("você não está em um jogo")
	}
	fmt.Printf("A latência é de %d millisegundos.\n", c.gameRepository.Delay)
	return nil
}
func (c *ClientService) HandleOver(params []string) error {
	if !c.state.inGame {
		return errors.New("você não está em um jogo")
	}
	c.gameRepository.SendDraw(c.state.username)
	c.gameRepository.Disconnect(c.state.conn)
	c.state.inGame = false
	c.state.conn = nil
	return nil
}

func (c *ClientService) HandleTableChanged() {
	c.state.game.PrintTable()
	switch c.state.game.State() {
	case game.Playing:
		return
	case game.Won:
		fmt.Println("Você ganhou!")
		c.gameRepository.SendWon(c.state.username)
	case game.Draw:
		fmt.Println("Deu velha...")
		c.gameRepository.SendDraw(c.state.username)
	case game.Lost:
		fmt.Println("Você perdeu...")
	}
	c.gameRepository.Disconnect(c.state.conn)
}

// /////////////////////////////////////////////////////////////////////
// CONCURRENCY
// /////////////////////////////////////////////////////////////////////

const heartbeatPeriod = 8 * time.Second

func (c *ClientService) heartbeat(quit chan int) {
	for {
		select {
		case <-time.After(heartbeatPeriod):
			// TODO: send and receive heartbeats (maximum 3 minutes)
			fmt.Println("\nHeartbeat...")
			c.heartbeatRepository.Send(c.state.username)
		case <-quit:
			return
		}
	}
}

func (c *ClientService) listenCalled() {
	// TODO
}

func (c *ClientService) listenOponent() {
	// TODO
}
