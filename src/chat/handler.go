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

	// add user to room presence
	room.RoomMembers.Join(roomID, userID)
	defer room.RoomMembers.Leave(roomID, userID)

	WS.Register <- client

	go client.ReadPump()
	go client.WritePump()
}

func (c *Client) ReadPump() {
	chatRepo := NewRepository(common.MongoConnect())
	defer func() {
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
