package main

import (
	"fmt"
	"log"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/zniksbot/mongo"
	"github.com/zneix/zniksbot/utils"
)

func connectToChannels(client *twitch.Client, channels []string) {
	for i := 0; i < len(channels); i++ {
		client.Join(channels[i])
		client.Say(channels[i], "HONEYDETECTED ❗")
	}
}

func registerEventHandlers(client *twitch.Client) {
	client.OnConnect(func() {
		log.Println("connected to IRC")
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Println(fmt.Sprintf("[#%s] %s: %s", message.Channel, message.User.DisplayName, message.Message))
	})
}

func main() {
	log.Println("Starting zniksbot!")

	mongo.Connect()

	oauth, _ := utils.GetEnv("OAUTH", true)
	twitchClient := twitch.NewClient("zniksbot", oauth)

	registerEventHandlers(twitchClient)

	connectToChannels(twitchClient, []string{"supinic", "zniksbot"})

	err := twitchClient.Connect()

	if err != nil {
		log.Fatalln(err.Error())
	}
}
