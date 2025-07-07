package main

import (
	"flag"
	"fmt"
	"net"

	"client-server-app/pkg/logger"
)

func main() {
	verbose := flag.Bool("v", false, "включить подробный вывод")
	flag.Parse()
	logDir, err := logger.СreateSessionLogDir()
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
		go logger.HandleConnection(conn, logDir, *verbose)
	}
}
