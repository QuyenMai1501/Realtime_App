package room

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *Repository) getRoom() ([]Room, error) {
	//NOTE: Homework
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.Rooms.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	rooms := make([]Room, 0)
	if err = cursor.All(context.TODO(), &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
