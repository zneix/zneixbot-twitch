package main

import (
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	. "github.com/zneix/zniksbot/bot"
	db "github.com/zneix/zniksbot/mongo"
	"github.com/zneix/zniksbot/utils"
)

var channels = map[string]*Channel{
	"supinic":  {Name: "supinic", Cooldowns: make(map[string]time.Time)},
	"zniksbot": {Name: "zniksbot", Cooldowns: make(map[string]time.Time)},
}

func connectToChannels() {
	for i := range channels {
		Zniksbot.Client.Join(i)
		Zniksbot.Client.Say(i, "HONEYDETECTED ❗")
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
	connectToChannels()

	err := Zniksbot.Client.Connect()

	if err != nil {
		log.Fatalln(err.Error())
	}
}
