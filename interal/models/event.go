package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OwnerID   primitive.ObjectID `bson:"owner_id,omitempty"`
	EventTime string             `bson:"event_time,omitempty"`
	EventName string             `bson:"event_name,omitempty"`
}
