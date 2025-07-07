package utils

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func SendLine(conn net.Conn, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	full := fmt.Sprintf("[%s] %s\n", timestamp, strings.TrimRight(msg, "\n"))
	_, err := conn.Write([]byte(full))
	return err
}

func SendLog(logConn *net.Conn, serverName string, verbose *bool, msg string, args ...interface{}) {
	msgf := fmt.Sprintf(msg, args...)
	timestamp := time.Now().Format("15:04:05")
	line := fmt.Sprintf("[%s] [%s] %s\n", serverName, timestamp, msgf)
	if *logConn != nil {
		_, err := (*logConn).Write([]byte(line))
		if err != nil {
			fmt.Println("Error accessing logger, switching local logs in terminal")
			*logConn = nil
			*verbose = false
		}
		if *verbose || *logConn == nil {
			fmt.Print(line)
		}
	} else {
		fmt.Print(line)
	}
}

func MakeLogger(logConn *net.Conn, serverName string, verbose *bool) func(string, ...interface{}) {
	return func(msg string, args ...interface{}) {
		SendLog(logConn, serverName, verbose, msg, args...)
	}
}
