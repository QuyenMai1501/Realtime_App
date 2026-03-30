package room

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	Rooms *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Rooms: db.Collection("rooms"),
	}
}

func (r *Repository) createRoom(name string, ownerID bson.ObjectID) (*Room, error) {
	room := &Room{
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
	}

	res, err := r.Rooms.InsertOne(context.TODO(), room)
	if err != nil {
		return nil, err
	}

	room.ID = res.InsertedID.(bson.ObjectID)
	return room, nil
}

func getRoom() {
	//NOTE: Homework
}
