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

				fmt.Println(len(str))

				for i := 6; i < len(str); i++ {
					message += " " + str[i]
				}

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
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}
var streamers = make(map[string]*websocket.Conn)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Websocket error:", err)
		return
	}

	defer conn.Close()

	//fmt.Println(conn)

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
			if _, exists := streamers[username]; exists {
				conn.WriteMessage(websocket.TextMessage, []byte("Username already taken."))
				continue
			}

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
	username := r.URL.Path[8:]
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>` + username + ` | Streaming</title>
		<script src="https://cdn.tailwindcss.com"></script>
	</head>
	<body>
		<div class="my-4 mx-4 flex">
			<div class="w-8/12">
				<div class="h-[32rem] mb-2 rounded-2xl bg-black text-white text-center place-content-center">
					<h1 class="text-2xl font-bold">You're streaming!</h1>
					<p>Pretend this is a screen.</p>
				</div>
				<div>
					<h1 class="text-2xl font-bold">` + username + `'s stream</h2>
				</div>
			</div>
			<div class="w-4/12 h-[32rem] ml-4">
				<div class="h-10 bg-slate-300 rounded-t-2xl px-4 py-2 font-bold">Chat donations</div>
				<div class="h-[29.5rem] bg-slate-100 rounded-b-2xl px-4 py-2 overflow-y-auto" id="donation-container">
				</div>
			</div>
		</div>

		<script>
			const socket = new WebSocket("ws://localhost:8080/ws/");

			socket.onmessage = function (e) {
				const data = JSON.parse(e.data);
				const donationContainer = document.getElementById("donation-container");

				console.log(data);

				if (isNaN(data.amount)) {
					console.error("Invalid donation:", data.amount);
				}

				data.message = data.message.replaceAll("_", " ");
				
				const donation = document.createElement("div");
				donation.className = "bg-rose-200 px-2 py-2 rounded-2xl mb-2";
				donation.innerHTML = '<div class="flex items-center mb-1"><div class="w-3/4 flex items-center"><span class="h-8 w-8 rounded-full bg-rose-600 py-0.5 text-center text-white">'+data.from.charAt(0).toUpperCase()+'</span><p class="ml-2 w-auto">'+data.from+'</p></div><div class="w-1/4"><p class="text-right">Rp'+data.amount+'</p></div></div><hr class="border-rose-600" /><p>'+data.message+'</p>'
			
				donationContainer.appendChild(donation);
			}

			socket.onopen = () => {
				const regMessage = {
					type: "streamer",
					username: "` + username + `"
				};

				socket.send(JSON.stringify(regMessage));
			}

			socket.onerror = (error) => {
				console.error("WebSocket error:", error);
			};

			window.onbeforeunload = () => {
				socket.close();
			};
		</script>
	</body>
	</html>
	`
	fmt.Fprintf(w, html)
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
