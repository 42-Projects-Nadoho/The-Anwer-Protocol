package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"the_answer_protocol/internal/protocol"
)

func main() {
	fmt.Println("=== Welcome to TAP ===")
	fmt.Println("Connecting to localhost:8080...")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection Error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected! You can now type your commands.")

	readerDone := make(chan struct{})

	go func() {
		defer close(readerDone)
		serverScanner := bufio.NewScanner(conn)
		for serverScanner.Scan() {
			fmt.Printf("\r%s\n> ", serverScanner.Text())
		}
		if err := serverScanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from server: %v\n", err)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		rawText := scanner.Text()
		cmd := protocol.Parse(rawText)

		if cmd.Action == "UNKNOWN" {
			continue
		}

		fmt.Fprintf(conn, "%s\n", rawText)

		if cmd.Action == "QUIT" {
			// Server closes the connection after replying; wait for that.
			fmt.Println("Goodbye!")
			<-readerDone
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
