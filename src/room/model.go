package room

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Room struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name"`
	OwnerID   bson.ObjectID `bson:"owner_id" json:"owner_id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
