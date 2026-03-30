package room

import (
	"WS_GIN_GOZIL/src/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Controller struct {
	Repo *Repository
}

func NewController(r *Repository) *Controller {
	return &Controller{Repo: r}
}

func (ctrl *Controller) Create(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}

	c.BindJSON(&input)

	userIDHex := c.MustGet(auth.UserIDKey).(string)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	room, err := ctrl.Repo.createRoom(input.Name, userID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (ctrl *Controller) List(c *gin.Context) {
	rooms, err := ctrl.Repo.getRoom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rooms)
}