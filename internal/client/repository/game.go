package repository

import "net"

// Talks to the client and server
type GameRepository struct {
	// TODO: Implement send play
	Delay int
}

func NewGameRepository() *GameRepository {
	return &GameRepository{Delay: 0}
}

func (g *GameRepository) Connect(ip string, port int) (net.Conn, error) {
	// TODO: Implement connect
	return nil, nil
}

func (g *GameRepository) Disconnect(conn net.Conn) {
	// TODO: Implement connect
}

func (g *GameRepository) SendPlay(i int, j int) {
	// TODO: Implement send play
}

func (g *GameRepository) SendWon(username string) {
	// TODO: Implement send play
}

func (g *GameRepository) SendDraw(username string) {
	// TODO: Implement send play
}
