package server

import (
	"bufio"
	"net"
	"time"

	"the_answer_protocol/tap/src/internal/protocol"
)

type Client struct {
	hub  *Hub
	conn net.Conn
	send chan []byte

	remoteAddr    string
	username      string
	authenticated bool
	currentRoomID string
	groupName     string
	isInvited     []string
	hp            int
	inventory     []string
	quests        []string

	// recentCmds tracks this client's own recent command timestamps, for
	// the command-flooding abuse signal. Only readPump's own goroutine
	// touches it, so no synchronization is needed.
	recentCmds []time.Time
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
			c.hub.BroadcastGlobal([]byte(protocol.FormatEvt("GLOBAL", "PRESENCE LEAVE", c.username)))
		}
		c.hub.logger.Info("client_disconnected",
			"username", c.username,
			"remote_addr", c.remoteAddr,
		)
		c.hub.unregister <- c
		c.conn.Close()
	}()

	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		text := scanner.Text()
		cmd := protocol.Parse(text)

		c.hub.logger.Info("command_received",
			"username", c.username,
			"remote_addr", c.remoteAddr,
			"action", cmd.Action,
			"args", cmd.Args,
		)
		c.checkFlood()

		if cmd.Action != "CONNECT" && cmd.Action != "QUIT" && !c.authenticated {
			c.reply([]byte(protocol.FormatErr(protocol.ErrNotAuthenticated, "NOT_AUTHENTICATED")))
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
		case "TAKE":
			c.handleTake(cmd.Args)
		case "DROP":
			c.handleDrop(cmd.Args)
		case "INVENTORY":
			c.handleInventory()
		case "TALK":
			c.handleTalk(cmd.Args)
		case "ATTACK":
			c.handleAttack(cmd.Args)
		case "STATUS":
			c.handleStatus()
		case "QUEST":
			c.handleQuest(cmd.Args)
		case "QUESTS":
			c.handleQuests()
		case "UNKNOWN":
			c.reply([]byte(protocol.FormatErr(protocol.ErrUnknownCommand, "UNKNOWN_COMMAND")))
		default:
			c.reply([]byte(protocol.FormatErr(protocol.ErrUnknownCommand, "UNKNOWN_COMMAND")))
		}
	}

	if err := scanner.Err(); err != nil {
		c.hub.logger.Error("read_error",
			"username", c.username,
			"remote_addr", c.remoteAddr,
			"error", err.Error(),
		)
	}
}

// checkFlood logs a warning if this client is sending commands faster than
// a reasonable human-driven rate.
func (c *Client) checkFlood() {
	const (
		window    = 1 * time.Second
		threshold = 10
	)
	now := time.Now()
	cutoff := now.Add(-window)

	kept := c.recentCmds[:0]
	for _, t := range c.recentCmds {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	c.recentCmds = append(kept, now)

	if len(c.recentCmds) > threshold {
		c.hub.logger.Warn("possible_abuse",
			"type", "command_flooding",
			"username", c.username,
			"remote_addr", c.remoteAddr,
			"commands_in_window", len(c.recentCmds),
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
			c.hub.logger.Error("write_error",
				"username", c.username,
				"remote_addr", c.remoteAddr,
				"error", err.Error(),
			)
			return
		}
	}
}

func ServeClient(hub *Hub, conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()

	client := &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan []byte, 256),
		remoteAddr:    remoteAddr,
		currentRoomID: hub.worldMap.StartRoomID,
		isInvited:     []string{},
		hp:            100,
		inventory:     []string{},
		quests:        []string{},
	}

	client.hub.register <- client

	hub.logger.Info("client_connected", "remote_addr", remoteAddr)

	ip, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		if attempts := hub.RecordConnection(ip); attempts > 5 {
			hub.logger.Warn("possible_abuse",
				"type", "rapid_connections",
				"remote_addr", remoteAddr,
				"attempts_in_window", attempts,
			)
		}
	}

	client.reply([]byte(protocol.FormatOK("hello proto=1")))

	go client.writePump()
	go client.readPump()
}
