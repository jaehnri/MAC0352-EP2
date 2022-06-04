package health

import (
	"bytes"
	"ep2/internal/server/router"
	"log"
	"net"
)

const (
	UDPHost   = "172.17.0.3"
	UDPPort   = 8081
	UDPPrefix = "udp"
)

type HeartbeatUDPServer struct {
	Router *router.Router
}

func NewHeartbeatUDPServer() *HeartbeatUDPServer {
	return &HeartbeatUDPServer{
		Router: router.NewRouter(),
	}
}

func (udp *HeartbeatUDPServer) StartHeartbeatUDPServer() {
	conn, err := net.ListenUDP(UDPPrefix, &net.UDPAddr{
		Port: UDPPort,
		IP:   net.ParseIP(UDPHost),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Escutando heartbeats UDP em %s:%d", UDPHost, UDPPort)
	for {
		buf := make([]byte, 2048)
		length, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		payload := parseUDPPayload(buf, length)
		response := udp.Router.Route(payload, remote.String())

		conn.Write([]byte(response))
	}

	conn.Close()
}

func parseUDPPayload(payload []byte, length int) string {
	return bytes.NewBuffer(payload[:length-1]).String()
}
