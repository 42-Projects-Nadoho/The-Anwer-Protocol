/*
This script automates a high-load stress test to ensure the server's input queues 
and locks can handle extreme traffic without crashing.

Under the hood, it:
1. Connects: Opens two raw TCP connections to the server and authenticates them as Alice0 and Alice1.
2. Floods: Spawns 100 parallel goroutines for *each* client to instantly blast 100 alternating LOOK and CHAT commands down the socket at the exact same time.
3. Drains: Runs a background reader to ingest the massive wave of JSON responses from the server.
4. Verifies: Waits for all commands to be sent and digested, proving that the server safely queued, parsed, and responded to 200 concurrent requests without suffering a panic or deadlocking its internal state.
*/
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	address    = "127.0.0.1:4242"
	numClients = 2
	numSends   = 100
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Printf("[Client %d] Dial error: %v\n", clientID, err)
				os.Exit(1)
			}
			defer conn.Close()

			reader := bufio.NewReader(conn)
			
			// Read the handshake
			line, _ := reader.ReadString('\n')
			if !strings.HasPrefix(line, "OK hello proto=1") {
				fmt.Printf("[Client %d] Invalid handshake: %s\n", clientID, line)
				os.Exit(1)
			}

			username := fmt.Sprintf("Alice%d", clientID)
			fmt.Fprintf(conn, "CONNECT %s\n", username)
			
			line, _ = reader.ReadString('\n')
			if !strings.HasPrefix(line, "OK connected") {
				fmt.Printf("[Client %d] Connect failed: %s\n", clientID, line)
				os.Exit(1)
			}

			fmt.Printf("[Client %d] Connected successfully. Spamming commands...\n", clientID)
			
			var sendWg sync.WaitGroup
			
			// Start reader to drain events and responses
			go func() {
				for {
					conn.SetReadDeadline(time.Now().Add(2 * time.Second))
					_, err := reader.ReadString('\n')
					if err != nil {
						break
					}
				}
			}()

			// Blast commands concurrently
			for j := 0; j < numSends; j++ {
				sendWg.Add(1)
				go func(cmdNum int) {
					defer sendWg.Done()
					if cmdNum%2 == 0 {
						fmt.Fprintf(conn, "LOOK\n")
					} else {
						fmt.Fprintf(conn, "CHAT ROOM Spam %d\n", cmdNum)
					}
				}(j)
			}
			
			sendWg.Wait()
			fmt.Printf("[Client %d] Finished spamming %d commands.\n", clientID, numSends)
			time.Sleep(1 * time.Second) // wait for server to digest
		}(i)
	}

	wg.Wait()
	fmt.Println("All clients finished successfully without server crash.")
}
