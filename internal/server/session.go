package server


import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"the_answer_protocol/internal/protocol"
	"the_answer_protocol/internal/world"
)

type Client struct {
	hub  *Hub
	conn net.Conn
	send chan []byte

	currentRoomID string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		text := scanner.Text()

		cmd := protocol.Parse(text)

		switch cmd.Action {
		case "LOOK":
			c.handleLook()

		case "MOVE":
			if len(cmd.Args) > 0 {
				c.handleMove(cmd.Args[0])
			} else {
				c.send <- []byte("Move where? (e.g., MOVE north)\n")
			}

		case "CHAT":
			message := fmt.Appendf(
				nil,
				"[%s]: %s\n",
				c.conn.RemoteAddr(),
				text,
			)
			c.hub.broadcast <- message

		case "UNKNOWN":
			c.send <- []byte(
				"Unknown command. Try LOOK, MOVE <direction>, or CHAT <message>\n",
			)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf(
			"Read error for client %s: %v\n",
			c.conn.RemoteAddr(),
			err,
		)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		_, err := c.conn.Write(message)
		if err != nil {
			fmt.Printf(
				"Write error for client %s: %v\n",
				c.conn.RemoteAddr(),
				err,
			)
			return
		}
	}
}

func ServeClient(hub *Hub, conn net.Conn) {
	client := &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan []byte, 256),
		currentRoomID: "town_square",
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) handleLook() {
	var currentRoom *world.Room
	for _, r := range c.hub.worldMap.Rooms {
		if r.ID == c.currentRoomID {
			currentRoom = &r
			break
		}
	}

	if currentRoom == nil {
		c.send <- []byte("You are floating in the void...\n")
		return
	}

	response := fmt.Sprintf(
		"\n--- %s ---\n%s\n",
		currentRoom.Name,
		currentRoom.Description,
	)

	if len(currentRoom.Exits) > 0 {
		response += "Exits: "
		for direction := range currentRoom.Exits {
			response += direction + " "
		}
		response += "\n"
	}

	c.send <- []byte(response)
}


func (c *Client) handleMove(direction string) {
	direction = strings.ToLower(strings.TrimSpace(direction))

	var currentRoom *world.Room
	for _, r := range c.hub.worldMap.Rooms {
		if r.ID == c.currentRoomID {
			currentRoom = &r
			break
		}
	}

	if currentRoom == nil {
		c.send <- []byte("Error: Current room not found.\n")
		return
	}

	nextRoomID, exists := currentRoom.Exits[direction]
	if !exists {
		c.send <- []byte("You cannot go that way.\n")
		return
	}

	c.currentRoomID = nextRoomID
	c.send <- fmt.Appendf(nil, "You move %s.\n", direction)

	c.handleLook()
}
