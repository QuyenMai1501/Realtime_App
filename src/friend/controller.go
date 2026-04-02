package friend

import (
	"WS_GIN_GOZIL/src/auth"
	"WS_GIN_GOZIL/src/notify"
	"WS_GIN_GOZIL/src/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Controller struct {
	repo *Repository
	UserRepo *user.Repository
}

func NewController(r *Repository, ur *user.Repository) *Controller {
	return &Controller{repo: r, UserRepo: ur}
}

func (ctrl *Controller) SendRequest(c *gin.Context) {
	var input struct {
		ToUserID string `json:"to_user_id"`
	}

	c.BindJSON(&input)

	fromID := c.MustGet(auth.UserIDKey).(string)
	fromObjID, _ := bson.ObjectIDFromHex(fromID)
	toObjID, _ := bson.ObjectIDFromHex(input.ToUserID)

	if err := ctrl.repo.SendRequest(fromObjID, toObjID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can not send friend request"})
		return
	}

	// Gửi notify cho người nhận
	notify.SendToUser(toObjID.Hex(), "Bạn có một lời mời kết bạn!")

	c.JSON(http.StatusOK, gin.H{"message": "Friend Request Sent!"})
}

func (ctrl *Controller) AcceptRequest(c *gin.Context) {
	var input struct {
		RequestID string `json:"request_id"`
	}

	c.BindJSON(&input)

	requestObjID, _ := bson.ObjectIDFromHex(input.RequestID)

	if err := ctrl.repo.AcceptRequest(requestObjID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can not accept friend request"})
		return
	}
	req, err := ctrl.repo.GetRequestByID(requestObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	notify.SendToUser(req.FromUserID.Hex(), "Lời mời kết bạn của bạn đã được chấp nhận!")
	c.JSON(http.StatusOK, gin.H{"message": "Friend Request Accepted!"})
}

func (ctrl *Controller) RejectRequest(c *gin.Context) {
	var input struct {
		RequestID string `json:"request_id"`
	}

	c.BindJSON(&input)

	requestObjID, _ := bson.ObjectIDFromHex(input.RequestID)

	if err := ctrl.repo.RejectRequest(requestObjID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can not reject friend request"})
		return
	}
	req, err := ctrl.repo.GetRequestByID(requestObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	notify.SendToUser(req.FromUserID.Hex(), "Lời mời kết bạn của bạn đã bị từ chối!")

	c.JSON(http.StatusOK, gin.H{"message": "Friend Request Rejected!"})
}

func (ctrl *Controller) ListMyFriend(c *gin.Context) {
	userIDStr := c.MustGet(auth.UserIDKey).(string)
	userID, _ := bson.ObjectIDFromHex(userIDStr)
	friendIDs, _ := ctrl.repo.ListFriends(userID)
	if len(friendIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"friends": []gin.H{},
		})
		return
	}
	users := ctrl.UserRepo.FindManyByIDs(friendIDs)
	result := make([]gin.H, 0, len(users))
	for _, u := range users {
		result = append(result, gin.H{
			"id":       u.ID.Hex(),
			"username": u.Username,
			"email":    u.Email,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"friends": result,
	})
}