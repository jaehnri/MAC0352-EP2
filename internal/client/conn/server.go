package conn

import (
	"bufio"
	"bytes"
	"encoding/json"
	"ep2/pkg/config"
	"ep2/pkg/model"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type IServerConnection interface {
	SendStartedGame(username string, oponentUsername string) error
	SendWon(username string, oponent string) error
	SendDraw(username string, oponent string) error
	SendOver(username string, oponent string) error
	CreateUser(username string, password string) error
	ChangePassword(username string, oldpassword, newpassword string) error
	Login(username string, password string) error
	Logout(username string) error
	OnlineUsers() ([]model.UserData, error)
	AllUsers() ([]model.UserData, error)
	GetUser(username string) (model.UserData, error)
	ReadHeartbeat() (string, error)
	SendHeartbeat(username string) error
	Disconnect() error
}

type ServerConnection struct {
	conn            net.Conn
	writer          bufio.Writer
	reader          bufio.Reader
	heartbeatConn   net.Conn
	heartbeatWriter bufio.Writer
	heartbeatReader bufio.Reader
}

/////////////////////////////////////////////////////////////////////////////////////
// New
/////////////////////////////////////////////////////////////////////////////////////

func TcpConnectToServer(ip string, port int) (*ServerConnection, error) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	heartbeatConn, err := net.Dial("tcp", ip+":"+strconv.Itoa(config.ServerHeartbeatPort))
	if err != nil {
		return nil, err
	}
	return newServerConnection(conn, heartbeatConn), nil
}

func UdpConnectToServer(ip string, port int) (*ServerConnection, error) {
	conn, err := net.Dial("udp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	heartbeatConn, err := net.Dial("udp", ip+":"+strconv.Itoa(config.ServerHeartbeatPort))
	if err != nil {
		return nil, err
	}
	return newServerConnection(conn, heartbeatConn), nil
}

func newServerConnection(conn net.Conn, heartbeatConn net.Conn) *ServerConnection {
	return &ServerConnection{
		conn:            conn,
		writer:          *bufio.NewWriter(conn),
		reader:          *bufio.NewReader(conn),
		heartbeatConn:   heartbeatConn,
		heartbeatWriter: *bufio.NewWriter(heartbeatConn),
		heartbeatReader: *bufio.NewReader(heartbeatConn),
	}
}

/////////////////////////////////////////////////////////////////////////////////////
// Methods
/////////////////////////////////////////////////////////////////////////////////////

func (c *ServerConnection) SendStartedGame(username string, oponentUsername string) error {
	response, err := c.request(fmt.Sprintf("play %s %s", username, oponentUsername))
	return handleVoidResponse(response, err)
}

func (c *ServerConnection) SendWon(username string, oponent string) error {
	return c.sendOver(username, 3, oponent, 0)
}

func (c *ServerConnection) SendDraw(username string, oponent string) error {
	return c.sendOver(username, 1, oponent, 1)
}

func (c *ServerConnection) SendOver(username string, oponent string) error {
	return c.sendOver(username, 0, oponent, 0)
}

func (c *ServerConnection) sendOver(username string, usernamePoints int, oponent string, oponentPoints int) error {
	response, err := c.request(fmt.Sprintf("over %s %d %s %d", username, usernamePoints, oponent, oponentPoints))
	return handleVoidResponse(response, err)
}

/////////////////////////////////////////////////////////////////////////////////////
// User
/////////////////////////////////////////////////////////////////////////////////////

func (c *ServerConnection) CreateUser(username string, password string) error {
	response, err := c.request(fmt.Sprintf("new %s %s", username, password))
	return handleVoidResponse(response, err)
}

func (c *ServerConnection) ChangePassword(username string, oldpassword, newpassword string) error {
	response, err := c.request(fmt.Sprintf("pass %s %s %s", username, oldpassword, newpassword))
	return handleVoidResponse(response, err)
}

func (c *ServerConnection) Login(username string, password string) error {
	response, err := c.request(fmt.Sprintf("in %s %s", username, password))
	return handleVoidResponse(response, err)
}

func (c *ServerConnection) Logout(username string) error {
	response, err := c.request(fmt.Sprintf("out %s", username))
	return handleVoidResponse(response, err)
}

func (c *ServerConnection) OnlineUsers() ([]model.UserData, error) {
	response, err := c.request("l")
	if err != nil {
		return nil, err
	}
	var users []model.UserData
	jsonErr := json.Unmarshal([]byte(response), &users)
	if jsonErr != nil {
		return nil, errors.New(response)
	}
	return users, nil
}

func (c *ServerConnection) AllUsers() ([]model.UserData, error) {
	response, err := c.request("halloffame")
	if err != nil {
		return nil, err
	}
	var users []model.UserData
	jsonErr := json.Unmarshal([]byte(response), &users)
	if jsonErr != nil {
		return nil, errors.New(response)
	}
	return users, nil
}

func (c *ServerConnection) GetUser(username string) (model.UserData, error) {
	response, err := c.request(fmt.Sprintf("get %s", username))
	if err != nil {
		return model.UserData{}, err
	}
	var user model.UserData
	jsonErr := json.Unmarshal([]byte(response), &user)
	if jsonErr != nil {
		return model.UserData{}, errors.New(response)
	}
	return user, nil
}

/////////////////////////////////////////////////////////////////////////////////////
// Heartbeat
/////////////////////////////////////////////////////////////////////////////////////

func (c *ServerConnection) ReadHeartbeat() (string, error) {
	str, err := c.heartbeatReader.ReadString(config.MessageDelim)
	return str, err
}

func (c *ServerConnection) SendHeartbeat(username string) error {
	text := "heartbeat"
	if username != "" {
		text += " " + username
	}
	_, err := c.heartbeatWriter.WriteString(config.ParseWriteMessage(text))
	return err
}

/////////////////////////////////////////////////////////////////////////////////////
// Core
/////////////////////////////////////////////////////////////////////////////////////

func (c *ServerConnection) request(str string) (string, error) {
	err := c.send(str)
	if err != nil {
		return "", err
	}

	response, err := c.read()
	return response, err
}

func (c *ServerConnection) send(str string) error {
	_, err := c.writer.WriteString(config.ParseWriteMessage(str))
	if err != nil {
		return err
	}

	err = c.writer.Flush()
	return err
}

func handleVoidResponse(response string, err error) error {
	if err != nil {
		return err
	}
	if response != config.OK {
		return errors.New(response)
	}
	return nil
}

func (c *ServerConnection) read() (string, error) {
	reply := make([]byte, 1024)
	length, err := c.conn.Read(reply)
	if err != nil {
		return "", err
	}

	response := bytes.NewBuffer(reply[:length]).String()
	return response, err
}

func (c *ServerConnection) Disconnect() error {
	c.send("bye")
	return c.conn.Close()
}
