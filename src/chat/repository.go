package chat

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	Collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Collection: db.Collection("messages"),
	}
}

func (r *Repository) SaveMessage(m *Message) error {
	m.CreatedAt = time.Now()
	_, err := r.Collection.InsertOne(context.TODO(), m)
	return err
}