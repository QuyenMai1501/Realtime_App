package main

import (
	"WS_GIN_GOZIL/src/auth"
	"WS_GIN_GOZIL/src/callsignal"
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
		log.Fatalf("Không thể đánh index vì document đã tồn tại trong collection.")
	}

	userController := user.NewController(userRepo)
	authController := auth.NewController(userRepo)
	friendController := friend.NewController(friendRepo, userRepo)
	roomController := room.NewController(roomRepo)

	go chat.WS.Run()
	go notify.NotifyWS.Run()

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Welcome to my chat server version 0.0.0.1")
	})
	r.POST("/api/register", userController.Register)

	r.POST("/api/login", authController.Login)
	r.GET("/api/myprofile", auth.JWTMiddleware(), authController.MyProfile)

	r.POST("/api/friend/request", auth.JWTMiddleware(), friendController.SendRequest)
	r.POST("/api/friend/accept", auth.JWTMiddleware(), friendController.AcceptRequest)
	r.GET("/api/friends", auth.JWTMiddleware(), friendController.ListMyFriend)

	r.POST("/api/room", auth.JWTMiddleware(), roomController.Create)

	r.Static("/ui", "./ui")

	r.GET("/ws", chat.ServerWS)
	r.GET("/ws/notify", notify.ServerWS)
	r.GET("/ws/signaling", callsignal.ServeSignalingWS)

	port := common.GetEnv("PORT")
	fmt.Println("Server is running at http://localhost" + port)
	// r.Run(port)

	err := r.RunTLS(port, "cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
}
