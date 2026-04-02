package callsignal

import (
	"WS_GIN_GOZIL/src/chat"
	"log"

	"github.com/gin-gonic/gin"
)

func ServeSignalingWS(c *gin.Context) {
	userID := c.Query("user")

	conn, err := chat.UPGRADER.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Websocket upgrade failed.", err)
		return
	}

	client := &chat.Client{
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []byte),
	}

	// Register client to signalingClients
	signalingClients[userID] = client

	go client.WritePump()

	// Read message
	go SignalMessage(client)
}
