package main

import (
	"ep2/internal/game"
	"ep2/internal/services"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type stateStruct struct {
	// login
	isLogged bool
	username string
	// connection
	inGame bool
	conn   net.Conn
	game   game.Game
}

type Client struct {
	state       stateStruct
	userService services.UserService
	gameService services.GameService
}

func NewClient() Client {
	return Client{
		state: stateStruct{
			isLogged: false,
			inGame:   false,
		},
		userService: *services.NewUserService(),
		gameService: *services.NewGameService(),
	}
}

// /////////////////////////////////////////////////////////////////////
// USER
// /////////////////////////////////////////////////////////////////////

func (c *Client) handleNew(params []string) error {
	// TODO:
	// port, errPort := strconv.ParseInt(params[0], 10, 32)
	// if errPort != nil {
	// 	return errPort
	// }
	// err := c.userService.Create(params[0], port)
	err := c.userService.Create(params)
	if err != nil {
		return err
	}
	fmt.Println("Usuário criado com sucesso")
	return nil
}
func (c *Client) handleIn(params []string) error {
	if c.state.isLogged {
		return errors.New("você já está logado, faça logout para trocar de usuário")
	}
	username := params[0]
	err := c.userService.Login(username, params[1])
	if err != nil {
		return err
	}
	go c.listenCalled()
	go c.heartbeat()
	c.state.isLogged = true
	c.state.username = username
	fmt.Printf("Você está logado como '%s'\n", username)
	return nil
}
func (c *Client) handlePass(params []string) error {
	if !c.state.isLogged {
		return errors.New("você não está logado")
	}
	err := c.userService.ChangePassword(c.state.username, params[0], params[1])
	if err != nil {
		return err
	}
	fmt.Println("Sua senha foi alterada.")
	return nil
}
func (c *Client) handleOut(params []string) error {
	if c.state.inGame {
		return errors.New("você está em um jogo")
	}
	if !c.state.isLogged {
		return errors.New("você não está logado")
	}
	c.userService.Logout(c.state.username)
	c.state.isLogged = false
	c.state.username = ""
	return nil
}
func (c *Client) handleL(params []string) error {
	fmt.Println("Usuários conectados:")
	for _, user := range c.userService.ListConnected() {
		fmt.Printf("• %s (%s)", user.Username, user.State)
	}
	return nil
}
func (c *Client) handleHalloffame(params []string) error {
	fmt.Println("Usuários conectados:")
	for i, user := range c.userService.ListConnected() {
		fmt.Printf("%d. %s (%d pts)", i, user.Username, user.Points)
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////
// GAME
// /////////////////////////////////////////////////////////////////////

func (c *Client) handleCall(params []string) error {
	if !c.state.inGame {
		return errors.New("faça login antes de iniciar um jogo")
	}
	user, err := c.userService.Get(params[0])
	if err != nil {
		return err
	}
	if user.State != services.Available {
		return fmt.Errorf("o usuário '%s' não está disponível", user.Username)
	}
	c.state.conn, err = c.gameService.Connect(user.ConnectedIp, user.ConnectedPort)
	if err != nil {
		return err
	}
	go c.listenOponent()
	c.state.inGame = true
	c.state.game = game.NewGame(game.X) // TODO
	return nil

}
func (c *Client) handlePlay(params []string) error {
	if !c.state.inGame {
		return errors.New("você não está em um jogo")
	}
	i, erri := strconv.ParseInt(params[0], 10, 32)
	j, errj := strconv.ParseInt(params[1], 10, 32)
	if erri != nil || errj != nil {
		return errors.New("posição inválida")
	}
	c.gameService.SendPlay(int(i), int(j))
	err := c.state.game.Play(int(i), int(j))
	if err != nil {
		return err
	}
	fmt.Printf("Você colocou %s em (%d,%d).\n", c.state.game.User, i, j)
	return nil
}
func (c *Client) handleDelay(params []string) error {
	if !c.state.inGame {
		return errors.New("você não está em um jogo")
	}
	fmt.Printf("A latência é de %d millisegundos.\n", c.gameService.Delay)
	return nil
}
func (c *Client) handleOver(params []string) error {
	if !c.state.inGame {
		return errors.New("você não está em um jogo")
	}
	c.gameService.SendDraw(c.state.username)
	c.gameService.Disconnect(c.state.conn)
	c.state.inGame = false
	c.state.conn = nil
	return nil
}

func (c *Client) handleTableChanged() {
	c.state.game.PrintTable()
	switch c.state.game.State() {
	case game.Playing:
		return
	case game.Won:
		fmt.Println("Você ganhou!")
		c.gameService.SendWon(c.state.username)
	case game.Draw:
		fmt.Println("Deu velha...")
		c.gameService.SendDraw(c.state.username)
	case game.Lost:
		fmt.Println("Você perdeu...")
	}
	c.gameService.Disconnect(c.state.conn)
}

// /////////////////////////////////////////////////////////////////////
// CONCURRENCY
// /////////////////////////////////////////////////////////////////////

func (c *Client) heartbeat() {
	// TODO: send and receive heartbeats (maximum 3 minutes)
}
func (c *Client) listenCalled() {
	// TODO
}
func (c *Client) listenOponent() {
	// TODO
}
