package servers

import (
	"bufio"
	"ep2/internal/server/router"
	"fmt"
	"net"
	"os"
)

const (
	ConnHost       = "172.17.0.3"
	ConnPort       = "8080"
	ConnType       = "tcp"
	MaxPayloadSize = 2048
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
		// Handle connections in a new goroutine.
		go tcp.handleRequest(conn)
	}
}

func (tcp *TCPServer) handleRequest(conn net.Conn) {
	for {
		// Read the incoming data into a variable.
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if err != nil {
			fmt.Println("Houve um erro ao ler um payload: ", err.Error())
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
