package server

import (
	"bufio"
	"fmt"
	"net"
)


type Client struct {
	hub *Hub

	conn net.Conn

	send chan []byte
}


func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		text := scanner.Text()

		message := fmt.Appendf(nil, "[%s]: %s\n", c.conn.RemoteAddr(), text)
		c.hub.broadcast <- message
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Read error for client %s: %v\n", c.conn.RemoteAddr(), err)
	}
}


func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		_, err := c.conn.Write(message)
		if err != nil {
			fmt.Printf("Write error for client %s: %v\n", c.conn.RemoteAddr(), err)
			return
		}
	}
}

func ServeClient(hub *Hub, conn net.Conn) {
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256)}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
