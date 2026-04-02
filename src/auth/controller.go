package auth

import (
	"WS_GIN_GOZIL/src/common"
	"WS_GIN_GOZIL/src/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserController struct {
	UserRepo *user.Repository
}

func NewController(repo *user.Repository) *UserController {
	return &UserController{
		UserRepo: repo,
	}
}

func (ctrl *UserController) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid input"})
		return
	}

	u, err := ctrl.UserRepo.FindByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Email"})
		return
	}

	if !common.CheckPassword(u.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	token, err := GenerateToken(u.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (ctrl UserController) MyProfile(ctx *gin.Context) {
	userIDStr := ctx.MustGet("user_id").(string)
	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid User ID"})
		return
	}

	u := ctrl.UserRepo.FindByID(userID)
	ctx.JSON(http.StatusOK, gin.H{"id": u.ID.Hex(), "username": u.Username, "email": u.Email, "created_at": u.CreatedAt})
}
