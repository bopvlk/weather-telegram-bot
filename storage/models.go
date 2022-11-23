package storage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	client *mongo.Client
	User   *User
	Events []Event
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TelegramUserID int64              `bson:"telegram_user_id,omitempty"`
	PasswordHash   string             `bson:"password_hash,omitempty"`
	City           string             `bson:"city,omitempty"`
}

type Event struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OwnerID   primitive.ObjectID `bson:"owner_id,omitempty"`
	EventTime string             `bson:"event_time,omitempty"`
	EventName string             `bson:"event_name,omitempty"`
}
