package health

import (
	"bufio"
	"ep2/internal/server/router"
	"fmt"
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
		fmt.Println("Erro ao iniciar escuta:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Escutando em " + ConnHost + ":" + ConnPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Houve um erro ao aceitar: ", err.Error())
			os.Exit(1)
		}

		fmt.Printf("Conexão TCP iniciada com %s!", conn.RemoteAddr().String())
		// Handle connections in a new goroutine.
		go tcp.handleRequest(conn)
	}
}

func (tcp *HeartbeatTCPServer) handleRequest(conn net.Conn) {
	for {
		// Read the incoming data into a variable.
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Houve um erro ao ler um payload: ", err.Error())
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
