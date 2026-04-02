package chat

import (
	"WS_GIN_GOZIL/src/common"
	"WS_GIN_GOZIL/src/room"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var UPGRADER = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServerWS(c *gin.Context) {
	roomID := c.Query("room")
	userID := c.Query("user")

	conn, err := UPGRADER.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Websocket upgrade failed", err)
		return
	}

	client := &Client{
		Conn:   conn,
		UserID: userID,
		RoomID: roomID,
		Send:   make(chan []byte),
	}

	WS.Register <- client

	room.RoomMembers.Join(roomID, userID)
	joinMsg, _ := json.Marshal(map[string]any{
		"type":   "join",
		"roomID": roomID,
		"userID": userID,
		"users":  room.RoomMembers.GetUsers(roomID),
	})

	WS.Broadcast <- &MessagePayload{
		RoomID:  roomID,
		Message: joinMsg,
	}

	go client.ReadPump()
	client.WritePump()
}

func (c *Client) ReadPump() {
	chatRepo := NewRepository(common.MongoConnect())
	defer func() {
		// 1. Delete user from room present
		room.RoomMembers.Leave(c.RoomID, c.UserID)
		// 2. Broadcast news to all clients inside the room
		leaveMsg, _ := json.Marshal(map[string]any{
			"type":   "leave",
			"roomID": c.RoomID,
			"userID": c.UserID,
			"users":  room.RoomMembers.GetUsers(c.RoomID),
		})

		WS.Broadcast <- &MessagePayload{
			RoomID:  c.RoomID,
			Message: leaveMsg,
		}
		WS.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		if m.Content == "" {
			continue
		}

		m.SenderID, _ = bson.ObjectIDFromHex(c.UserID)
		m.RoomID = c.RoomID

		chatRepo.SaveMessage(&m)

		WS.Broadcast <- &MessagePayload{
			RoomID:  c.RoomID,
			Message: msg,
		}
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
