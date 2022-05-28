package services

import (
	"ep2/internal/data"
)

const (
	Offline   = "offline"
	Available = "online-available"
	Playing   = "online-playing"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) Create(name string, password string) error {
	// TODO: Implement create
	return nil
}

func (u *UserRepository) ChangePassword(name string, oldpassword, newpassword string) error {
	// TODO: Implement password change
	return nil
}

func (u *UserRepository) Login(name string, password string) error {
	// TODO: Implement login
	return nil
}

func (u *UserRepository) Logout(name string) {
	// TODO: Implement logout
}

func (u *UserRepository) Connected() []data.UserData {
	// TODO: Implement list connected
	return nil
}

func (u *UserRepository) All() []data.UserData {
	// TODO: Implement list all
	return nil
}

func (u *UserRepository) Get(username string) (data.UserData, error) {
	return data.UserData{}, nil
}
