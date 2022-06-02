package conn

import (
	"bufio"
	"encoding/json"
	"ep2/internal"
	"ep2/pkg/model"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type ServerConnection struct {
	conn   net.Conn
	writer bufio.Writer
	reader bufio.Reader
}

/////////////////////////////////////////////////////////////////////////////////////
// New
/////////////////////////////////////////////////////////////////////////////////////

func TcpConnectToServer(ip string, port int) (*ServerConnection, error) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	return newServerConnection(conn), nil
}

func UdpConnectToServer(ip string, port int) (*ServerConnection, error) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	return newServerConnection(conn), nil
}

func newServerConnection(conn net.Conn) *ServerConnection {
	return &ServerConnection{
		writer: *bufio.NewWriter(conn),
		reader: *bufio.NewReader(conn),
		conn:   conn,
	}
}

/////////////////////////////////////////////////////////////////////////////////////
// Methods
/////////////////////////////////////////////////////////////////////////////////////

func (c ServerConnection) SendWon(username string) error {
	return c.send(fmt.Sprintf("won %s", username))
}

func (c ServerConnection) SendDraw(username string) error {
	return c.send(fmt.Sprintf("draw %s", username))
}

func (c ServerConnection) SendHeartbeat(username string) error {
	return c.send(fmt.Sprintf("heartbeat %s", username))
}

func (c ServerConnection) Create(username string, password string) error {
	response, err := c.request(fmt.Sprintf("create %s %s", username, password))
	if err != nil {
		return err
	}
	if len(response) != 0 {
		return errors.New(response)
	}
	return nil
}

func (c ServerConnection) ChangePassword(username string, oldpassword, newpassword string) error {
	response, err := c.request(fmt.Sprintf("pass %s %s %s", username, oldpassword, newpassword))
	if err != nil {
		return err
	}
	if len(response) != 0 {
		return errors.New(response)
	}
	return nil
}

func (c ServerConnection) Login(username string, password string) error {
	response, err := c.request(fmt.Sprintf("login %s %s", username, password))
	if err != nil {
		return err
	}
	if len(response) != 0 {
		return errors.New(response)
	}
	return nil
}

func (c ServerConnection) Logout(name string) error {
	return c.send("logout")
}

func (c ServerConnection) Connected() ([]model.UserData, error) {
	response, err := c.request("connected")
	if err != nil {
		return nil, err
	}
	var users []model.UserData
	err = json.Unmarshal([]byte(response), &users)
	return users, err
}

func (c ServerConnection) All() ([]model.UserData, error) {
	response, err := c.request("get")
	if err != nil {
		return nil, err
	}
	var users []model.UserData
	err = json.Unmarshal([]byte(response), &users)
	return users, err
}

func (c ServerConnection) Get(username string) (model.UserData, error) {
	response, err := c.request(fmt.Sprintf("get %s", username))
	if err != nil {
		return model.UserData{}, err
	}
	var user model.UserData
	err = json.Unmarshal([]byte(response), &user)
	return user, err
}

/////////////////////////////////////////////////////////////////////////////////////
// Core
/////////////////////////////////////////////////////////////////////////////////////

func (s ServerConnection) request(str string) (string, error) {
	err := s.send(str)
	if err != nil {
		return "", err
	}

	respose, err := s.read()
	return respose, err
}

func (c ServerConnection) send(str string) error {
	_, err := c.writer.WriteString(str + strconv.QuoteRune(internal.MessageDelim))
	return err
}

func (c ServerConnection) read() (string, error) {
	str, err := c.reader.ReadString(internal.MessageDelim)
	return str, err
}

func (c ServerConnection) Disconnect() error {
	return c.conn.Close()
}
