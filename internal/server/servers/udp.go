package servers

import (
	"bytes"
	"ep2/internal/server/router"
	"log"
	"net"
)

const (
	UDPHost   = "172.17.0.3"
	UDPPort   = 8080
	UDPPrefix = "udp"
)

type UDPServer struct {
	Router *router.Router
}

func NewUDPServer() *UDPServer {
	return &UDPServer{
		Router: router.NewRouter(),
	}
}

func (udp *UDPServer) StartUDPServer() {
	conn, err := net.ListenUDP(UDPPrefix, &net.UDPAddr{
		Port: UDPPort,
		IP:   net.ParseIP(UDPHost),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Escutando UDP em %s:%d", UDPHost, UDPPort)
	for {
		buf := make([]byte, 2048)
		length, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		payload := parseUDPPayload(buf, length)
		response := udp.Router.Route(payload, remote.String())

		conn.WriteToUDP([]byte(response), remote)
	}

	conn.Close()
}

func parseUDPPayload(payload []byte, length int) string {
	return bytes.NewBuffer(payload[:length-1]).String()
}
