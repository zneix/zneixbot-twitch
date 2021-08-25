package main

import (
	"context"
	"log"
	"time"

	"github.com/zneix/zneixbot-twitch/internal/eventsub"
	"github.com/zneix/zneixbot-twitch/pkg/bot"
	db "github.com/zneix/zneixbot-twitch/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// initChannels queries the database instance and creates new map of all the channels inside it
func initChannels(parentCtx context.Context) (channels map[string]*bot.Channel) {
	channels = make(map[string]*bot.Channel)

	// Query channels from the database
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	cur, err := zb.Mongo.Collection(db.CollectionNameChannels).Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		// Deserialize channel data
		var channel bot.Channel
		err := cur.Decode(&channel)
		if err != nil {
			log.Fatal(err)
		}

		// In case "eventsub" wasn't found, init empty object so that we won't panic
		if channel.EventSub == nil {
			channel.EventSub = &bot.ChannelEventSub{
				Enabled:       false,
				Subscriptions: make([]*bot.ChannelEventSubSubscription, 0),
			}
		}

		// Initialize default values
		channel.Cooldowns = make(map[string]time.Time)
		channel.QueueChannel = make(chan *bot.QueuedMessage)

		channels[(&channel).ID] = &channel
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

// joinChannels performs startup actions for all the channels that are already loaded
func joinChannels() {
	for ID, channel := range zb.Channels {
		// Set the ID in map translating login names back to IDs
		zb.Logins[channel.Login] = ID

		// Start message queue routine
		go channel.Write(zb)

		// JOIN the channel
		zb.TwitchIRC.Join(channel.Login)
		//channel.Send("HONEYDETECTED ❗")

		// Subscribe to EventSub subscriptions
		if channel.EventSub.Enabled {
			for _, sub := range channel.EventSub.Subscriptions {
				eventsub.CreateChannelSubscription(zb, sub, channel.ID)
			}
		}
		// end of the loop
	}
}
