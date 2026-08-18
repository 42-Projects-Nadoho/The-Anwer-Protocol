package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"the_answer_protocol/internal/protocol"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "Server address to connect to")
	flag.Parse()

	fmt.Println("=== Welcome to TAP ===")
	fmt.Printf("Connecting to %s...\n", *addr)

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR 900 CONNECTION_FAILED\n")
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected! You can now type your commands.")

	readerDone := make(chan struct{})
	quitting := false

	go func() {
		defer close(readerDone)
		serverScanner := bufio.NewScanner(conn)
		for serverScanner.Scan() {
			fmt.Printf("\r%s\n> ", serverScanner.Text())
		}
		if !quitting {
			fmt.Fprintf(os.Stderr, "\rERR 900 CONNECTION_FAILED\n")
			os.Exit(1)
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

		_, err := fmt.Fprintf(conn, "%s\n", rawText)
		if err != nil && !quitting {
			fmt.Fprintf(os.Stderr, "\rERR 900 CONNECTION_FAILED\n")
			os.Exit(1)
		}

		if cmd.Action == "QUIT" {
			quitting = true
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
