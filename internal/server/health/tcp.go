package health

import (
	"bufio"
	"ep2/internal/server/router"
	"ep2/pkg/config"
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

	log.Printf("Escutando heartbeats TCP em %s:%s", ConnHost, ConnPort)
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

		payload := config.ParseMessageRead(netData)
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
	log.Printf("Conexão TCP fechada com cliente %s.", conn.RemoteAddr().String())
}

func shouldCloseTheConnection(response string) bool {
	return response == "BYE"
}
