package callsignal

import (
	"WS_GIN_GOZIL/src/chat"
	"encoding/json"
)

type SignalPayload struct {
	ToUserID string          `json:"to_user_id"`
	FromUser string          `json:"from_user"`
	Type     string          `json:"type"` // "offer", "answer", "ice"
	Data     json.RawMessage `json:"data"` // sdp hoặc candidate
}

var signalingClients = make(map[string]*chat.Client)
