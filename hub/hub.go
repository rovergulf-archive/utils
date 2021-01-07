package hub

import (
	"context"
	"sync"
)

type Hub struct {
	sync.RWMutex
	Clients    map[int32]map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[int32]map[*Client]bool),
	}
}

func (h *Hub) Run(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case client := <-h.Register:
			// register connection // TODO check this time
			h.Lock()
			connections := h.Clients[client.UserId]
			if connections == nil {
				connections = make(map[*Client]bool)
			}
			connections[client] = true
			h.Clients[client.UserId] = connections
			h.Unlock()
		case client := <-h.Unregister:
			// unregister conn // TODO check this time
			h.Lock()
			if _, ok := h.Clients[client.UserId]; ok {
				delete(h.Clients[client.UserId], client)
				close(client.Send)
			}
			h.Unlock()
		case message := <-h.Broadcast:
			// send message
			h.Lock()
			for _, userClients := range h.Clients {
				for client := range userClients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients[client.UserId], client)
					}
				}
			}
			h.Unlock()
		}
	}
}

func (h *Hub) GetUserConnectionsByUserId(userId int32) map[*Client]bool {
	return h.Clients[userId]
}
