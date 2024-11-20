package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
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

		res := string(buf[0:n])
		reqType := string(buf[0:5])

		switch reqType {
		case "check":
			clientAddr := res[6:]
			var viewer User = viewers[clientAddr]

			fmt.Println("u:", clientAddr, "v:", viewer.username, "b:", viewer.balance)

			str := "Your balance is " + strconv.Itoa(viewer.balance) + "."
			ln.WriteToUDP([]byte(str), conn)
		case "topup":
			str := res[6:]
			strSlice := strings.Split(str, ":")
			clientAddr := strSlice[0] + ":" + strSlice[1]
			amountStr := strSlice[2]
			amount, err := strconv.Atoi(amountStr)

			var str2 string
			var viewer User = viewers[clientAddr]

			if err != nil {
				str2 = "Amount invalid. Please try again."
			} else {
				viewer.balance += amount

				viewers[clientAddr] = viewer
				str2 = "Added " + amountStr + " to your balance."
			}

			fmt.Println("u:", clientAddr, "a:", amount, "v:", viewer.username, "b:", viewer.balance)
			ln.WriteToUDP([]byte(str2), conn)
		}

		fmt.Println(res, reqType, conn)
	}
}

// Websocket Server
var upgrader = websocket.Upgrader{}

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Websocket error:", err)
		return
	}

	defer c.Close()
}

func startWebsocketServer() {

}

func main() {
	go startTCPServer()
	go startUDPServer()
	startWebsocketServer()
}
