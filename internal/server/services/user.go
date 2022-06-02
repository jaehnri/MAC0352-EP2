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

	currentPasswordFromDatabase, err := u.repository.GetCurrentPassword(user)
	if err != nil {
		return err
	}

	if currentPasswordFromDatabase != currentPassword {
		return fmt.Errorf("ERRO: Usuário %s tentou alterar a sua senha mas errou a senha atual.\n", user)
	}

	return u.repository.ChangePassword(user, newPassword)
}

func (u *UserService) Login(args []string, address string) error {
	if len(args) != 2 {
		return fmt.Errorf("ERRO: formato esperado é: in <user> <senha>.\n")
	}

	name := args[0]
	password := args[1]

	currentPasswordFromDatabase, err := u.repository.GetCurrentPassword(name)
	if err != nil {
		return err
	}

	if currentPasswordFromDatabase != password {
		return fmt.Errorf("ERRO: Usuário <%s> errou a senha.\n", name)
	}

	err = u.repository.ChangeStatus(name, address, Available)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) Logout(args []string, address string) error {
	if len(args) != 1 {
		return fmt.Errorf("ERRO: formato esperado é: out <user>.\n")
	}
	name := args[0]

	err := u.repository.ChangeStatus(name, address, Offline)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetHallOfFame() ([]model.UserData, error) {
	return u.repository.HallOfFame()
}

func (u *UserService) GetOnlineUsers() ([]model.UserData, error) {
	return u.repository.GetOnlineUsers()
}

func (u *UserService) Play(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("ERRO: formato esperado é: play <user1> <user2>.\n")
	}
	user1 := args[0]
	user2 := args[1]

	return u.repository.Play(user1, user2, Playing)
}
