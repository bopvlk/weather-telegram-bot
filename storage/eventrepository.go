package storage

import (
	"context"
	"fmt"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepository struct {
	eventCollections *mongo.Collection
}

func (er *EventRepository) SaveEvent(ctx context.Context, userID primitive.ObjectID, startTime, name string) (*mongo.InsertOneResult, error) {
	event := models.Event{
		OwnerID:   userID,
		EventTime: startTime,
		EventName: name,
	}
	id, err := er.eventCollections.InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("s.userCollection.InsertOne(ctx, usr) in SaveEvent(...) falied  %v", err)
	}
	return id, nil
}

func (er *EventRepository) FindEvents(ctx context.Context, userID primitive.ObjectID) ([]models.Event, error) {
	filterCursor, err := er.eventCollections.Find(ctx, bson.M{"owner_id": userID})
	if err != nil {
		return nil, fmt.Errorf("s.eventsCollection.Find(ctx, bson.M{\"owner_id\": s.event.OwnerID})) in the FindEvent(...) falied  %v", err)
	}
	var events []models.Event

	if err = filterCursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("filterCursor.All(ctx, &s.event) in the FindEvent(...) falied  %v", err)
	}
	return events, nil
}
