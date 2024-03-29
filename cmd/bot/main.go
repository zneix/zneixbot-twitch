package main

import (
	"context"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/zneixbot-twitch/internal/eventsub"
	"github.com/zneix/zneixbot-twitch/internal/helixclient"
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

	helixClient, err := helixclient.New()
	if err != nil {
		log.Fatalf("[Helix] Error while initializing client: %s\n", err)
	}

	self := &bot.Self{
		Login:   "zneixbot",
		OAuth:   oauth,
		BotType: bot.BotTypeVerified,
	}

	zb = &bot.Bot{
		TwitchIRC: twitch.NewClient(self.Login, self.OAuth),
		Helix:     helixClient,
		Mongo:     db.NewMongoConnection(),
		Logins:    make(map[string]string),
		Commands:  initCommands(),
		Self:      self,
		StartTime: time.Now(),
	}

	registerEventHandlers()
	zb.Users = initUsers(ctx)
	zb.Channels = initChannels(ctx)

	waitForEventsubWebServerClose := make(chan struct{})
	go eventsub.InitializeWebServer(zb, waitForEventsubWebServerClose)
	// TODO: Handle waitForEventsubWebServerClose being closed without blocking

	err = zb.TwitchIRC.Connect()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
