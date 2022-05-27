package services

import (
	"ep2/internal/data"
	"ep2/internal/repository"
	"fmt"
)

const (
	Offline   = "offline"
	Available = "online-available"
	Playing   = "online-playing"
)

type User interface {
	Create(name string, password string)
	ChangePassword(name string, password string)
	Login(name string)
	Logout(name string)
	ListConnected()
	ListAll()
}

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repository: repository.NewUserRepository(),
	}
}

func (u *UserService) Create(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("ERRO: formato esperado é: new <user> <senha>.\n")
	}

	return u.repository.Create(args[0], args[1])
}

func (u *UserService) ChangePassword(name string, oldpassword, newpassword string) error {
	// TODO: Implement password change
	return nil
}

func (u *UserService) Login(name string, password string) error {
	// TODO: Implement login
	return nil
}

func (u *UserService) Logout(name string) {
	// TODO: Implement logout
}

func (u *UserService) ListConnected() []data.UserData {
	// TODO: Implement list connected
	return nil
}

func (u *UserService) ListAll() []data.UserData {
	// TODO: Implement list all
	return nil
}

func (u *UserService) Get(username string) (data.UserData, error) {
	return data.UserData{}, nil
}
