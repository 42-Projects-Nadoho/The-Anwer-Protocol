package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

/*
test_abuse.go

This script simulates abuse patterns and Denial of Service (DoS) attacks
against the server to verify the logging, monitoring, and stability systems
are functioning correctly.

Tests performed:
1. Rapid Connection Cycling (Port Exhaustion / Resource Starvation): 
   Repeatedly opening and immediately closing TCP connections. This tests 
   if the server properly handles early EOFs and prevents goroutine leaks 
   or race conditions in the connection lifecycle.

2. Command Flooding (Application Layer Spam): 
   Spamming hundreds of commands in a fraction of a second. This verifies 
   that the rate-limiting and abuse monitoring (checkFlood) logic successfully 
   flags malicious clients.
*/

func main() {
	addr := "127.0.0.1:4242"
	fmt.Printf("Connecting to %s...\n", addr)

	fmt.Println("--- Test: Rapid Connection Cycling ---")
	// Fire off 20 extremely rapid connections, abandoning them immediately
	// without waiting for the server's handshake. This is designed to trigger
	// race conditions during the server's ServeClient initialization phase.
	for i := 0; i < 20; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("Connection %d failed: %v\n", i, err)
			continue
		}
		// Instantly close to simulate connection spam/port exhaustion attempts
		conn.Close()
	}
	fmt.Println("Sent 20 rapid connections.")

	fmt.Println("\n--- Test: Command Flooding ---")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Authenticate the spammer client
	fmt.Fprintf(conn, "CONNECT abuser\n")
	time.Sleep(100 * time.Millisecond) // Wait for server to process auth

	// Flood the server with 200 LOOK commands instantly.
	// This will trigger the checkFlood() window threshold.
	for i := 0; i < 200; i++ {
		fmt.Fprintf(conn, "LOOK\n")
	}

	fmt.Println("Fired 200 commands instantly.")
	fmt.Println("Check the server logs to verify 'possible_abuse' WARN logs were triggered!")
	
	// Wait a moment so the server processes before we drop connection
	time.Sleep(500 * time.Millisecond)
	fmt.Println("\nAbuse test finished executing!")
}
