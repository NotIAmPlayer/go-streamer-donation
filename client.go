package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

var username string

func connectTCPServer() {
	conn, err := net.Dial("tcp", ":80")

	if err != nil {
		fmt.Println("Error: ", err)
	}

	conn.Write([]byte("client:" + username + "\n"))

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)
		//fmt.Println("done")
		done <- struct{}{}
	}()

	var input string

	for {
		fmt.Println("SELECT AN OPTION")
		fmt.Println("=================")
		fmt.Println("1. Check Balance")
		fmt.Println("2. Top-Up Balance")
		fmt.Println("3. Send Donation")
		fmt.Println("0. EXIT")
		fmt.Print("> ")

		fmt.Scan(&input)

		switch input {
		case "0":
			break
		case "1": // UDP
			fmt.Println("Checking your balance...")
			checkBalance(conn.LocalAddr().String())
		case "2": // UDP
			var amountStr string
			var amountInt int

			for {
				fmt.Print("Input the amount you want to top-up: ")
				fmt.Scan(&amountStr)

				temp, err := strconv.Atoi(amountStr)

				if err != nil {
					fmt.Println("Amount invalid. Please try again.")
				} else if temp < 0 {
					fmt.Println("Amount must be a positive value. Please try again.")
				} else {
					amountInt = temp
					break
				}
			}

			topUpBalance(conn.LocalAddr().String(), amountInt)
		case "3": //TCP
			fmt.Println("Preparing sending donation...")
		default:
			fmt.Println("Input invalid. Please try again.")
		}

		if input == "0" {
			break
		}

		fmt.Println("")
	}

	conn.Close()
	<-done
}

func checkBalance(addr string) {
	connectUDPServer("check:" + addr)
}

func topUpBalance(addr string, amount int) {
	connectUDPServer("topup:" + addr + ":" + strconv.Itoa(amount))
}

func connectUDPServer(str string) {
	addr, err := net.ResolveUDPAddr("udp", ":4080")

	if err != nil {
		fmt.Println("Client - Address Resolving Error:", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)

	defer conn.Close()

	conn.Write([]byte(str))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Client - Error Reading Response:", err)
	}

	fmt.Println(string(buf[0:n]))
}

func main() {
	fmt.Print("Enter your username: ")
	fmt.Scan(&username)

	fmt.Println("Welcome, " + username + "!")

	connectTCPServer()
}
