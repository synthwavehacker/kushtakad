package state

import (
	"github.com/asdine/storm"
)

// ServerHub maintains the set of active clients and broadcasts messages to the
// clients.
type ServerHub struct {
	clients    map[*Client]bool
	broadcast  chan *ClientMsg
	register   chan *Client
	unregister chan *Client
	db         *storm.DB
}

func NewServerHub(db *storm.DB) *ServerHub {
	return &ServerHub{
		broadcast:  make(chan *ClientMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *ServerHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.sender)
				close(client.receiver)
			}
		case clientMsg := <-h.broadcast:
			for client := range h.clients {
				if client.SensorID == clientMsg.SensorID {
					select {
					case client.sender <- clientMsg:
					default:
						close(client.sender)
						close(client.receiver)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
