package conn

import (
	"bufio"
	"ep2/pkg/config"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const amountLatency = 3

type OponentConnection struct {
	conn    net.Conn
	writer  bufio.Writer
	reader  bufio.Reader
	Latency []time.Duration
}

func ConnectToClient(ip string, port int) (*OponentConnection, error) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	return newOponentConnection(conn), nil
}
func WaitForOponentConnection() *OponentConnection {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(config.ClientPort))
	if err != nil {
		fmt.Println("Erro ao iniciar escuta:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	conn, _ := l.Accept()
	return newOponentConnection(conn)
}
func newOponentConnection(conn net.Conn) *OponentConnection {
	return &OponentConnection{
		conn:   conn,
		writer: *bufio.NewWriter(conn),
		reader: *bufio.NewReader(conn),
	}
}

//////////////////////////////////////////////////////////////
// SEND
//////////////////////////////////////////////////////////////

func (c *OponentConnection) SendOver() error {
	return c.send("overed")
}

func (c *OponentConnection) SendPlay(i int, j int) error {
	return c.send(fmt.Sprintf("played %d %d", i, j))
}

//////////////////////////////////////////////////////////////
// ACCEPT CONNECTION
//////////////////////////////////////////////////////////////

const (
	acceptGame = "accept"
	rejectGame = "reject"
)

func (c *OponentConnection) ReadGameAcceptance() (bool, error) {
	str, err := c.Read()
	return str == acceptGame, err
}

func (c *OponentConnection) SendAcceptGame() error {
	return c.send(acceptGame)
}

func (c *OponentConnection) SendRejectGame() error {
	return c.send(rejectGame)
}

func (c *OponentConnection) SendUsername(username string) error {
	return c.send(username)
}

//////////////////////////////////////////////////////////////
// CORE
//////////////////////////////////////////////////////////////

func (c *OponentConnection) send(command string) error {
	before := time.Now()

	_, err := c.writer.WriteString(config.ParseWriteMessage(command))

	after := time.Now()
	c.updateLatency(after.Sub(before))

	return err
}
func (c *OponentConnection) updateLatency(latency time.Duration) {
	if len(c.Latency) < amountLatency {
		c.Latency = append(c.Latency, latency)
		return
	}

	for i := 0; i < amountLatency-1; i++ {
		c.Latency[i] = c.Latency[i+1]
	}
	c.Latency[amountLatency-1] = latency
}

func (c *OponentConnection) Read() (string, error) {
	str, err := c.reader.ReadString(config.MessageDelim)
	return str, err
}

func (c *OponentConnection) Disconnect() error {
	return c.conn.Close()
}
