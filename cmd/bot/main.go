package main

import (
	"context"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/zneixbot-twitch/pkg/bot"
	db "github.com/zneix/zneixbot-twitch/pkg/mongo"
	"github.com/zneix/zneixbot-twitch/pkg/utils"
)

// Global namespace available in main package and passed where necessary
var zb *bot.Bot

func main() {
	log.Println("Starting zneixbot " + VERSION)

	ctx := context.Background()

	oauth, _ := utils.GetEnv("OAUTH", true)

	zb = &bot.Bot{
		Client:    twitch.NewClient("zneixbot", oauth),
		Mongo:     db.NewMongoConnection(),
		Logins:    make(map[string]string),
		Commands:  initCommands(),
		StartTime: time.Now(),
	}

	registerEventHandlers()
	zb.Channels = initChannels(ctx)

	err := zb.Client.Connect()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
