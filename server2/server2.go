package main

import (
	"bufio"
	"client-server-app/utils"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
	"syscall"
)

const (
	addr2      = ":8002"
	maxClients = 5
	logAddr    = "localhost:9000"
	serverName = "Server2"
)

var (
	logConn net.Conn
	verbose bool
	logger  func(string, ...interface{})
)

func main() {
	flag.BoolVar(&verbose, "v", false, "включить подробный вывод")
	flag.Parse()
	// Подключаемся к серверу логирования
	var err error
	logConn, err = net.Dial("tcp", logAddr)
	logger = MakeLogger(&logConn, serverName, &verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot connect to logger: %v\n", err)
		// продолжаем без логирования
		fmt.Println("Logs will be printed here")
		logger("Running server 2 on architecture: %s", runtime.GOARCH)
	} else {
		logger("Started listening on %s", addr2)
		logger("Running server 2 on architecture: %s", runtime.GOARCH)
		if verbose {
			fmt.Println("Verbose enabled")
		}
	}

	listener, err := net.Listen("tcp", addr2)
	if err != nil {
		logger("listen error: " + err.Error())
		os.Exit(1)
	}
	logger("Server running")

	// Подобие семафора для ограничения числа клиентов
	sem := make(chan struct{}, maxClients)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		select {
		case sem <- struct{}{}:
			// Accepted
			go handleConn(conn, sem)
		default:
			// Too many clients
			utils.SendLine(conn, "ERROR: server busy")
			conn.Close()
		}
	}
}

func handleConn(conn net.Conn, sem chan struct{}) {
	defer conn.Close()
	defer func() { <-sem }()
	reader := bufio.NewReader(conn)
	logger("Client %s connected", conn.RemoteAddr().String())
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			logger("Client %s disconnected", conn.RemoteAddr().String())
			return
		}
		if err != nil {
			logger("Read error: %s", err)
			return
		}
		cmd := strings.TrimSpace(line)
		cmd = strings.ToLower(cmd)
		logger("Received cmd: %s", cmd)
		if strings.HasPrefix(cmd, "mem") {
			_, arg, _ := strings.Cut(cmd, "mem:")
			sendMemInfo(conn, arg)
		} else {
			utils.SendLine(conn, "UNKNOWN CMD")
		}
	}
}

func sendMemInfo(conn net.Conn, arg string) {
	if arg == "" {
		arg = "all"
	}
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		utils.SendLine(conn, "ERROR executing command: %s", err.Error())
		return
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
		utils.SendLine(conn, "Physical Memory: %.2f%%, Swap: %.2f%%", percRAM, percSwap)
	case "phy":
		utils.SendLine(conn, "Physical Memory: %.2f%%", percRAM)
	case "vir":
		utils.SendLine(conn, "Swap: %.2f%%", percSwap)
	default:
		utils.SendLine(conn, "WRONG ARGUMENT")
	}
}

func MakeLogger(logConn *net.Conn, serverName string, verbose *bool) func(string, ...interface{}) {
	return func(msg string, args ...interface{}) {
		utils.SendLog(logConn, serverName, verbose, msg, args...)
	}
}
