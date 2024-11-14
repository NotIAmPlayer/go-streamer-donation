package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var username string

func connectTCPServer() {
	conn, err := net.Dial("tcp", ":80")

	if err != nil {
		fmt.Println("Error: ", err)
	}

	conn.Write([]byte("client:" + username))

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)
		fmt.Println("done")
		done <- struct{}{}
	}()

	copyInput(conn, os.Stdin)

	conn.Close()
	<-done
}

func copyInput(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	fmt.Print("Enter your username: ")
	fmt.Scan(&username)

	fmt.Println("Welcome, " + username + "!")

	go connectTCPServer()
}
