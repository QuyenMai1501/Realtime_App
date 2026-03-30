package notify

import (
	"WS_GIN_GOZIL/src/chat"
	"sync"
)

type NotifyHub struct {
	Clients    map[string]map[*chat.Client]bool
	Register   chan *chat.Client
	Unregister chan *chat.Client
	Broadcast  chan *NotifyPayload
	mu         sync.Mutex
}

type NotifyPayload struct {
	UserID  string
	Message []byte
}

var NotifyWS = NewNotifyHub()

func NewNotifyHub() *NotifyHub {
	return &NotifyHub{
		Clients:    map[string]map[*chat.Client]bool{},
		Register:   make(chan *chat.Client),
		Unregister: make(chan *chat.Client),
		Broadcast:  make(chan *NotifyPayload),
	}
}

func (h *NotifyHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.UserID] == nil {
				h.Clients[client.UserID] = make(map[*chat.Client]bool)
			}
			h.Clients[client.UserID][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients[client.UserID], client)
				close(client.Send)
			}
			h.mu.Unlock()

		case payload := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients[payload.UserID] {
				select {
				case client.Send <- payload.Message:
				default:
					close(client.Send)
					delete(h.Clients[payload.UserID], client)
				}
			}
			h.mu.Unlock()
		}
	}
}