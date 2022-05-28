package router

import (
	"ep2/internal/server/services"
	"fmt"
	"strings"
)

type Router struct {
	userService *services.UserService
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
		err := r.userService.Create(args)
		if err != nil {
			return err.Error()
		}
		return "OK"

	default:
		fmt.Printf("'%s' não é um comando conhecido.\n", command)
		return "ERROR"
	}
}
