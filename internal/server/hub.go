package server

import (
	"the_answer_protocol/internal/world"
)

type Hub struct {
	clients map[*Client]bool
	broadcast	chan []byte
	register	chan *Client
	unregister	chan *Client
	worldMap	*world.World
}


func NewHub(w *world.World) *Hub {
	return &Hub{
		broadcast:	make(chan []byte),
		register:	make(chan *Client),
		unregister:	make(chan *Client),
		clients:	make(map[*Client]bool),
		worldMap:	w,
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
				close(client.send)
			}

		case message := <- h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
