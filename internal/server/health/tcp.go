package health

import (
	"bufio"
	"ep2/internal/server/router"
	"log"
	"net"
	"os"
)

const (
	ConnHost = "172.17.0.3"
	ConnPort = "8081"
	ConnType = "tcp"
)

type HeartbeatTCPServer struct {
	Router *router.Router
}

func NewHeartbeatTCPServer() *HeartbeatTCPServer {
	return &HeartbeatTCPServer{
		Router: router.NewRouter(),
	}
}

func (tcp *HeartbeatTCPServer) StartHeartbeatTCPServer() {
	// Listen for incoming connections.
	l, err := net.Listen(ConnType, ConnHost+":"+ConnPort)
	if err != nil {
		log.Printf("Erro ao iniciar escuta %s:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	log.Printf("Servidor TCP escutando em %s:%s", ConnHost, ConnPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Houve um erro ao aceitar: %s", err.Error())
			os.Exit(1)
		}

		log.Printf("Conexão TCP iniciada com %s!", conn.RemoteAddr().String())
		// Handle connections in a new goroutine.
		go tcp.handleRequest(conn)
	}
}

func (tcp *HeartbeatTCPServer) handleRequest(conn net.Conn) {
	for {
		// Read the incoming data into a variable.
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Houve um erro ao ler um payload: %s", err.Error())
			break
		}

		payload := parseTCPPayload(netData)
		response := tcp.Router.Route(payload, conn.RemoteAddr().String())

		// Check if connection should end.
		if shouldCloseTheConnection(response) {
			break
		}

		// Send a response back to client.
		conn.Write([]byte(response))
	}

	// Close the TCP connection.
	conn.Close()
}

// Here, we remove the last 2 characters as they are a carriage feed (\r) and a line break (\n).
func parseTCPPayload(buf string) string {
	return buf[:len(buf)-2]
}

func shouldCloseTheConnection(response string) bool {
	return response == "BYE"
}
