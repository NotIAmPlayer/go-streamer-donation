package main

import (
	"bufio"
	"fmt"
	"net"
)

var viewers = make(map[string]User)

type User struct {
	username string
	balance  int
}

func startTCPServer() {
	ln, err := net.Listen("tcp", ":80")

	if err != nil {
		fmt.Println("TCP Server - Error:", err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("TCP Server - Error:", err)
		}

		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	input := bufio.NewScanner(conn)

	for input.Scan() {
		txt := input.Text()

		fmt.Println(txt)

		if len(txt) >= 7 {
			if txt[0:7] == "client" {
				viewers[addr] = User{username: txt[8:], balance: 0}
			}
		}
	}
}

func main() {
	startTCPServer()
	//go startUDPServer()
	//startWebsocketServer()
}
