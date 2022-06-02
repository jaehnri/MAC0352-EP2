package servers

import (
	"bytes"
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
	// Make a buffer to hold incoming data.
	buf := make([]byte, MaxPayloadSize)

	// Read the incoming connection into the buffer.
	qtdBytesRead, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Houve um erro ao ler um payload:", err.Error())
	}

	payload := parseTCPPayload(buf, qtdBytesRead)
	response := tcp.Router.Route(payload, conn.RemoteAddr().String())

	// Send a response back to client.
	conn.Write([]byte(response))

	// Close the TCP connection.
	conn.Close()
}

// Buf is always initialized as an array of 2048 bytes. However, most packets are sent with less than that.
// Here, we trim the byte array to use only how many bytes were actually read by conn.Read()
// and we remove the last 2 characters as they are a carriage feed (\r) and a line break (\n).
func parseTCPPayload(buf []byte, qtdBytesRead int) string {
	return bytes.NewBuffer(buf).String()[:qtdBytesRead-2]
}
