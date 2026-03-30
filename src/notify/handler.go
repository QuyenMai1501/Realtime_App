package notify

import (
	"WS_GIN_GOZIL/src/chat"
	"WS_GIN_GOZIL/src/room"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	client := &chat.Client{
		Conn:   conn,
		UserID: userID,
		RoomID: roomID,
		Send:   make(chan []byte),
	}

	// add user to room presence
	room.RoomMembers.Join(roomID, userID)
	defer room.RoomMembers.Leave(roomID, userID)

	NotifyWS.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
