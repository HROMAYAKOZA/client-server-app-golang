package main

import (
	"bufio"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"sync"
)

func ddos(wg *sync.WaitGroup) {
	defer wg.Done()
	addr := "localhost"
	port := 8001
	conn, err := net.Dial("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		fmt.Println("connection error")
		return
	}
	defer conn.Close()

	srvReader := bufio.NewReader(conn)
	conn.Write([]byte("gpu\n"))

	resp, err := srvReader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Println("Response:", resp)
	select {}
}

func main() {
	fmt.Println("Running ddos on architecture:", runtime.GOARCH)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go ddos(&wg)
	}

	wg.Wait()
}
