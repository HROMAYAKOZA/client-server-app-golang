package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strings"

	"client-server-app/pkg/utils"
)

type Server struct {
	Addr       string
	MaxClients int
	LogAddr    string
	Name       string
	Dispatcher *Dispatcher
	Verbose    bool

	Logger func(string, ...interface{})
}

func NewServer(d *Dispatcher, opts ...ServerOption) *Server {
	cfg := newConfig(opts...)
	s := &Server{
		Addr:       cfg.addr,
		MaxClients: cfg.maxClients,
		LogAddr:    cfg.logAddr,
		Name:       cfg.name,
		Verbose:    cfg.verbose,
		Dispatcher: d,
	}

	conn, err := net.Dial("tcp", s.LogAddr)
	s.Logger = utils.MakeLogger(&conn, s.Name, &s.Verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot connect to logger: %v\n", err)
		fmt.Println("Logs will be printed here")
	} else {
		s.Logger("Started listening on %s", s.Name)
		if s.Verbose {
			fmt.Println("Verbose enabled")
		}
	}
	s.Logger("Running server %s on architecture: %s", s.Name, runtime.GOARCH)
	s.Logger("Running server %s on system: %s", s.Name, runtime.GOOS)

	s.Dispatcher.Logger = s.Logger
	return s
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		s.Logger("listen error: " + err.Error())
		os.Exit(1)
	}
	s.Logger("Server running")
	sem := make(chan struct{}, s.MaxClients)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		select {
		case sem <- struct{}{}:
			go s.handleConn(conn, sem)
		default:
			utils.SendLine(conn, "ERROR: server busy")
			conn.Close()
		}
	}
}

func (s *Server) handleConn(conn net.Conn, sem chan struct{}) {
	defer conn.Close()
	defer func() { <-sem }()
	reader := bufio.NewReader(conn)
	s.Logger("Client %s connected", conn.RemoteAddr().String())
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			s.Logger("Client %s disconnected", conn.RemoteAddr().String())
			return
		}
		if err != nil {
			s.Logger("Read error: %s", err)
			return
		}
		s.Logger("Received cmd: %s", line)
		s.Dispatcher.Dispatch(s, strings.TrimSpace(line), conn)
	}
}

func (s *Server) SendLine(conn net.Conn, format string, args ...interface{}) error {
	err := utils.SendLine(conn, format, args...)
	msg := fmt.Sprintf(format, args...)
	s.Logger("Sent to client %s: %s", conn.RemoteAddr(), msg)
	return err
}
