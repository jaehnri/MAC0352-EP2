package services

import "net"

// Talks to the client and server
type GameService struct {
	// TODO: Implement send play
	Delay int
}

func NewGameService() *GameService {
	return &GameService{Delay: 0}
}

func (g *GameService) Connect(ip string, port string) net.Conn {
	// TODO: Implement connect
	return nil
}

func (g *GameService) Disconnect(conn net.Conn) {
	// TODO: Implement connect
}

func (g *GameService) SendPlay(i int, j int) {
	// TODO: Implement send play
}
