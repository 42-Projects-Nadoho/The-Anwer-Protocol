package server

import (
	"fmt"
	"log/slog"
	"time"

	"the_answer_protocol/tap/src/internal/protocol"
	"the_answer_protocol/tap/src/internal/world"
)

type DynamicNPC struct {
	ID      string
	NPCType string
	HP      int
}

type Hub struct {
	clients   map[*Client]bool
	usernames map[string]*Client
	groups    map[string]map[*Client]bool
	nextGroup int

	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	// ops carries closures that need exclusive access to hub state; running
	// them inside Run() keeps the hub the single writer, no mutex needed.
	ops chan func()

	worldMap *world.World
	logger   *slog.Logger

	// connAttempts tracks recent connection timestamps per IP, for the
	// rapid-connections abuse signal.
	connAttempts map[string][]time.Time

	// dynamic world state
	roomItems map[string]map[string]bool
	roomNPCs  map[string]map[string]*DynamicNPC
}

func NewHub(w *world.World, logger *slog.Logger) *Hub {
	h := &Hub{
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		ops:          make(chan func()),
		clients:      make(map[*Client]bool),
		usernames:    make(map[string]*Client),
		groups:       make(map[string]map[*Client]bool),
		worldMap:     w,
		logger:       logger,
		nextGroup:    0,
		connAttempts: make(map[string][]time.Time),
		roomItems:    make(map[string]map[string]bool),
		roomNPCs:     make(map[string]map[string]*DynamicNPC),
	}

	// Initialize dynamic state from world
	for roomID, r := range w.Rooms {
		h.roomItems[roomID] = make(map[string]bool)
		for _, itemID := range r.Items {
			h.roomItems[roomID][itemID] = true
		}

		h.roomNPCs[roomID] = make(map[string]*DynamicNPC)
		for _, spawn := range r.Spawns {
			npcData, ok := w.NPCs[spawn.NPCType]
			if !ok {
				continue
			}
			for i := 0; i < spawn.Count; i++ {
				id := spawn.NPCType
				if spawn.Count > 1 {
					id = fmt.Sprintf("%s_%d", spawn.NPCType, i+1)
				}
				h.roomNPCs[roomID][id] = &DynamicNPC{
					ID:      id,
					NPCType: spawn.NPCType,
					HP:      npcData.Stats.HP,
				}
			}
		}
	}
	return h
}

// RecordConnection registers a connection attempt from ip and returns how
// many attempts from that same ip landed within the last window — a simple
// rapid-connections abuse signal for the caller to act on.
func (h *Hub) RecordConnection(ip string) int {
	const window = 10 * time.Second
	count := 0
	h.do(func() {
		now := time.Now()
		cutoff := now.Add(-window)
		var kept []time.Time
		for _, t := range h.connAttempts[ip] {
			if t.After(cutoff) {
				kept = append(kept, t)
			}
		}
		kept = append(kept, now)
		h.connAttempts[ip] = kept
		count = len(kept)
	})
	return count
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.username != "" {
					delete(h.usernames, client.username)
				}
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		case fn := <-h.ops:
			fn()
		}
	}
}

// do runs fn inside the hub goroutine, blocking the caller until it's done.
func (h *Hub) do(fn func()) {
	done := make(chan struct{})
	h.ops <- func() {
		fn()
		close(done)
	}
	<-done
}

func (h *Hub) TryConnect(c *Client, username string) bool {
	ok := false
	h.do(func() {
		if _, taken := h.usernames[username]; taken {
			return
		}
		h.usernames[username] = c
		c.username = username
		c.authenticated = true
		ok = true
	})
	return ok
}

func (h *Hub) PlayersInRoom(roomID string) []string {
	players := []string{}
	h.do(func() {
		for c := range h.clients {
			if c.authenticated && c.currentRoomID == roomID {
				players = append(players, c.username)
			}
		}
	})
	return players
}

func (h *Hub) PlayerCount() int {
	n := 0
	h.do(func() {
		for c := range h.clients {
			if c.authenticated {
				n++
			}
		}
	})
	return n
}

// exclude may be nil to include everyone in the room.
func (h *Hub) BroadcastRoom(roomID string, message []byte, exclude *Client) {
	h.do(func() {
		for c := range h.clients {
			if c == exclude || !c.authenticated || c.currentRoomID != roomID {
				continue
			}
			select {
			case c.send <- message:
			default:
			}
		}
	})
}

func (h *Hub) BroadcastGlobal(message []byte) {
	h.broadcast <- message
}

func (h *Hub) CreateGroup(c *Client) string {
	groupName := ""
	h.do(func() {
		idGroup := h.nextGroup
		groupName = fmt.Sprintf("group-%d", idGroup)
		h.groups[groupName] = map[*Client]bool{c: true}
		c.groupName = groupName
		h.nextGroup++
	})
	return groupName
}

// JoinGroup takes a group id directly (not a username) — a deliberate RFC
// deviation, documented in the README: since groups have no leader role and
// can outlive their creator, resolving "join by naming a member" breaks as
// soon as that member has left, even though the group itself may still
// exist.
func (h *Hub) JoinGroup(c *Client, groupID string) (string, bool) {
	ok := false

	h.do(func() {
		members, exists := h.groups[groupID]
		if !exists {
			return
		}

		invited := false
		for i, g := range c.isInvited {
			if g == groupID {
				invited = true
				c.isInvited = append(c.isInvited[:i], c.isInvited[i+1:]...)
				break
			}
		}
		if !invited {
			return
		}

		evt := []byte(protocol.FormatEvt("GROUP", "JOIN", c.username))
		for member := range members {
			select {
			case member.send <- evt:
			default:
			}
		}
		members[c] = true
		c.groupName = groupID
		ok = true
	})

	return groupID, ok
}

// InviteGroup's event carries the group id in addition to the inviter
// (EVT GROUP INVITE <inviter> <group-id>, not the RFC's EVT GROUP INVITE
// <leader>) — a direct consequence of JoinGroup taking a group id: the
// invitee has to learn that id from somewhere, and this event is the only
// place that can carry it.
func (h *Hub) InviteGroup(c *Client, targetUsername string) bool {
	ok := false
	h.do(func() {
		target, exists := h.usernames[targetUsername]
		if !exists {
			return
		}
		for _, g := range target.isInvited {
			if g == c.groupName {
				return
			}
		}
		evt := []byte(protocol.FormatEvt("GROUP", "INVITE", c.username+" "+c.groupName))
		target.isInvited = append(target.isInvited, c.groupName)
		select {
		case target.send <- evt:
		default:
		}
		ok = true
	})
	return ok
}

func (h *Hub) LeaveGroup(c *Client) bool {
	ok := false
	h.do(func() {
		id := c.groupName
		if id == "" {
			return
		}
		ok = true
		c.groupName = ""
		delete(h.groups[id], c)
		evt := []byte(protocol.FormatEvt("GROUP", "LEAVE", c.username))
		for member := range h.groups[id] {
			select {
			case member.send <- evt:
			default:
			}
		}
		if len(h.groups[id]) == 0 {
			delete(h.groups, id)
		}
	})
	return ok
}

func (h *Hub) BroadcastGroup(groupID string, message []byte, exclude *Client) {
	h.do(func() {
		for member := range h.groups[groupID] {
			if member == exclude {
				continue
			}
			select {
			case member.send <- message:
			default:
			}
		}
	})
}
