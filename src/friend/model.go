package friend

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FriendRequest struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	FromUserID bson.ObjectID `bson:"from_user_id" json:"from_user_id"`
	ToUserID   bson.ObjectID `bson:"to_user_id" json:"to_user_id"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`

	// pending, accepted, rejected
	Status string `bson:"status" json:"status"`
}

type Friend struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	User1     bson.ObjectID `bson:"user1" json:"user1"`
	User2     bson.ObjectID `bson:"user2" json:"user2"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
