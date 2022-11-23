package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewStorage(cfg *models.Config) (*Storage, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.f4rrgdp.mongodb.net/?retryWrites=true&w=majority", cfg.DBUser, cfg.DBPassword)).
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &Storage{
		client: client,
	}, nil
}

func (s *Storage) SaveUser(ctx context.Context, telegramUserID int64, password, city string) (*mongo.InsertOneResult, error) {
	usr := User{
		TelegramUserID: telegramUserID,
		PasswordHash:   password,
		City:           city,
	}
	id, err := s.getUserColection().InsertOne(context.Background(), usr)
	if err != nil {
		return nil, fmt.Errorf("s.userCollection.InsertOne(ctx, usr) in SaveUser(...) falied  %v", err)
	}
	return id, nil
}

func (s *Storage) FindUserPerTGId(ctx context.Context, telegramUserID int64) (*User, error) {
	filterCursor, err := s.getUserColection().Find(context.Background(), bson.M{"telegram_user_id": telegramUserID})
	if err != nil {
		return nil, fmt.Errorf("s.userCollection.Find(ctx, bson.M{\"telegram_user_id\": telegramUserID}) in the FindUser(...) falied  %v", err)
	}
	var users []User
	if err = filterCursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	if len(users) > 1 {
		return nil, errors.New("users in database are more then 1 at the FindUser(....)")
	} else if len(users) == 0 {
		return nil, nil
	}
	s.User = &users[0]
	return s.User, nil
}

func (s *Storage) SaveEvent(ctx context.Context, startTime, name string) (*mongo.InsertOneResult, error) {
	evn := Event{
		OwnerID:   s.User.ID,
		EventTime: startTime,
		EventName: name,
	}
	id, err := s.getEventsColection().InsertOne(context.Background(), evn)
	if err != nil {
		return nil, fmt.Errorf("s.userCollection.InsertOne(ctx, usr) in SaveEvent(...) falied  %v", err)
	}
	return id, nil
}

func (s *Storage) FindEvents(ctx context.Context) ([]Event, error) {
	filterCursor, err := s.getEventsColection().Find(context.Background(), bson.M{"owner_id": s.User.ID})
	if err != nil {
		return nil, fmt.Errorf("s.eventsCollection.Find(ctx, bson.M{\"owner_id\": s.event.OwnerID})) in the FindEvent(...) falied  %v", err)
	}
	var events []Event

	if err = filterCursor.All(context.Background(), &events); err != nil {
		return nil, fmt.Errorf("filterCursor.All(ctx, &s.event) in the FindEvent(...) falied  %v", err)
	}
	s.Events = events
	return s.Events, nil
}

func (s *Storage) getUserColection() *mongo.Collection {
	return s.client.Database("forecast_users").Collection("users")
}

func (s *Storage) getEventsColection() *mongo.Collection {
	return s.client.Database("forecast_users").Collection("users")
}
