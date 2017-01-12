package model

import (
	"time"

	"github.com/alexivanenko/my_simple_reminder_bot/config"
	"gopkg.in/mgo.v2"
)

var db *mgo.Database
var session *mgo.Session

func init() {
	info := &mgo.DialInfo{
		Addrs:    []string{config.String("db", "server")},
		Timeout:  60 * time.Second,
		Database: config.String("db", "name"),
		Username: config.String("db", "user"),
		Password: config.String("db", "pass"),
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db = session.DB(config.String("db", "name"))
}

func GetDB() *mgo.Database {
	return db
}

func GetSession() *mgo.Session {
	return session
}
