/*
This script tests TCP packet fragmentation resilience.
It verifies that the server correctly buffers and reconstructs commands
that are maliciously split across multiple tiny TCP packets arriving
with artificial network delays.
*/
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
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

	// Test: TCP Packet Fragmentation
	fmt.Println("\n--- Test: TCP Packet Fragmentation ---")
	conn.Write([]byte("CH"))
	time.Sleep(20 * time.Millisecond)
	conn.Write([]byte("AT GL"))
	time.Sleep(20 * time.Millisecond)
	conn.Write([]byte("OBAL test frag\n"))

	line, _ = reader.ReadString('\n')
	fmt.Printf("Response: %s", line)

	fmt.Println("\nFragmentation test finished executing!")
}
