package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type ServConnection struct {
	Conn   net.Conn
	Reader *bufio.Reader
	ID     int
}

type Client struct {
	maxConns int
	conns    []ServConnection
}

func New(maxConns int) *Client {
	return &Client{maxConns: maxConns}
}

func (c *Client) Len() int {
	return len(c.conns)
}

func (c *Client) Connect(server int, addr string) error {
	if server != 1 && server != 2 {
		return errors.New("SERVER NUMBER NOT SUPPORTED")
	}
	if len(c.conns) >= c.maxConns {
		return errors.New("CONNECTIONS LIMIT EXCEEDED")
	}
	port := 8000 + server
	conn, err := net.Dial("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	c.conns = append(c.conns, ServConnection{Conn: conn, Reader: bufio.NewReader(conn), ID: server})
	fmt.Printf("Connected to server %d\n", server)
	return nil
}

func (c *Client) DisconnectAll() error {
	for i, conn := range c.conns {
		err := conn.Conn.Close()
		if err != nil {
			return err
		}
		(c.conns)[i] = ServConnection{}
	}
	c.conns = (c.conns)[:0]
	return nil
}
func (c *Client) Disconnect(idx int) error {
	if idx < 0 || idx >= len(c.conns) {
		return errors.New("WRONG INDEX")
	}
	err := (c.conns)[0].Conn.Close()
	if err != nil {
		return err
	}
	c.conns = append((c.conns)[:idx], (c.conns)[idx+1:]...)
	return nil
}
func (c *Client) List() {
	for i, c := range c.conns {
		fmt.Println("Server", c.ID, "#", i)
	}
	if len(c.conns) == 0 {
		fmt.Println("No active connections")
	}
}
