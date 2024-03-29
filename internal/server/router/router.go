package router

import (
	"encoding/json"
	"ep2/internal/server/services"
	"ep2/pkg/config"
	"log"
	"strings"
)

type Router struct {
	userService *services.UserService
}

type ServerRouter interface {
	HandleNew(params []string) string
	HandlePass(params []string) string
	HandleIn(params []string, address string) string
	HandleOut(params []string, address string) string
	HandleHallOfFame() string
	HandleL() string
	HandleGet(params []string) string
	HandlePlay(params []string) string
	HandleOver(params []string) string
	HandleHeartbeat(params []string) string
	HandleBye() string
}

func NewRouter() *Router {
	return &Router{
		userService: services.NewUserService(),
	}
}

func (r *Router) Route(packet string, address string) string {
	splitPacket := strings.Split(packet, " ")
	command := splitPacket[0]
	args := splitPacket[1:]

	switch command {
	case "new":
		return r.HandleNew(args)

	case "pass":
		return r.HandlePass(args)

	case "in":
		return r.HandleIn(args, address)

	case "out":
		return r.HandleOut(args, address)

	case "halloffame":
		return r.HandleHallOfFame()

	case "l":
		return r.HandleL()

	case "get":
		return r.HandleGet(args)

	case "play":
		return r.HandlePlay(args)

	case "over":
		return r.HandleOver(args)

	case "bye":
		return r.HandleBye()

	case "heartbeat":
		return r.HandleHeartbeat(args, address)

	default:
		log.Printf("'%s' não é um comando conhecido.", command)
		return "ERROR"
	}
}

func (r *Router) HandleNew(params []string) string {
	err := r.userService.Create(params)
	if err != nil {
		return err.Error()
	}

	return config.OK
}

func (r *Router) HandlePass(params []string) string {
	err := r.userService.ChangePassword(params)
	if err != nil {
		return err.Error()
	}

	return config.OK
}

func (r *Router) HandleIn(params []string, address string) string {
	err := r.userService.Login(params, address)
	if err != nil {
		return err.Error()
	}

	return config.OK
}

func (r *Router) HandleOut(params []string, address string) string {
	err := r.userService.Logout(params, address)
	if err != nil {
		return err.Error()
	}

	return config.OK
}

func (r *Router) HandleHallOfFame() string {
	users, err := r.userService.GetHallOfFame()
	if err != nil {
		return err.Error()
	}

	jsonResponse, err := json.Marshal(users)
	return string(jsonResponse)
}

func (r *Router) HandleL() string {
	users, err := r.userService.GetOnlineUsers()
	if err != nil {
		return err.Error()
	}

	jsonResponse, err := json.Marshal(users)
	return string(jsonResponse)
}

func (r *Router) HandleGet(params []string) string {
	user, err := r.userService.GetUser(params)
	if err != nil {
		return err.Error()
	}

	jsonResponse, err := json.Marshal(user)
	return string(jsonResponse)
}

func (r *Router) HandlePlay(params []string) string {
	err := r.userService.Play(params)
	if err != nil {
		return err.Error()
	}

	return config.OK
}

func (r *Router) HandleOver(params []string) string {
	err := r.userService.Over(params)
	if err != nil {
		return err.Error()
	}

	return config.OK
}

func (r *Router) HandleHeartbeat(params []string, address string) string {
	err := r.userService.UpdateHeartbeat(params, address)
	if err != nil {
		return err.Error()
	}

	return config.OK
}

func (r *Router) HandleBye() string {
	return "BYE"
}
