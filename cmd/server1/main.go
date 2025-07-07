package main

import (
	"flag"

	"client-server-app/internal/commands"
	"client-server-app/internal/server"
)

func main() {
	var v bool
	flag.BoolVar(&v, "v", false, "")
	flag.Parse()

	disp := server.NewDispatcher(commands.GPU{}, commands.Hide{})
	srv := server.NewServer(disp,
		server.WithAddr(":8001"),
		server.WithMaxClients(105),
		server.WithLogger("localhost:9000", "Server1", v),
	)
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
