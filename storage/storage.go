package storage

import (
	"context"
	"fmt"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client  *mongo.Client
	Storage Storage
}

type Storage interface {
	SaveEvent(ctx context.Context, startTime, name string) (*mongo.InsertOneResult, error)
	FindEvents(ctx context.Context) ([]models.Event, error)
	SaveUser(ctx context.Context, telegramUserID int64, password, city string) (*mongo.InsertOneResult, error)
	FindOneUser(ctx context.Context, telegramUserID int64) (*models.User, error)
}

func NewStorage(ctx context.Context, cfg *models.Config) (*DB, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@mongodb", cfg.DBUser, cfg.DBPassword)).
		SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	var storage Storage
	return &DB{
		client:  client,
		Storage: storage,
	}, nil
}
