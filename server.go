package main

/*
import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\\n", message)

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// TCP: payment
func handleConnection(conn net.Conn) {
	addr := conn.RemoteAddr().String()

	fmt.Println(addr)

	input := bufio.NewScanner(conn)

	for input.Scan() {
		fmt.Println(input.Text())
	}
}

func startTCPServer() {
	ln, err := net.Listen("tcp", ":80")

	if err != nil {
		fmt.Println("TCP Server error: ", err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("TCP error: ", err)
		}

		go handleConnection(conn)
	}
}

// UDP: notification
func startUDPServer() {
	udpLn, err := net.ListenPacket("udp", ":4040")

	if err != nil {
		fmt.Println("UDP server error: ", err)
	}

	for {
		buf := make([]byte, 1024)
		_, addr, err := udpLn.ReadFrom(buf)

		if err != nil {
			continue
		}

		go response(udpLn, addr, buf)
	}
}

func response(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	time := time.Now().Format(time.ANSIC)
	response := fmt.Sprintf("time received: %v. Your message: %v", time, string(buf))

	udpServer.WriteTo([]byte(response), addr)
}

// Websocket: chat
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		mType, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Reader error: ", err)
			return
		}

		fmt.Println(string(message))

		if err := conn.WriteMessage(mType, message); err != nil {
			fmt.Println("Writer error: ", err)
			return
		}
	}
}

func startWebsocketServer() {
	http.HandleFunc("/ws/", handleWebsocket)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error upgrading: ", err)
		return
	}

	fmt.Println("Client has connected!")

	err = ws.WriteMessage(1, []byte("Hi Client!"))

	if err != nil {
		fmt.Println("Error writing: ", err)
		return
	}

	reader(ws)
}

func main() {
	go startTCPServer()
	go startUDPServer()
	startWebsocketServer()
}
*/
