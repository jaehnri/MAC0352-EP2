package health

import (
	"ep2/internal/server/services"
	"log"
	"time"
)

const (
	HeartbeatFrequency = 30 * time.Second
)

type HeartbeatCronjob struct {
	userService *services.UserService
}

func NewHeartbeatCronjob() *HeartbeatCronjob {
	return &HeartbeatCronjob{
		userService: services.NewUserService(),
	}
}

func (h *HeartbeatCronjob) StartHeartbeatCronjob() {
	log.Printf("Checando usuários conectados a cada %s.", HeartbeatFrequency.String())
	for {
		time.Sleep(HeartbeatFrequency)
		h.checkOnlineUsers()
	}
}

func (h *HeartbeatCronjob) checkOnlineUsers() {
	onlineUsers, err := h.userService.GetOnlineUsers()
	if err != nil {
		log.Printf("Não foi possível resgatar usuários online do banco.")
		return
	}

	for _, user := range onlineUsers {
		now := time.Now().UTC()
		lastHeartbeat := user.LastHeartbeat.UTC()

		diff := now.Sub(lastHeartbeat)
		if diff > HeartbeatFrequency {
			err = h.userService.DisconnectUser(user.Username)
			if err != nil {
				log.Printf("Não foi possível desconectar o usuário <%s>.", user.Username)
			}
			log.Printf("Usuário <%s> desconectado via checagem de heartbeat. Último heartbeat foi há %s.", user.Username, diff)
		}
	}

	log.Print("Checagem de heartbeat completa!")
}
