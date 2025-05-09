package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"

	"github.com/peterh/liner"
)

type ServConnection struct {
	conn       net.Conn
	connReader *bufio.Reader
	server     int
}
type ConsList []ServConnection

func (con *ConsList) Append(c ServConnection) error {
	if len(*con) >= MAX_CONNECTIONS {
		return errors.New("CONNECTIONS LIMIT EXCEEDED")
	}
	*con = append(*con, c)
	return nil
}

func (con *ConsList) Len() int {
	return len(*con)
}

func (con *ConsList) List() {
	for i, c := range cons {
		fmt.Println("Server", c.server, "#", i)
	}
	if len(cons) == 0 {
		fmt.Println("No active connections")
	}
}

func (con *ConsList) RemoveByIndex(index int) error {
	if index < 0 || index >= len(*con) {
		return errors.New("WRONG INDEX")
	}
	err := (*con)[0].conn.Close()
	if err != nil {
		return err
	}
	*con = append((*con)[:index], (*con)[index+1:]...)
	return nil
}

func (con *ConsList) Clear() error {
	for i, c := range *con {
		err := c.conn.Close()
		if err != nil {
			return err
		}
		(*con)[i] = ServConnection{}
	}
	*con = (*con)[:0]
	return nil
}

const MAX_CONNECTIONS = 5

var cons ConsList

func new_con(server int, addr string) error {
	if server != 1 && server != 2 {
		return errors.New("SERVER NUMBER NOT SUPPORTED")
	}
	if len(cons) >= MAX_CONNECTIONS {
		return errors.New("CONNECTIONS LIMIT EXCEEDED")
	}
	port := 8000 + server
	conn, err := net.Dial("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		err_str := fmt.Sprint("Connection error:", err)
		return errors.New(err_str)
	}
	// defer conn.Close()
	fmt.Printf("Connected to server %d at %s:%d\n", server, addr, port)
	srvReader := bufio.NewReader(conn)
	err = cons.Append(ServConnection{conn, srvReader, server})
	if err != nil {
		err_str := fmt.Sprint("Connection error:", err)
		return errors.New(err_str)
	}
	return nil
}

func send_response(cmd string, arg string) {
	var mess string
	if arg != "" {
		mess = cmd + ":" + arg + "\n"
	} else {
		mess = cmd + "\n"
	}
	if len(cons) == 0 {
		fmt.Println("No active connections")
	} else if len(cons) == 1 {
		// Send command to server
		cons[0].conn.Write([]byte(mess))
		// Read response
		resp, err := cons[0].connReader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("EOF ERROR, disconnecting from server")
				cons.RemoveByIndex(0)
			} else {
				fmt.Println("Error reading response:", err)
			}
			return
		}
		fmt.Print("Response: ", resp)
	} else {
		fl := true
		for i, con := range cons {
			con.conn.Write([]byte(mess))
			resp, err := con.connReader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Printf("EOF ERROR, disconnecting from server #%d\n", i)
					cons.RemoveByIndex(i)
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

func main() {
	server := flag.Int("server", 0, "Server number to connect (1 or 2)")
	addr := flag.String("addr", "", "Server address")
	flag.Parse()
	if *server != 0 && *addr != "" {
		err := new_con(*server, *addr)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Current connections:", len(cons))
	fmt.Println("Running client on architecture:", runtime.GOARCH)

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	for {
		input, err := line.Prompt("Enter command (help, quit etc): ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				fmt.Println("\nAborted")
				break
			}
			fmt.Println("Error reading line:", err)
			break
		}
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)
		if input == "" {
			continue
		}
		line.AppendHistory(input)

		cmd, arg, _ := strings.Cut(input, " ")

		switch cmd {
		case "q", "quit":
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
		case "c", "connect":
			server, addr, fl := strings.Cut(arg, " ")
			if !fl {
				addr = "localhost"
			}
			serv_num, err := strconv.Atoi(server)
			if err != nil {
				fmt.Println("Correct server number should be given")
				continue
			}
			err = new_con(serv_num, addr)
			if err != nil {
				fmt.Println(err)
			}
		case "clear":
			err := cons.Clear()
			if err != nil {
				fmt.Println("Error while closing connections:", err)
			}
		case "list":
			cons.List()
		case "d", "disconnect":
			var serv_num int
			var err error
			if cons.Len() == 0 {
				fmt.Println("No connections to disconnect")
				continue
			} else if arg == "" && cons.Len() == 1 {
				serv_num = 0
			} else {
				serv_num, err = strconv.Atoi(arg)
				if err != nil {
					fmt.Println("Correct server number should be given")
					continue
				}
			}
			err = cons.RemoveByIndex(serv_num)
			if err != nil {
				fmt.Println("Error while closing connection:", err)
			}
		default:
			send_response(cmd, arg)
		}
	}
}
