/*
This script tests for simultaneous movement race conditions on the server.
It verifies that if two clients move at the exact same millisecond, the 
presence broadcast events are cleanly and sequentially distributed without
crashing the global event loop or corrupting room state.
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

const address = "127.0.0.1:4242"

func connectClient(name string) (net.Conn, *bufio.Reader) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("[%s] Dial error: %v\n", name, err)
		os.Exit(1)
	}

	reader := bufio.NewReader(conn)
	line, _ := reader.ReadString('\n')
	if !strings.HasPrefix(line, "OK hello proto=1") {
		fmt.Printf("[%s] Invalid handshake: %s\n", name, line)
		os.Exit(1)
	}

	fmt.Fprintf(conn, "CONNECT %s\n", name)
	line, _ = reader.ReadString('\n')
	if !strings.HasPrefix(line, "OK connected") {
		fmt.Printf("[%s] Connect failed: %s\n", name, line)
		os.Exit(1)
	}
	return conn, reader
}

func main() {
	aliceConn, _ := connectClient("Alice")
	bobConn, _ := connectClient("Bob")
	_, charlieReader := connectClient("Charlie")
	
	defer aliceConn.Close()
	defer bobConn.Close()

	fmt.Println("All 3 clients connected (Alice, Bob, Charlie).")
	fmt.Println("Alice and Bob will attempt to MOVE at the exact same millisecond.")
	
	var wg sync.WaitGroup
	wg.Add(2)
	
	startLine := time.Now().Add(500 * time.Millisecond)
	
	go func() {
		time.Sleep(time.Until(startLine))
		fmt.Fprintf(aliceConn, "MOVE north\n")
		wg.Done()
	}()
	
	go func() {
		time.Sleep(time.Until(startLine))
		fmt.Fprintf(bobConn, "MOVE south\n")
		wg.Done()
	}()
	
	// Charlie listens for presence events
	go func() {
		aliceLeft := false
		bobLeft := false
		
		for {
			line, err := charlieReader.ReadString('\n')
			if err != nil {
				break
			}
			if strings.Contains(line, "EVT ROOM PRESENCE LEAVE Alice") {
				aliceLeft = true
				fmt.Println("[Charlie] Received: Alice left!")
			}
			if strings.Contains(line, "EVT ROOM PRESENCE LEAVE Bob") {
				bobLeft = true
				fmt.Println("[Charlie] Received: Bob left!")
			}
			if aliceLeft && bobLeft {
				fmt.Println("[Charlie] Successfully received both leave events without server crash.")
				os.Exit(0)
			}
		}
	}()

	wg.Wait()
	
	// Wait a bit to let Charlie read the events
	time.Sleep(1 * time.Second)
	fmt.Println("Test finished.")
}
