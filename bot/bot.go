package bot

import (
	"github.com/gempir/go-twitch-irc/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bot struct {
	Client   *twitch.Client
	Mongo    *mongo.Client
	LastMsgs map[string]string
}

var Zniksbot *Bot
