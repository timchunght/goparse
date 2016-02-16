package models

import (
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"time"
)

type Trivia struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	EventId     string        `bson:"event_id" json:"event_id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
}
