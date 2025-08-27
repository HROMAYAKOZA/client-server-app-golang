package commands

import (
	"errors"
	"net"
	"os/exec"
	"strconv"

	"github.com/HROMAYAKOZA/client-server-app-golang/internal/server"
)

const (
	hideMinMS = 1000
	hideMaxMS = 10000
)

type Hide struct{}

func (Hide) Name() string { return "hide" }
func (Hide) Handle(s *server.Server, conn net.Conn, arg string) error {
	ms, err := strconv.Atoi(arg)
	if err != nil || ms < hideMinMS || ms > hideMaxMS {
		s.SendLine(conn, "ERROR: wrong timeout given")
		return errors.New("invalid input")
	}
	cmd := exec.Command("sleep", strconv.Itoa(ms/1000))
	cmdErr := cmd.Run()
	if cmdErr == nil {
		s.SendLine(conn, "OK")
	} else {
		s.SendLine(conn, "ERROR: hiding failed")
		return cmdErr
	}
	return nil
}
