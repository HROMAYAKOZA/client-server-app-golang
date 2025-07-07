package client

import (
	"fmt"
	"strings"
)

func (c *Client) Send(cmd, arg string) {
	msg := cmd
	if arg != "" {
		msg += ":" + arg
	}
	msg += "\n"

	if len(c.conns) == 0 {
		fmt.Println("No connections")
		return
	}

	if len(c.conns) == 1 {
		c.conns[0].Conn.Write([]byte(msg))
		resp, err := c.conns[0].Reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("EOF ERROR, disconnecting from server")
				c.Disconnect(0)
			} else {
				fmt.Println("Error reading response:", err)
			}
			return
		}
		fmt.Print("Response: ", resp)
	} else {
		fl := true
		for i, con := range c.conns {
			con.Conn.Write([]byte(msg))
			resp, err := con.Reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Printf("EOF ERROR, disconnecting from server #%d\n", i)
					c.Disconnect(i)
				} else {
					fmt.Println("Error reading response:", err)
				}
				continue
			}
			endOfTimeIdx := strings.Index(resp, "]")
			if endOfTimeIdx == -1 {
				fmt.Println("Incorrect answer format")
				return
			}
			message := strings.TrimSpace(resp[endOfTimeIdx+1:])
			if strings.HasPrefix(message, "UNKNOWN CMD") {
				continue
			}
			fl = false
			fmt.Print("Response: ", resp)
		}
		if fl {
			fmt.Println("UNKNOWN CMD")
		}
	}
}
