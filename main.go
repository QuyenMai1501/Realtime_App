package main

import (
	"WS_GIN_GOZIL/src/auth"
	"WS_GIN_GOZIL/src/chat"
	"WS_GIN_GOZIL/src/common"
	"WS_GIN_GOZIL/src/friend"
	"WS_GIN_GOZIL/src/notify"
	"WS_GIN_GOZIL/src/room"
	"WS_GIN_GOZIL/src/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	common.LoadEnv()
	db := common.MongoConnect()
	userRepo := user.NewRepository(db)
	friendRepo := friend.NewRepository(db)
	roomRepo := room.NewRepository(db)

	if err := room.EnsureRoomIndex(roomRepo.Rooms); err != nil {
		log.Fatalf("Không thể đánh index vì document đã tồn tại trong collection")
	}

	userController := user.NewController(userRepo)
	authController := auth.NewController(userRepo)
	friendController := friend.NewController(friendRepo)
	roomController := room.NewController(roomRepo)

	go chat.WS.Run()
	go notify.NotifyWS.Run()

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Welcome to my chat server version 0.1")
	})

	r.POST("/api/register", userController.Register)
	r.PATCH("/api/user", auth.JWTMiddleware(), userController.UpdatePatch)
	r.POST("/api/login", authController.Login)
	r.POST("/api/friend/request", auth.JWTMiddleware(), friendController.SendRequest)
	r.POST("/api/friend/accept", auth.JWTMiddleware(), friendController.AcceptRequest)
	r.POST("/api/friend/reject", auth.JWTMiddleware(), friendController.RejectRequest)
	r.GET("/api/friend/list", auth.JWTMiddleware(), friendController.ShowListFriend)
	r.POST("/api/room", auth.JWTMiddleware(), roomController.Create)
	r.GET("/api/rooms", auth.JWTMiddleware(), roomController.List)

	r.GET("/ws", chat.ServerWS)
	r.GET("/ws/notify", notify.ServerWS)

	port := common.GetEnv("PORT")
	fmt.Println("Server is running at http://localhost" + port)
	r.Run(port)
}
