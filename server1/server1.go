package main

import (
	"bufio"
	"bytes"
	"client-server-app/utils"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	addr1      = ":8001"
	hideMinMS  = 1000
	hideMaxMS  = 10000
	maxClients = 105
	logAddr    = "localhost:9000"
	serverName = "Server1"
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
		logger("Running server 1 on architecture: %s", runtime.GOARCH)
	} else {
		logger("Started listening on %s", addr1)
		// t := fmt.Sprint("Running server 1 on architecture:", runtime.GOARCH)
		logger("Running server 1 on architecture: %s", runtime.GOARCH)
		if verbose {
			fmt.Println("Verbose enabled")
		}
	}

	listener, err := net.Listen("tcp", addr1)
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
	defer func() { <-sem }() // release slot
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
		switch {
		case cmd == "gpu":
			sendGPUInfo(conn)
		case strings.HasPrefix(cmd, "hide:"):
			parts := strings.Split(cmd, ":")
			ms, err := strconv.Atoi(parts[1])
			if err != nil || ms < hideMinMS || ms > hideMaxMS {
				utils.SendLine(conn, "ERROR: time out of range")
			} else {
				if hideWindow(ms) {
					utils.SendLine(conn, "OK")
				} else {
					utils.SendLine(conn, "ERROR: hiding failed")
				}
			}
		default:
			utils.SendLine(conn, "UNKNOWN CMD")
		}
	}
}

func sendGPUInfo(conn net.Conn) {
	out, err := exec.Command("lspci").Output()
	if err != nil {
		utils.SendLine(conn, "ERROR executing command: %s", err.Error())
		return
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
		utils.SendLine(conn, "GPU: not found")
	} else {
		utils.SendLine(conn, "GPU: %s", strings.Join(gpus, ", "))
	}
}

func hideWindow(ms int) bool {
	// Simulate hide by sleeping
	cmd := exec.Command("sleep", strconv.Itoa(ms/1000))
	return cmd.Run() == nil
}

func MakeLogger(logConn *net.Conn, serverName string, verbose *bool) func(string, ...interface{}) {
	return func(msg string, args ...interface{}) {
		utils.SendLog(logConn, serverName, verbose, msg, args...)
	}
}
