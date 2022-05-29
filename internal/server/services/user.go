package services

import (
	"ep2/internal/server/repository"
	"ep2/pkg/model"
	"fmt"
)

const (
	Offline   = "offline"
	Available = "online-available"
	Playing   = "online-playing"
)

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

func (u *UserService) ChangePassword(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("ERRO: formato esperado é: pass <user> <senha antiga> <senha nova>.\n")
	}

	user := args[0]
	currentPassword := args[1]
	newPassword := args[2]

	currentPasswordFromDatabase, err := u.repository.GetOldPassword(user)
	if err != nil {
		return err
	}

	if currentPasswordFromDatabase != currentPassword {
		return fmt.Errorf("ERRO: Usuário %s tentou alterar a sua senha mas errou a senha atual.\n", user)
	}

	return u.repository.ChangePassword(user, newPassword)
}

func (u *UserService) Login(name string, password string) error {
	// TODO: Implement login
	return nil
}

func (u *UserService) Logout(name string) {
	// TODO: Implement logout
}

func (u *UserService) ListConnected() []model.UserData {
	// TODO: Implement list connected
	return nil
}

func (u *UserService) ListAll() []model.UserData {
	// TODO: Implement list all
	return nil
}

func (u *UserService) Get(username string) (model.UserData, error) {
	return model.UserData{}, nil
}
