package commands

import (
	"net"
	"syscall"

	"client-server-app/internal/server"
)

type Mem struct{}

func (Mem) Name() string { return "mem" }
func (Mem) Handle(s *server.Server, conn net.Conn, arg string) error {
	if arg == "" {
		arg = "all"
	}
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		s.SendLine(conn, "ERROR executing command: %s", err.Error())
		return err
	}
	totalRAM := float64(info.Totalram) * float64(info.Unit)
	freeRAM := float64(info.Freeram) * float64(info.Unit)
	usedRAM := totalRAM - freeRAM
	percRAM := usedRAM / totalRAM * 100

	totalSwap := float64(info.Totalswap) * float64(info.Unit)
	freeSwap := float64(info.Freeswap) * float64(info.Unit)
	usedSwap := totalSwap - freeSwap
	var percSwap float64
	if totalSwap > 0 {
		percSwap = usedSwap / totalSwap * 100
	}
	switch arg {
	case "all":
		s.SendLine(conn, "Physical Memory: %.2f%%, Swap: %.2f%%", percRAM, percSwap)
	case "phy":
		s.SendLine(conn, "Physical Memory: %.2f%%", percRAM)
	case "vir":
		s.SendLine(conn, "Swap: %.2f%%", percSwap)
	default:
		s.SendLine(conn, "WRONG ARGUMENT")
	}
	return nil
}
