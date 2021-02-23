package main

import (
	"log"

	"github.com/gempir/go-twitch-irc/v2"
	. "github.com/zneix/zniksbot/bot"
	db "github.com/zneix/zniksbot/mongo"
	"github.com/zneix/zniksbot/utils"
)

func connectToChannels(client *twitch.Client, channels []string) {
	for i := 0; i < len(channels); i++ {
		client.Join(channels[i])
		client.Say(channels[i], "HONEYDETECTED ❗")
	}
}

func main() {
	log.Println("Starting zniksbot!")

	mongoClient := db.Connect()

	oauth, _ := utils.GetEnv("OAUTH", true)
	twitchClient := twitch.NewClient("zniksbot", oauth)

	Zniksbot = &Bot{
		Client:   twitchClient,
		Mongo:    mongoClient,
		LastMsgs: make(map[string]string),
	}

	registerEventHandlers()

	connectToChannels(twitchClient, []string{"supinic", "zniksbot"})

	err := twitchClient.Connect()

	if err != nil {
		log.Fatalln(err.Error())
	}
}
