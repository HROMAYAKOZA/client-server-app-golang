package main

import (
	"flag"

	"github.com/HROMAYAKOZA/client-server-app-golang/internal/commands"
	"github.com/HROMAYAKOZA/client-server-app-golang/internal/server"
)

func main() {
	var v bool
	flag.BoolVar(&v, "v", false, "")
	flag.Parse()

	disp := server.NewDispatcher(commands.Mem{})
	srv := server.NewServer(disp,
		server.WithAddr(":8002"),
		server.WithMaxClients(5),
		server.WithLogger("localhost:9000", "Server2", v),
	)
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
