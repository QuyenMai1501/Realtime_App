package notify

import "encoding/json"

func SendToUser(userID, content string) {
	msg := map[string]string{
		"type":    "notify",
		"content": content,
	}

	data, _ := json.Marshal(msg)
	NotifyWS.Broadcast <- &NotifyPayload{
		UserID:  userID,
		Message: data,
	}
}