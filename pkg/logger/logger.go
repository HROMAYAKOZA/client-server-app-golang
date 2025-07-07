package logger

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func СreateSessionLogDir() (string, error) {
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

func HandleConnection(conn net.Conn, logDir string, verbose bool) {
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
