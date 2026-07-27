package notification

import (
	"sync"

	"github.com/google/uuid"
)

type Notification struct {
	ID        string      `json:"id"`
	Text      string      `json:"text"`
	Time      string      `json:"time"`
	Read      bool        `json:"read"`
	MeetingID string      `json:"meeting_id,omitempty"`
	Type      string      `json:"type,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
}

type Hub struct {
	mu      sync.Mutex
	clients map[uuid.UUID][]chan Notification
}

var GlobalHub = &Hub{
	clients: make(map[uuid.UUID][]chan Notification),
}

func (h *Hub) Register(userID uuid.UUID) chan Notification {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan Notification, 10)
	h.clients[userID] = append(h.clients[userID], ch)
	return ch
}

func (h *Hub) Unregister(userID uuid.UUID, ch chan Notification) {
	h.mu.Lock()
	defer h.mu.Unlock()

	channels := h.clients[userID]
	for i, c := range channels {
		if c == ch {
			h.clients[userID] = append(channels[:i], channels[i+1:]...)
			close(c)
			break
		}
	}
	if len(h.clients[userID]) == 0 {
		delete(h.clients, userID)
	}
}

func (h *Hub) Publish(userID uuid.UUID, notif Notification) {
	h.mu.Lock()
	defer h.mu.Unlock()

	channels := h.clients[userID]
	for _, c := range channels {
		select {
		case c <- notif:
		default:
			// channel full or closed, skip
		}
	}
}
