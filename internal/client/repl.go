package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/peterh/liner"
)

func (c *Client) RunREPL() {
	cli := liner.NewLiner()
	defer cli.Close()
	cli.SetCtrlCAborts(true)

	for {
		input, err := cli.Prompt("Enter command (help, quit etc):")
		if err != nil {
			if err == liner.ErrPromptAborted {
				fmt.Println("\nAborted")
				break
			}
			fmt.Println("Error reading line:", err)
			break
		}
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "" {
			continue
		}
		cli.AppendHistory(input)

		cmd, arg, _ := strings.Cut(input, " ")
		switch cmd {
		case "quit", "q":
			return
		case "h", "help":
			fmt.Println(`Available commands:
	MEM all|phy|vir     -- display memory info(def:ALL)
	GPU                 -- display GPU info
	HIDE ms             -- hide window for n miliseconds
	c|connect serv addr -- make new connection (server_number address(localhost))
	l|list              -- list active connections
	clear               -- close all connections
	d|disconnect n      -- close one connection(server idx from list)
	q|quit              -- exit this program
	h|help              -- display this text`)
		case "connect", "c":
			parts := strings.Fields(arg)
			if len(parts) == 0 {
				fmt.Println("Usage: connect <server> <addr>")
				continue
			}
			srv := parts[0]
			addr := "localhost"
			if len(parts) > 1 {
				addr = parts[1]
			}
			serv_num, err := strconv.Atoi(srv)
			if err != nil {
				fmt.Println("Correct server number should be given")
				continue
			}
			c.Connect(serv_num, addr)
		case "list", "l":
			c.List()
		case "disconnect", "d":
			var serv_num int
			var err error
			if c.Len() == 0 {
				fmt.Println("No connections to disconnect")
				continue
			} else if arg == "" && c.Len() == 1 {
				serv_num = 0
			} else {
				serv_num, err = strconv.Atoi(arg)
				if err != nil {
					fmt.Println("Correct server number should be given")
					continue
				}
			}
			err = c.Disconnect(serv_num)
			if err != nil {
				fmt.Println("Error while closing connection:", err)
			}
		case "clear":
			err := c.DisconnectAll()
			if err != nil {
				fmt.Println("Error while closing connections:", err)
			}
		default:
			c.Send(cmd, arg)
		}
	}
}
