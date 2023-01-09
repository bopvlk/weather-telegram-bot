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

type client struct {
	c *mongo.Client
}

type Store interface {
	UserRepository() *UserRepository
	EventRepository() *EventRepository
}

func NewStorage(ctx context.Context, cfg *models.Config) (Store, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@mongodb", cfg.DBUser, cfg.DBPassword)).
		SetServerAPIOptions(serverAPIOptions)

	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err = c.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &client{c: c}, nil
}

func (c *client) UserRepository() *UserRepository {
	return &UserRepository{
		userCollections: c.c.Database(database).Collection(userCollection),
	}
}

func (c *client) EventRepository() *EventRepository {
	return &EventRepository{
		eventCollections: c.c.Database(database).Collection(eventsCollection),
	}
}
