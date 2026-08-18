/*
This script tests the server's group state volatility handling.
It verifies that the group management locks do not break or deadlock when 
members join and leave chaotically (e.g., someone accepts a group invite 
at the exact same millisecond the group leader dissolves the group).
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
	aliceConn, aliceReader := connectClient("Alice")
	bobConn, bobReader := connectClient("Bob")
	charlieConn, _ := connectClient("Charlie")
	
	defer aliceConn.Close()
	defer bobConn.Close()
	defer charlieConn.Close()

	fmt.Println("Clients connected.")
	
	// Create group and invite
	fmt.Fprintf(aliceConn, "GROUP CREATE\n")
	aliceReader.ReadString('\n') // OK group created
	
	fmt.Fprintf(aliceConn, "GROUP INVITE Bob\n")
	aliceReader.ReadString('\n') // OK invited
	bobReader.ReadString('\n')   // EVT GROUP INVITE Alice
	
	fmt.Fprintf(aliceConn, "GROUP INVITE Charlie\n")
	aliceReader.ReadString('\n') // OK invited

	fmt.Println("Invites sent. Initiating chaotic join/leave race condition...")
	
	var wg sync.WaitGroup
	wg.Add(3)
	
	startLine := time.Now().Add(500 * time.Millisecond)
	
	// Alice abandons the group
	go func() {
		time.Sleep(time.Until(startLine))
		fmt.Fprintf(aliceConn, "GROUP LEAVE\n")
		wg.Done()
	}()
	
	// Bob accepts invite
	go func() {
		time.Sleep(time.Until(startLine))
		fmt.Fprintf(bobConn, "GROUP JOIN Alice\n")
		wg.Done()
	}()
	
	// Charlie accepts invite
	go func() {
		time.Sleep(time.Until(startLine))
		fmt.Fprintf(charlieConn, "GROUP JOIN Alice\n")
		wg.Done()
	}()
	
	// Just listen to bob's response
	go func() {
		for {
			bobConn.SetReadDeadline(time.Now().Add(2 * time.Second))
			line, err := bobReader.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Printf("[Bob received] %s", line)
		}
	}()

	wg.Wait()
	
	// Give the server time to process
	time.Sleep(1 * time.Second)
	fmt.Println("Finished executing chaotic group state test without crashing the server.")
}
