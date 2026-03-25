package main

import (
	"WS_GIN_GOZIL/src/auth"
	"WS_GIN_GOZIL/src/common"
	"WS_GIN_GOZIL/src/friend"
	"WS_GIN_GOZIL/src/user"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	common.LoadEnv()
	db := common.MongoConnect()
	userRepo := user.NewRepository(db)
	friendRepo := friend.NewRepository(db)

	userController := user.NewController(userRepo)
	authController := auth.NewController(userRepo)
	friendController := friend.NewController(friendRepo)

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

	port := common.GetEnv("PORT")
	fmt.Println("Server is running at http://localhost" + port)
	r.Run(port)
}
