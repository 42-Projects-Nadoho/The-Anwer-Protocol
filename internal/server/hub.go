package server

import (
	"fmt"
	"the_answer_protocol/internal/protocol"
	"the_answer_protocol/internal/world"
)

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
}

func NewHub(w *world.World) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		ops:        make(chan func()),
		clients:    make(map[*Client]bool),
		usernames:  make(map[string]*Client),
		groups:     make(map[string]map[*Client]bool),
		worldMap:   w,
		nextGroup:  0,
	}
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

func (h *Hub) JoinGroup(c *Client, targetUsername string) (string, bool) {
	var (
		groupName string
		ok        bool
	)

	h.do(func() {
		target, exists := h.usernames[targetUsername]
		if !exists {
			return
		}

		groupName = target.groupName
		if groupName == "" {
			return
		}

		members, exists := h.groups[groupName]
		if !exists {
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
		c.groupName = groupName
		ok = true
	})

	return groupName, ok
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
