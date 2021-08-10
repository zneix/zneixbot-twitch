package main

import (
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/zneixbot-twitch/pkg/bot"
	db "github.com/zneix/zneixbot-twitch/pkg/mongo"
	"github.com/zneix/zneixbot-twitch/pkg/utils"
)

// TODO: Store chnnels in e.g. database instead of hardcoding them
var (
	zb *bot.Bot

	channels = map[string]*bot.Channel{
		"463521670": {
			Login: "zneixbot",
		},
		"99631238": {
			Login: "zneix",
		},
		"31400525": {
			Login: "supinic",
		},
	}
)

func initChannels() {
	for ID, chn := range channels {
		// Set the ID in map translating login names back to IDs
		zb.Logins[chn.Login] = ID

		// Initialize default values
		chn.Cooldowns = make(map[string]time.Time)
		chn.QueueChannel = make(chan *bot.QueuedMessage)
		chn.Ratelimit = bot.RatelimitMsgNormal

		// Start message queue routine
		go bot.SendToChannel(zb, ID)

		// JOIN the channel
		zb.Client.Join(chn.Login)
		bot.SendTwitchMessage(chn, "HONEYDETECTED ❗")
	}
}

func main() {
	log.Println("Starting zneixbot!")

	mongoClient := db.Connect()

	oauth, _ := utils.GetEnv("OAUTH", true)

	zb = &bot.Bot{
		Client:    twitch.NewClient("zneixbot", oauth),
		Mongo:     mongoClient,
		Logins:    make(map[string]string),
		Channels:  channels,
		Commands:  initCommands(),
		StartTime: time.Now(),
	}

	registerEventHandlers()
	initChannels()

	err := zb.Client.Connect()

	if err != nil {
		log.Fatalln(err.Error())
	}
}
