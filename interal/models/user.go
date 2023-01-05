package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TelegramUserID int64              `bson:"telegram_user_id,omitempty"`
	PasswordHash   string             `bson:"password_hash,omitempty"`
	City           string             `bson:"city,omitempty"`
}
