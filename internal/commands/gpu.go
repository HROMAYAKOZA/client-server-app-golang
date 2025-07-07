package commands

import (
	"bufio"
	"bytes"
	"net"
	"os/exec"
	"strings"

	"client-server-app/internal/server"
)

type GPU struct{}

func (GPU) Name() string { return "gpu" }
func (GPU) Handle(s *server.Server, conn net.Conn, _ string) error {
	out, err := exec.Command("lspci").Output()
	if err != nil {
		s.SendLine(conn, "ERROR executing command: %s", err.Error())
		return err
	}
	var gpus []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "VGA compatible controller") || strings.Contains(line, "3D controller") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				gpus = append(gpus, strings.TrimSpace(parts[1]))
			}
		}
	}
	if len(gpus) == 0 {
		s.SendLine(conn, "GPU: not found")
	} else {
		s.SendLine(conn, "GPU: %s", strings.Join(gpus, ", "))
	}
	return nil
}
