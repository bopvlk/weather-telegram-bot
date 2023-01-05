package storage

import (
	"context"
	"errors"
	"fmt"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	bd *DB
}

func (db *DB) getUserColection() *mongo.Collection {
	return db.client.Database("forecast_users").Collection("users")
}

func (us *UserRepository) SaveUser(ctx context.Context, telegramUserID int64, password, city string) (*mongo.InsertOneResult, error) {
	usr := models.User{
		TelegramUserID: telegramUserID,
		PasswordHash:   password,
		City:           city,
	}

	id, err := us.bd.getUserColection().InsertOne(ctx, usr)
	if err != nil {
		return nil, fmt.Errorf("s.userCollection.InsertOne(ctx, usr) in SaveUser(...) falied  %v", err)
	}
	return id, nil
}

func (ur *UserRepository) FindOneUser(ctx context.Context, telegramUserID int64) (*models.User, error) {
	filterCursor, err := ur.bd.getUserColection().Find(ctx, bson.M{"telegram_user_id": telegramUserID})
	if err != nil {
		return nil, fmt.Errorf("s.userCollection.Find(ctx, bson.M{\"telegram_user_id\": telegramUserID}) in the FindUser(...) falied  %v", err)
	}
	var users []models.User
	if err = filterCursor.All(ctx, &users); err != nil {
		return nil, err
	}
	if len(users) > 1 {
		return nil, errors.New("users in database are more then 1 at the FindUser(....)")
	} else if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}
