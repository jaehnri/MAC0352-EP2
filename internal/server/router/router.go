package router

import (
	"ep2/internal/server/services"
	"fmt"
	"strings"
)

type Router struct {
	userService *services.UserService
}

type ServerRouter interface {
	HandleNew(params []string) string
	HandlePass(params []string) string
	//HandleIn(params []string) error
	//HandleOut(params []string) error
	//HandleL(params []string) error
	//HandlePlay(params []string) error
	//HandleHallOfFame(params []string) error
	//HandleCall(params []string) error
	//HandleOver(params []string) error
}

func NewRouter() *Router {
	return &Router{
		userService: services.NewUserService(),
	}
}

func (r *Router) Route(packet string) string {
	splitPacket := strings.Split(packet, " ")
	command := splitPacket[0]
	args := splitPacket[1:]

	switch command {
	case "new":
		return r.HandleNew(args)

	case "pass":
		return r.HandlePass(args)

	default:
		fmt.Printf("'%s' não é um comando conhecido.\n", command)
		return "ERROR"
	}
}

func (r *Router) HandleNew(params []string) string {
	err := r.userService.Create(params)
	if err != nil {
		return err.Error()
	}

	return "OK"
}

func (r *Router) HandlePass(params []string) string {
	err := r.userService.ChangePassword(params)
	if err != nil {
		return err.Error()
	}

	return "OK"
}
