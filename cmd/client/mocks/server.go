package mocks

import (
	"ep2/internal/server/services"
	"ep2/pkg/model"
	"errors"
	"time"
)

type MockServerConnection struct {
	users []model.UserData
}

func NewServerConnection() *MockServerConnection {
	conn := &MockServerConnection{
		users: make([]model.UserData, 0),
	}
	conn.users = append(conn.users, model.UserData{
		Username: "other",
		State:    services.Available,
		Address:  "localhost",
		Points:   0,
	})
	return conn
}

//////////////////////////////////////////////////////////////////////////////
// UPDATE
//////////////////////////////////////////////////////////////////////////////

func (c *MockServerConnection) SendStartedGame(username string, oponentUsername string) error {
	for i := 0; i < len(c.users); i++ {
		if c.users[i].Username == username || c.users[i].Username == oponentUsername {
			c.users[i].State = services.Playing
		}
	}
	return nil
}

func (c *MockServerConnection) SendWon(username string, oponent string) error {
	return c.sendOver(username, 3, oponent, 0)
}
func (c *MockServerConnection) SendDraw(username string, oponent string) error {
	return c.sendOver(username, 1, oponent, 1)
}
func (c *MockServerConnection) SendOver(username string, oponent string) error {
	return c.sendOver(username, 0, oponent, 0)
}
func (c *MockServerConnection) sendOver(username string, usernamePoints int, oponent string, oponentPoints int) error {
	for i := 0; i < len(c.users); i++ {
		if c.users[i].Username == username {
			c.users[i].Points += usernamePoints
		} else if c.users[i].Username == oponent {
			c.users[i].Points += oponentPoints
		}
	}
	return nil
}
func (c *MockServerConnection) Login(username string, password string) error {
	c.users = append(c.users, model.UserData{
		Username: username,
		State:    services.Available,
		Address:  "localhost",
		Points:   0,
	})
	return nil
}

//////////////////////////////////////////////////////////////////////////////
// NULL
//////////////////////////////////////////////////////////////////////////////

func (c *MockServerConnection) CreateUser(username string, password string) error { return nil }
func (c *MockServerConnection) ChangePassword(username string, oldpassword, newpassword string) error {
	return nil
}
func (c *MockServerConnection) Logout(username string) error        { return nil }
func (c *MockServerConnection) SendHeartbeat(username string) error { return nil }
func (c *MockServerConnection) Disconnect() error                   { return nil }

//////////////////////////////////////////////////////////////////////////////
// DATA
//////////////////////////////////////////////////////////////////////////////

func (c *MockServerConnection) OnlineUsers() ([]model.UserData, error) {
	return c.users, nil
}
func (c *MockServerConnection) AllUsers() ([]model.UserData, error) {
	return c.users, nil
}
func (c *MockServerConnection) GetUser(username string) (model.UserData, error) {
	for _, v := range c.users {
		if v.Username == username {
			return v, nil
		}
	}
	return model.UserData{}, errors.New("esse usuário não existe")
}
func (c *MockServerConnection) ReadHeartbeat() (string, error) {
	time.Sleep(2 * time.Second)
	return "OK", nil
}
