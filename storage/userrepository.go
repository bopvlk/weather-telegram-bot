package storage

import (
	"context"
	"fmt"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	userCollections *mongo.Collection
}

func (ur *UserRepository) SaveUser(ctx context.Context, telegramUserID int64, password, city string) (*mongo.InsertOneResult, error) {
	usr := models.User{
		TelegramUserID: telegramUserID,
		PasswordHash:   password,
		City:           city,
	}

	id, err := ur.userCollections.InsertOne(ctx, usr)
	if err != nil {
		return nil, fmt.Errorf("s.userCollection.InsertOne(ctx, usr) in SaveUser(...) falied  %v", err)
	}
	return id, nil
}

func (ur *UserRepository) FindUser(ctx context.Context, telegramUserID int64) (*models.User, error) {
	var user models.User
	filter := bson.M{"telegram_user_id": telegramUserID}
	if err := ur.userCollections.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
