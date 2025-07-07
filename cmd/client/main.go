package main

import (
	"flag"
	"fmt"
	"runtime"

	"client-server-app/internal/client"
)

func main() {
	server := flag.Int("server", 0, "Server number to connect (1 or 2)")
	addr := flag.String("addr", "", "Server address")
	flag.Parse()

	cli := client.New(5) // MAX_CONNECTIONS
	if *server != 0 && *addr != "" {
		if err := cli.Connect(*server, *addr); err != nil {
			fmt.Println("Ошибка подключения:", err)
			return
		}
	}

	fmt.Printf("Connected: %d\nArch: %s\nOS: %s\n", cli.Len(), runtime.GOARCH, runtime.GOOS)
	cli.RunREPL()
}
