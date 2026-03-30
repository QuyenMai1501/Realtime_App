package room

import (
	"sync"
)

var RoomMembers = NewPresenceTracker()

type PresenceTracker struct {
	mu     sync.Mutex
	online map[string]map[string]bool // roomID -> map[userID]bool
}

func NewPresenceTracker() *PresenceTracker {
	return &PresenceTracker{
		online: make(map[string]map[string]bool),
	}
}

func (pt *PresenceTracker) Join(roomID, userID string)  {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if pt.online[roomID] == nil {
		pt.online[roomID] = make(map[string]bool)
	}
	pt.online[roomID][userID] = true
}

func (pt *PresenceTracker) Leave(roomID, userID string)  {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if pt.online[roomID] != nil {
		delete(pt.online[roomID], userID)
		if len(pt.online[roomID]) == 0 {
			delete(pt.online, roomID)
		}
	}
}

func (pt *PresenceTracker) GetUsers(roomID string) []string {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	var users []string
	for uid := range pt.online[roomID] {
		users = append(users, uid)
	}
	return users
}