package client

import (
	"ep2/internal/services"
	"fmt"
	"strings"
)

type Router struct {
	userService *services.UserService
}

type RouteMethods interface {
	HandleNew(params []string) error
	HandleIn(params []string) error
	HandlePass(params []string) error
	HandleOut(params []string) error
	HandleL(params []string) error
	HandleHalloffame(params []string) error
	HandleCall(params []string) error
	HandlePlay(params []string) error
	HandleDelay(params []string) error
	HandleOver(params []string) error
	HandleBye(params []string) error
	HandleHelp(params []string) error
}

func NewRouter() *Router {
	return &Router{
		userService: services.NewUserService(),
	}
}

func (r *Router) Route(line string, methods RouteMethods) error {
	words := strings.Split(line, " ")
	if line == "" || len(words) == 0 {
		return nil
	}
	params := words[1:]

	switch words[0] {
	// USER
	case "new":
		return methods.HandleNew(params)
	case "in":
		return methods.HandleIn(params)
	case "pass":
		return methods.HandlePass(params)
	case "out":
		return methods.HandleOut(params)
	case "l":
		return methods.HandleL(params)
	case "halloffame":
		return methods.HandleHalloffame(params)

	// GAME
	case "call":
		return methods.HandleCall(params)
	case "play":
		return methods.HandlePlay(params)
	case "delay":
		return methods.HandleDelay(params)
	case "over":
		return methods.HandleOver(params)

	// OTHER
	case "bye":
		return methods.HandleBye(params)
	case "help":
		return methods.HandleHelp(params)
	default:
		return fmt.Errorf("'%s' não é um comando conhecido", words[0])
	}
}
