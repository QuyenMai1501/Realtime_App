package friend

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	FriendRequest *mongo.Collection
	Friends       *mongo.Collection
	Users         *mongo.Collection
}

type FriendListItem struct {
	ID       bson.ObjectID `json:"id" bson:"_id"`
	Username string        `json:"username" bson:"username"`
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		FriendRequest: db.Collection("friend_requests"),
		Friends:       db.Collection("friends"),
		Users:         db.Collection("users"),
	}
}

func (r *Repository) SendRequest(fromId, toID bson.ObjectID) error {
	req := FriendRequest{
		FromUserID: fromId,
		ToUserID:   toID,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}
	_, err := r.FriendRequest.InsertOne(context.TODO(), req)
	return err
}

func (r *Repository) AcceptRequest(requestID bson.ObjectID) error {
	_, err := r.FriendRequest.UpdateByID(context.TODO(), requestID, bson.M{"$set": bson.M{
		"status": "accept",
	}})
	if err != nil {
		return err
	}

	var req FriendRequest
	err = r.FriendRequest.FindOne(context.TODO(), bson.M{"_id": requestID}).Decode(&req)
	if err != nil {
		return err
	}

	friend := Friend{
		User1:     req.FromUserID,
		User2:     req.ToUserID,
		CreatedAt: time.Now(),
	}
	_, err = r.Friends.InsertOne(context.TODO(), friend)
	return err
}

func (r *Repository) RejectRequest(requestID bson.ObjectID) error {
	_, err := r.FriendRequest.UpdateByID(context.TODO(), requestID, bson.M{"$set": bson.M{
		"status": "reject",
	}})
	return err
}

func (r *Repository) ListFriends(userID bson.ObjectID) ([]bson.ObjectID, error) {
	cursor, err := r.Friends.Find(context.TODO(), bson.M{
		"$or": []bson.M{
			{"user1": userID},
			{"user2": userID},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var ids []bson.ObjectID
	for cursor.Next(context.TODO()) {
		var f Friend
		if err := cursor.Decode(&f); err == nil {
			if f.User1 == userID {
				ids = append(ids, f.User2)
			} else {
				ids = append(ids, f.User1)
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return ids, nil

}

func (r *Repository) GetRequestByID(id bson.ObjectID) (*FriendRequest, error) {
	var req FriendRequest
	err := r.FriendRequest.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}