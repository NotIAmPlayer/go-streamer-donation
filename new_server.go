package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
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

		fmt.Println(txt, len(txt))

		if len(txt) >= 7 {
			str := strings.Split(txt, ":")
			reqType := str[0]

			fmt.Println(str)

			if reqType == "client" { // Client-sent usernames
				username := str[1]
				viewers[addr] = User{username: username, balance: 0}
			} else if reqType == "donation" {
				clientAddr := str[1] + ":" + str[2]
				viewer := viewers[clientAddr]
				streamer := str[3]
				amountStr := str[4]
				message := str[5]

				amountInt, err := strconv.Atoi(amountStr)
				var str2 string

				if err != nil {
					str2 = "Amount invalid. Please try again."
				} else if amountInt > viewer.balance {
					str2 = "Donation amount is higher than account balance. Please try again."
				} else {
					if conn, exists := streamers[streamer]; exists {
						viewer.balance -= amountInt
						str2 = "Donated " + amountStr + " to " + streamer + "."
						viewers[clientAddr] = viewer

						broadcastToWebsocket(conn, viewer.username, amountInt, streamer, message)
					} else {
						str2 = "Streamer not found."
					}
				}
				conn.Write([]byte(str2 + "\n"))
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
var streamers = make(map[string]*websocket.Conn)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Websocket error:", err)
		return
	}

	defer conn.Close()

	fmt.Println(conn)

	for {
		// Read from client
		_, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Websocket Server - Reading error:", err)
		}
		fmt.Println("Received:", message)

		var msg map[string]string
		json.Unmarshal(message, &msg)

		if msg["type"] == "streamer" {
			username := msg["username"]
			streamers[username] = conn
			fmt.Println("Registered streamer:", username)
		}
	}
}

func broadcastToWebsocket(conn *websocket.Conn, username string, amount int, streamer string, message string) {
	donationMessage := map[string]interface{}{
		"type":    "donation",
		"from":    username,
		"amount":  amount,
		"message": message,
	}

	msg, _ := json.Marshal(donationMessage)
	err := conn.WriteMessage(websocket.TextMessage, msg)

	if err != nil {
		fmt.Println("Websocket - Broadcast error:", err)
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	username = r.URL.Path[8:]
	//fmt.Fprintf(w, "Hello, %s! Ready to stream?", r.URL.Path[8:])
}

func startWebsocketServer() {
	http.HandleFunc("/ws/", websocketHandler)
	http.HandleFunc("/stream/", streamHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	go startTCPServer()
	go startUDPServer()
	startWebsocketServer()
}
