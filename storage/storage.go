package storage

import (
	"context"
	"fmt"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	database         = "forecast_users"
	userCollection   = "users"
	eventsCollection = "events"
)

type Storage interface {
	UserRepository() *UserRepository
	EventRepository() *EventRepository
}

func NewStorage(ctx context.Context, cfg *models.Config) (Storage, error) {
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
	return storage, nil
}
