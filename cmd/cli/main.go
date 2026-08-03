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

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		rawText := scanner.Text()
		cmd := protocol.Parse(rawText)

		if cmd.Action == "QUIT" {
			fmt.Println("Goodbye!")
			break
		}

		if cmd.Action == "UNKNOWN" {
			continue
		}

		fmt.Fprintf(conn, "%s\n", rawText)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
