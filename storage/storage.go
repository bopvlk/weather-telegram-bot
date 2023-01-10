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

type MongoStorage interface {
	UserRepository() *UserRepository
	EventRepository() *EventRepository
}

func NewStorage(ctx context.Context, cfg *models.Config) (MongoStorage, error) {
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

func newUserRepository(c *mongo.Client) *UserRepository {
	return &UserRepository{userCollections: c.Database(database).Collection(userCollection)}
}

func (c *client) UserRepository() *UserRepository {
	return newUserRepository(c.c)
}

func newEventRepository(c *mongo.Client) *EventRepository {
	return &EventRepository{eventCollections: c.Database(database).Collection(eventsCollection)}
}

func (c *client) EventRepository() *EventRepository {
	return newEventRepository(c.c)
}
