package servers

import (
	"bufio"
	"ep2/internal/server/router"
	"ep2/pkg/config"
	"log"
	"net"
	"os"
)

const (
	TCPHost   = "172.17.0.3"
	TCPPort   = "8080"
	TCPPrefix = "tcp"
)

type TCPServer struct {
	Router *router.Router
}

func NewTCPServer() *TCPServer {
	return &TCPServer{
		Router: router.NewRouter(),
	}
}

func (tcp *TCPServer) StartTCPServer() {
	// Listen for incoming connections.
	l, err := net.Listen(TCPPrefix, TCPHost+":"+TCPPort)
	if err != nil {
		log.Printf("Erro ao iniciar escuta TCP: %s", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	log.Printf("Escutando TCP em %s:%s", TCPHost, TCPPort)
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

func (tcp *TCPServer) handleRequest(conn net.Conn) {
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
}

func shouldCloseTheConnection(response string) bool {
	return response == "BYE"
}
