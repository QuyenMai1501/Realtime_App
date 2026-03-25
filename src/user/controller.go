package user

import (
	"WS_GIN_GOZIL/src/common"
	"errors"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Controller struct {
	Repo *Repository
}

func NewController(repo *Repository) *Controller {
	return &Controller{
		Repo: repo,
	}
}

func (ctrl *Controller) Register(ctx *gin.Context) {
	var input User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON body"})
		return
	}

	hashed, _ := common.HashPassword(input.Password)
	input.Password = hashed

	// Call Repo
	if err := ctrl.Repo.Create(&input); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Register success!"})
}

func (ctrl *Controller) UpdatePatch(ctx *gin.Context) {
	var input UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	rawUserID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := rawUserID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
		return
	}

	userObjID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	update := bson.M{}

	if input.Username != nil {
		username := strings.TrimSpace(*input.Username)
		if username == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "username cannot be empty"})
			return
		}
		update["username"] = username
	}

	if input.Email != nil {
		email := strings.TrimSpace(*input.Email)
		if email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email cannot be empty"})
			return
		}
		update["email"] = email
	}

	if input.Password != nil {
		password := strings.TrimSpace(*input.Password)
		if password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "password cannot be empty"})
			return
		}

		hashed, err := common.HashPassword(password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		update["password"] = hashed
	}

	if err := ctrl.Repo.UpdateByID(userObjID, update); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Update user success!"})
}
