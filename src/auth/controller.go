package auth

import (
	"WS_GIN_GOZIL/src/common"
	"WS_GIN_GOZIL/src/user"
	"net/http"

	"github.com/gin-gonic/gin"
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
