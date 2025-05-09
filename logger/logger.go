package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var verbose bool

func createSessionLogDir() (string, error) {
	baseDir := "logs"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	// уникальное имя сессии на основе времени запуска
	sessionName := "session_" + time.Now().Format("20060102_150405")
	sessionPath := filepath.Join(baseDir, sessionName)

	if err := os.MkdirAll(sessionPath, 0755); err != nil {
		return "", err
	}

	return sessionPath, nil
}

func main() {
	flag.BoolVar(&verbose, "v", false, "включить подробный вывод")
	flag.Parse()
	logDir, err := createSessionLogDir()
	if err != nil {
		fmt.Println("Ошибка создания каталога логов:", err)
		return
	}
	fmt.Println("Логгер запущен. Логи сохраняются в:", logDir)

	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Ошибка запуска логгера:", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			continue
		}
		go handleConnection(conn, logDir)
	}
}

func handleConnection(conn net.Conn, logDir string) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	var server string
	defer disconnect(logDir, server)

	for scanner.Scan() {
		text := scanner.Text()
		// Ожидаем, что лог в формате: [ServerName] message
		parts := strings.SplitN(text, "] ", 2)
		if len(parts) != 2 {
			continue
		}

		server = strings.TrimPrefix(parts[0], "[")
		message := strings.TrimSpace(parts[1])

		filename := filepath.Join(logDir, server+".log")
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Ошибка открытия файла лога:", err)
			continue
		}
		if _, err := f.WriteString(message); err != nil {
			fmt.Println("Ошибка записи лога:", err)
		}
		if _, err := f.WriteString("\n"); err != nil {
			fmt.Println("Ошибка записи лога:", err)
		}
		if verbose {
			fmt.Println(text)
		}
		f.Close()
	}
}

func disconnect(logDir string, server string) {
	if server == "" {
		return
	}
	filename := filepath.Join(logDir, server+".log")
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка открытия файла лога:", err)
		return
	}
	if _, err := f.WriteString("Connection closed\n"); err != nil {
		fmt.Println("Ошибка записи лога:", err)
	}
}
