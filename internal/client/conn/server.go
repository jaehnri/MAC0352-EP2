package conn

import (
	"bufio"
	"encoding/json"
	"ep2/pkg/config"
	"ep2/pkg/model"
	"errors"
	"fmt"
	"net"
	"strconv"
)

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

func (c ServerConnection) SendWon(username string) error {
	response, err := c.request(fmt.Sprintf("over %s 3", username))
	return handleVoidResponse(response, err)
}

func (c ServerConnection) SendDraw(username string) error {
	response, err := c.request(fmt.Sprintf("over %s 1", username))
	return handleVoidResponse(response, err)
}

/////////////////////////////////////////////////////////////////////////////////////
// User
/////////////////////////////////////////////////////////////////////////////////////

func (c ServerConnection) CreateUser(username string, password string) error {
	response, err := c.request(fmt.Sprintf("new %s %s", username, password))
	return handleVoidResponse(response, err)
}

func (c ServerConnection) ChangePassword(username string, oldpassword, newpassword string) error {
	response, err := c.request(fmt.Sprintf("pass %s %s %s", username, oldpassword, newpassword))
	return handleVoidResponse(response, err)
}

func (c ServerConnection) Login(username string, password string) error {
	response, err := c.request(fmt.Sprintf("in %s %s", username, password))
	return handleVoidResponse(response, err)
}

func (c ServerConnection) Logout(username string) error {
	response, err := c.request(fmt.Sprintf("out %s", username))
	return handleVoidResponse(response, err)
}

func (c ServerConnection) ConnectedUsers() ([]model.UserData, error) {
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

func (c ServerConnection) AllUsers() ([]model.UserData, error) {
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

func (c ServerConnection) GetUser(username string) (model.UserData, error) {
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

func (c ServerConnection) ReadHeartbeat() (string, error) {
	str, err := c.heartbeatReader.ReadString(config.MessageDelim)
	return str, err
}

func (c ServerConnection) SendHeartbeat(username string) error {
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

func (c ServerConnection) request(str string) (string, error) {
	err := c.send(str)
	if err != nil {
		return "", err
	}

	respose, err := c.read()
	return respose, err
}

func (c ServerConnection) send(str string) error {
	_, err := c.writer.WriteString(config.ParseWriteMessage(str))
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

func (c ServerConnection) read() (string, error) {
	str, err := c.reader.ReadString(config.MessageDelim)
	return config.ParseMessageRead(str), err
}

func (c ServerConnection) Disconnect() error {
	c.send("bye")
	return c.conn.Close()
}
