package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)


func main() {
	fmt.Println("Starting server on port 8080......")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lancement Error: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server waiting for player....")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Acceptation Error: %v\n", err)
			continue
		}

		fmt.Printf("New client connectd from : %s\n", conn.RemoteAddr())
		go manageClient(conn)
	}
}


func manageClient(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		texteRecu := scanner.Text()
		fmt.Printf("[Message arrived] : %s\n", texteRecu)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
	fmt.Println("The client is deconected.")
	conn.Close()
}
