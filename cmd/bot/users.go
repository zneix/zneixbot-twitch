package main

import (
	"context"
	"log"
	"time"

	"github.com/zneix/zneixbot-twitch/pkg/bot"
	db "github.com/zneix/zneixbot-twitch/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func initUsers(parentCtx context.Context) (users map[string]*bot.User) {
	users = make(map[string]*bot.User)

	// Query channels from the database
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	cur, err := zb.Mongo.Collection(db.CollectionNameUsers).Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		// Deserialize user data
		var channel bot.User
		err := cur.Decode(&channel)
		if err != nil {
			log.Fatal(err)
		}

		users[(&channel).ID] = &channel
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return
}
