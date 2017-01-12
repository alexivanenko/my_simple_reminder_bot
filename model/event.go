package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	ObjectID bson.ObjectId `json:"id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Date     time.Time     `json:"date" bson:"date"`
	TimeZone string        `json:"timezone" bson:"timezone"`
	ChatID   int64         `json:"chat_id" bson:"chat_id"`
}

func (event *Event) Save() error {
	db := GetDB()
	b := db.C("events")

	event.ObjectID = bson.NewObjectId()
	_, err := b.UpsertId(event.ObjectID, event)

	return err
}
