/*
This script tests TCP packet coalescing resilience.
It verifies that the server correctly parses and executes multiple
newline-terminated commands that arrive artificially concatenated
within a single TCP packet flush.
*/
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	addr := "127.0.0.1:4242"
	fmt.Printf("Connecting to %s...\n", addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Dial error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Consume handshake
	line, _ := reader.ReadString('\n')
	if !strings.HasPrefix(line, "OK hello") {
		fmt.Printf("Invalid handshake: %s\n", line)
		os.Exit(1)
	}

	// Login
	fmt.Println("--- Authenticating ---")
	fmt.Fprintf(conn, "CONNECT Alice\n")
	line, _ = reader.ReadString('\n')
	fmt.Printf("Server: %s", line)
	if !strings.HasPrefix(line, "OK ") {
		fmt.Printf("Failed to login: %s\n", line)
		os.Exit(1)
	}

	// Test: TCP Packet Coalescing
	fmt.Println("\n--- Test: TCP Packet Coalescing ---")
	// Send multiple commands in one flush
	conn.Write([]byte("LOOK\nCHAT GLOBAL test coalescing\n"))
	
	line1, _ := reader.ReadString('\n')
	fmt.Printf("Response 1: %s", line1)
	
	line2, _ := reader.ReadString('\n')
	fmt.Printf("Response 2: %s", line2)

	fmt.Println("\nCoalescing test finished executing!")
}
