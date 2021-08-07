package main

import (
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	. "github.com/zneix/zniksbot/bot"
	db "github.com/zneix/zniksbot/mongo"
	"github.com/zneix/zniksbot/utils"
)

// TODO: Store chnnels in e.g. database instead of hardcoding them
var channels = map[string]*Channel{
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

func initChannels() {
	for ID, chn := range channels {
		// Initialize default values
		chn.Cooldowns = make(map[string]time.Time)
		chn.QueueChannel = make(chan *QueuedMessage)

		// Start message queue routine
		go SendToChannel(chn.QueueChannel, ID)

		// JOIN the channel
		Zniksbot.Client.Join(chn.Login)
		SendTwitchMessage(ID, "HONEYDETECTED ❗")
	}
}

func main() {
	log.Println("Starting zniksbot!")

	mongoClient := db.Connect()

	oauth, _ := utils.GetEnv("OAUTH", true)

	Zniksbot = &Bot{
		Client:    twitch.NewClient("zniksbot", oauth),
		Mongo:     mongoClient,
		Channels:  channels,
		Commands:  initCommands(),
		StartTime: time.Now(),
	}

	registerEventHandlers()
	initChannels()

	err := Zniksbot.Client.Connect()

	if err != nil {
		log.Fatalln(err.Error())
	}
}
