package mocks

import "ep2/pkg/model"

type MockServerConnection struct {
}

func NewServerConnection() *MockServerConnection {
	return &MockServerConnection{}
}

//////////////////////////////////////////////////////////////////////////////
// NULL
//////////////////////////////////////////////////////////////////////////////

func (c *MockServerConnection) SendStartedGame(username string, oponentUsername string) error {
	return nil
}
func (c *MockServerConnection) SendWon(username string, oponent string) error     { return nil }
func (c *MockServerConnection) SendDraw(username string, oponent string) error    { return nil }
func (c *MockServerConnection) SendOver(username string, oponent string) error    { return nil }
func (c *MockServerConnection) CreateUser(username string, password string) error { return nil }
func (c *MockServerConnection) ChangePassword(username string, oldpassword, newpassword string) error {
	return nil
}
func (c *MockServerConnection) Login(username string, password string) error { return nil }
func (c *MockServerConnection) Logout(username string) error                 { return nil }
func (c *MockServerConnection) SendHeartbeat(username string) error          { return nil }
func (c *MockServerConnection) Disconnect() error                            { return nil }

//////////////////////////////////////////////////////////////////////////////
// DATA
//////////////////////////////////////////////////////////////////////////////

var oponent model.UserData = model.UserData{
	Username: "Luca",
	State:    "online",
	Address:  "localhost",
	Points:   0,
}

func (c *MockServerConnection) OnlineUsers() ([]model.UserData, error) {
	return []model.UserData{oponent}, nil
}
func (c *MockServerConnection) AllUsers() ([]model.UserData, error) {
	return []model.UserData{oponent}, nil
}
func (c *MockServerConnection) GetUser(username string) (model.UserData, error) {
	return model.UserData{}, nil
}
func (c *MockServerConnection) ReadHeartbeat() (string, error) { return "OK", nil }
