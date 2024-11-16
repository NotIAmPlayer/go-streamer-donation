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

// TCP Server
func startTCPServer() {
	ln, err := net.Listen("tcp", ":80")

	if err != nil {
		fmt.Println("TCP Server - Listening Error:", err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("TCP Server - Connection Error:", err)
		}

		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	input := bufio.NewScanner(conn)

	fmt.Println(addr)

	for input.Scan() {
		txt := input.Text()

		//fmt.Println(txt, len(txt))

		if len(txt) >= 7 {
			var reqType = txt[0:6]
			var username = txt[7:]

			fmt.Println(reqType, username)

			if reqType == "client" { // Client-sent usernames
				viewers[addr] = User{username: username, balance: 0}
			}
		}
	}
}

// UDP Server
func startUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", ":4080")

	if err != nil {
		fmt.Println("UDP Server - Address Resolving Error:", err)
	}

	ln, err := net.ListenUDP("udp", addr)

	if err != nil {
		fmt.Println("UDP Server - Listening Error:", err)
	}

	defer ln.Close()

	for {
		buf := make([]byte, 1024)
		n, conn, err := ln.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("UDP Server - Reading Error:", err)
		}

		var res = string(buf[0:n])
		var reqType = string(buf[0:5])

		switch reqType {
		case "check":
			var username = res[6:]
			var viewer User

			for _, v := range viewers {
				fmt.Println(v.username, v.balance)

				if v.username == username {
					viewer = v
					break
				}
			}

			fmt.Println("u:", username, "v:", viewer.username, "b:", viewer.balance)
		}

		fmt.Println(res, reqType, conn)
	}
}

func main() {
	go startTCPServer()
	startUDPServer()
	//startWebsocketServer()
}
