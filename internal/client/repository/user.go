package repository

import (
	"ep2/pkg/model"
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

func (u *UserRepository) Connected() []model.UserData {
	// TODO: Implement list connected
	return nil
}

func (u *UserRepository) All() []model.UserData {
	// TODO: Implement list all
	return nil
}

func (u *UserRepository) Get(username string) (model.UserData, error) {
	return model.UserData{}, nil
}
