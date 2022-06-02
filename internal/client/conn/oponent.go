package conn

import (
	"bufio"
	"ep2/internal"
	"fmt"
	"net"
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
	return &OponentConnection{
		conn:   conn,
		writer: *bufio.NewWriter(conn),
		reader: *bufio.NewReader(conn),
	}, nil
}

func (c *OponentConnection) SendOver() error {
	return c.send("overed")
}

func (c *OponentConnection) SendPlay(i int, j int) error {
	return c.send(fmt.Sprintf("played %d %d", i, j))
}

func (c *OponentConnection) send(command string) error {
	before := time.Now()

	_, err := c.writer.WriteString(command + strconv.QuoteRune(internal.MessageDelim))

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
	str, err := c.reader.ReadString(internal.MessageDelim)
	return str, err
}

func (c *OponentConnection) Disconnect() error {
	return c.conn.Close()
}
