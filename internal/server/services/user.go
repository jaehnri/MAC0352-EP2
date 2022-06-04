package services

import (
	"ep2/internal/server/repository"
	"ep2/pkg/model"
	"fmt"
	"log"
	"strconv"
	"strings"
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

func (u *UserService) GetUser(args []string) (*model.UserData, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ERRO: formato esperado é: get <user>.\n")
	}

	name := args[0]
	return u.repository.GetUser(name)
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

	err = u.repository.ChangeStatus(name, removePortFromIPAddress(address), Available)
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

	err := u.repository.ChangeStatus(name, removePortFromIPAddress(address), Offline)
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

func (u *UserService) Over(args []string) error {
	if len(args) != 4 {
		return fmt.Errorf("ERRO: formato esperado é: over <user1> <points1> <user2> <points2>.'\n")
	}

	user1 := args[0]
	user2 := args[2]
	pointsUser1, err := strconv.Atoi(args[1])
	if err != nil {
		log.Printf("ERRO: formato esperado de <points1> é inteiro.")
		return err
	}
	pointsUser2, err := strconv.Atoi(args[3])
	if err != nil {
		log.Printf("ERRO: formato esperado de <points2> é inteiro.")
		return err
	}

	err = u.repository.UpdatePoints(user1, pointsUser1)
	if err != nil {
		return err
	}
	err = u.repository.UpdatePoints(user2, pointsUser2)
	if err != nil {
		return err
	}

	err = u.repository.ChangeStatusWithoutAddress(user1, Available)
	if err != nil {
		return err
	}
	err = u.repository.ChangeStatusWithoutAddress(user1, Available)
	if err != nil {
		return err
	}

	printWinner(user1, user2, pointsUser1, pointsUser2)
	return nil
}

func (u *UserService) UpdateHeartbeat(args []string, address string) error {
	if len(args) == 0 {
		log.Printf("Heartbeat recebido de cliente em %s", removePortFromIPAddress(address))
		return nil
	}

	if len(args) != 1 {
		return fmt.Errorf("ERRO: formato esperado é: heartbeat <user1>.'\n")
	}

	name := args[0]
	return u.repository.UpdateHeartbeats(name, removePortFromIPAddress(address))
}

func printWinner(user1, user2 string, pointsUser1, pointsUser2 int) {
	if pointsUser1 > pointsUser2 {
		log.Printf("A partida entre <%s> e <%s> encerrou! O vencedor foi <%s>!", user1, user2, user1)
	}

	if pointsUser1 < pointsUser2 {
		log.Printf("A partida entre <%s> e <%s> encerrou! O vencedor foi <%s>!", user1, user2, user2)
	}

	if pointsUser1 == pointsUser2 {
		log.Printf("A partida entre <%s> e <%s> encerrou em empate.", user1, user2)
	}
}

func removePortFromIPAddress(address string) string {
	return strings.Split(address, ":")[0]
}
