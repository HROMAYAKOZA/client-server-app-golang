package server

import (
	"net"
	"strings"
)

type Command interface {
	Name() string
	Handle(s *Server, conn net.Conn, arg string) error
}

type Dispatcher struct {
	cmds   map[string]Command
	Logger func(string, ...interface{})
}

func NewDispatcher(cmds ...Command) *Dispatcher {
	m := make(map[string]Command, len(cmds))
	for _, c := range cmds {
		m[c.Name()] = c
	}
	return &Dispatcher{
		cmds:   m,
		Logger: nil,
	}
}

func (d *Dispatcher) Dispatch(s *Server, line string, conn net.Conn) {
	parts := strings.SplitN(line, ":", 2)
	name, arg := parts[0], ""
	if len(parts) == 2 {
		arg = parts[1]
	}
	if cmd, ok := d.cmds[name]; ok {
		err := cmd.Handle(s, conn, arg)
		if err != nil {
			d.Logger("ERROR on command %s: %s", cmd, err)
		}
	} else {
		s.SendLine(conn, "UNKNOWN CMD")
	}
}
