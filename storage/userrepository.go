package storage

import (
	"context"
	"fmt"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (ur *UserRepository) getColection() *mongo.Collection {
	return ur.client.Database(database).Collection(userCollection)
}

func (ur *UserRepository) SaveUser(ctx context.Context, telegramUserID int64, password, city string) (*mongo.InsertOneResult, error) {
	usr := models.User{
		TelegramUserID: telegramUserID,
		PasswordHash:   password,
		City:           city,
	}

	id, err := ur.getColection().InsertOne(ctx, usr)
	if err != nil {
		return nil, fmt.Errorf("s.userCollection.InsertOne(ctx, usr) in SaveUser(...) falied  %v", err)
	}
	return id, nil
}

func (ur *UserRepository) FindUser(ctx context.Context, telegramUserID int64) (*models.User, error) {
	var user models.User
	filter := bson.D{{Key: "telegram_user_id", Value: telegramUserID}}
	if err := ur.getColection().FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
