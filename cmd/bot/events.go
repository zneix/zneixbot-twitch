package main

import (
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
	. "github.com/zneix/zneixbot-twitch/pkg/bot"
)

const prefix = "z!"

func registerEventHandlers() {
	Zniksbot.Client.OnConnect(func() {
		log.Println("connected to IRC")
	})

	Zniksbot.Client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if !strings.HasPrefix(message.Message, prefix) {
			return
		}

		args := strings.Fields(message.Message)
		command := args[0][len(prefix):]
		args = args[1:]

		if len(command) == 0 {
			return
		}

		log.Printf("command: %s ; args: %s", command, args)

		handleCommands(message, command, args)
	})
}
