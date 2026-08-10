package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"the_answer_protocol/internal/protocol"
)

func (c *Client) handleConnect(args []string) {
	if c.authenticated {
		c.send <- []byte(protocol.FormatErr(protocol.ErrConnectionFailed, "ALREADY_CONNECTED"))
		return
	}
	if len(args) == 0 {
		c.send <- []byte(protocol.FormatErr(protocol.ErrConnectionFailed, "USERNAME_REQUIRED"))
		return
	}

	username := args[0]
	if !c.hub.TryConnect(c, username) {
		c.send <- []byte(protocol.FormatErr(protocol.ErrNameInUse, "NAME_IN_USE"))
		return
	}

	c.send <- []byte(protocol.FormatOK("connected"))

	c.hub.BroadcastRoom(
		c.currentRoomID,
		[]byte(protocol.FormatEvt("ROOM", "PRESENCE ENTER", c.username)),
		c,
	)
}

// Items/NPCs stay empty until that system exists (spacotto's lot).
type lookRoom struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
}

type lookResponse struct {
	Room    lookRoom `json:"room"`
	Players []string `json:"players"`
	Items   []string `json:"items"`
	NPCs    []string `json:"npcs"`
}

func (c *Client) handleLook() {
	room, ok := c.hub.worldMap.RoomByID(c.currentRoomID)
	if !ok {
		c.send <- []byte(protocol.FormatErr(protocol.ErrConnectionFailed, "ROOM_NOT_FOUND"))
		return
	}

	resp := lookResponse{
		Room: lookRoom{
			ID:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			Exits:       room.Exits,
		},
		Players: c.hub.PlayersInRoom(c.currentRoomID),
		Items:   []string{},
		NPCs:    []string{},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		c.send <- []byte(protocol.FormatErr(protocol.ErrSendFailed, "SERIALIZATION_FAILED"))
		return
	}

	c.send <- []byte(protocol.FormatOK(string(data)))
}

func (c *Client) handleMove(args []string) {
	if len(args) == 0 {
		c.send <- []byte(protocol.FormatErr(protocol.ErrNoExit, "NO_EXIT"))
		return
	}
	direction := strings.ToLower(strings.TrimSpace(args[0]))

	room, ok := c.hub.worldMap.RoomByID(c.currentRoomID)
	if !ok {
		c.send <- []byte(protocol.FormatErr(protocol.ErrConnectionFailed, "ROOM_NOT_FOUND"))
		return
	}

	nextRoomID, exists := room.Exits[direction]
	if !exists {
		c.send <- []byte(protocol.FormatErr(protocol.ErrNoExit, "NO_EXIT"))
		return
	}

	oldRoomID := c.currentRoomID
	c.hub.BroadcastRoom(
		oldRoomID,
		[]byte(protocol.FormatEvt("ROOM", "PRESENCE LEAVE", c.username)),
		c,
	)

	c.currentRoomID = nextRoomID
	c.send <- []byte(protocol.FormatOK("room=" + nextRoomID))

	c.hub.BroadcastRoom(
		nextRoomID,
		[]byte(protocol.FormatEvt("ROOM", "PRESENCE ENTER", c.username)),
		c,
	)
}

func (c *Client) handleChat(cmd protocol.Command) {
	if len(cmd.Args) < 2 {
		c.send <- []byte(protocol.FormatErr(protocol.ErrUnknownCommand, "USAGE: CHAT <scope> <message>"))
		return
	}

	scope := strings.ToUpper(cmd.Args[0])
	message := strings.TrimSpace(cmd.Raw[len(cmd.Args[0]):])

	evt := []byte(protocol.FormatEvt(scope, "CHAT", c.username+" "+message))

	switch scope {
	case "GLOBAL":
		c.hub.BroadcastGlobal(evt)

	case "ROOM":
		c.hub.BroadcastRoom(c.currentRoomID, evt, nil)

	case "GROUP":
		if c.groupName == "" {
			c.send <- []byte(protocol.FormatErr(protocol.ErrNotInGroup, "NOT_IN_GROUP"))
			return
		}
		c.hub.BroadcastGroup(c.groupName, evt, nil)

	default:
		c.send <- []byte(protocol.FormatErr(protocol.ErrUnknownCommand, "UNKNOWN_SCOPE"))
		return
	}

	c.send <- []byte(protocol.FormatOK(""))
}

func (c *Client) handleWho() {
	c.send <- []byte(protocol.FormatOK(fmt.Sprintf("players=%d", c.hub.PlayerCount())))
}

func (c *Client) handleQuit() {
	// Direct write: guarantees delivery before readPump's deferred Close().
	c.conn.Write([]byte(protocol.FormatOK("bye")))
}

func (c *Client) handleGroup(args []string) {
	if len(args) == 0 {
		c.send <- []byte(protocol.FormatErr(protocol.ErrUnknownCommand, "USAGE: GROUP <CREATE|INVITE|JOIN|LEAVE>"))
		return
	}
	scope := strings.ToUpper(args[0])

	switch scope {
	case "CREATE":
		if c.groupName != "" {
			c.send <- []byte(protocol.FormatErr(protocol.ErrAlreadyInGroup, "ALREADY_IN_GROUP"))
			return
		}
		groupID := c.hub.CreateGroup(c)
		c.send <- []byte(protocol.FormatOK("group=" + groupID))

	case "LEAVE":
		if c.groupName == "" {
			c.send <- []byte(protocol.FormatErr(protocol.ErrNotInGroup, "NOT_IN_GROUP"))
			return
		}
		c.hub.LeaveGroup(c)
		c.send <- []byte(protocol.FormatOK(""))

	case "INVITE":
		if c.groupName == "" {
			c.send <- []byte(protocol.FormatErr(protocol.ErrNotInGroup, "NOT_IN_GROUP"))
			return
		} else if len(args) < 2 {
			c.send <- []byte(protocol.FormatErr(protocol.ErrUnknownCommand, "USAGE: GROUP INVITE <username>"))
			return
		} else {
			c.hub.InviteGroup(c, args[1])
			c.send <- []byte(protocol.FormatOK(""))
		}

	case "JOIN":
		if c.groupName != "" {
			c.send <- []byte(protocol.FormatErr(protocol.ErrAlreadyInGroup, "ALREADY_IN_GROUP"))
			return
		} else if len(args) < 2 {
			c.send <- []byte(protocol.FormatErr(protocol.ErrUnknownCommand, "USAGE: GROUP JOIN <group-id>"))
			return
		} else if groupID, ok := c.hub.JoinGroup(c, args[1]); ok {
			c.send <- []byte(protocol.FormatOK("group=" + groupID))
			return
		} else {
			c.send <- []byte(protocol.FormatErr(protocol.ErrNotInGroup, "NOT_IN_GROUP"))
		}

	default:
		c.send <- []byte(protocol.FormatErr(protocol.ErrUnknownCommand, "UNKNOWN_SCOPE"))
		return
	}
}
