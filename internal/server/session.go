package server

import (
	"bufio"
	"fmt"
	"net"
	"the_answer_protocol/internal/protocol"
)

type Client struct {
	hub  *Hub
	conn net.Conn
	send chan []byte

	username      string
	authenticated bool
	currentRoomID string
	groupName     string
	isInvited     []string
}

func (c *Client) readPump() {
	defer func() {
		if c.authenticated {
			if c.groupName != "" {
				c.hub.LeaveGroup(c)
			}
			c.hub.BroadcastRoom(
				c.currentRoomID,
				[]byte(protocol.FormatEvt("ROOM", "PRESENCE LEAVE", c.username)),
				c,
			)
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()

	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		text := scanner.Text()
		cmd := protocol.Parse(text)

		if cmd.Action != "CONNECT" && cmd.Action != "QUIT" && !c.authenticated {
			c.send <- []byte(protocol.FormatErr(protocol.ErrNotAuthenticated, "NOT_AUTHENTICATED"))
			continue
		}

		switch cmd.Action {
		case "CONNECT":
			c.handleConnect(cmd.Args)

		case "LOOK":
			c.handleLook()

		case "MOVE":
			c.handleMove(cmd.Args)

		case "CHAT":
			c.handleChat(cmd)

		case "WHO":
			c.handleWho()

		case "QUIT":
			c.handleQuit()
			return

		case "GROUP":
			c.handleGroup(cmd.Args)

		case "UNKNOWN":
			c.send <- []byte(protocol.FormatErr(protocol.ErrUnknownCommand, "UNKNOWN_COMMAND"))

		default:
			c.send <- []byte(protocol.FormatErr(protocol.ErrUnknownCommand, "UNKNOWN_COMMAND"))
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
		isInvited:     []string{},
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()

	client.send <- []byte(protocol.FormatOK("hello proto=1"))
}
