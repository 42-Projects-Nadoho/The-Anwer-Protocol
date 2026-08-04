package main

import (
	"os"
	"fmt"
	"net"
	"the_answer_protocol/internal/server"
)


func main() {
	hub := server.NewHub()
	go hub.Run()

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
		server.ServeClient(hub, conn)
	}
}
