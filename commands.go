package main

import (
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/zniksbot/utils"
)

func handleCommands(msg twitch.PrivateMessage, cmd string, args []string) {
	switch cmd {
	case "ping":
		utils.SendMessage(msg.Channel, "hi KKona 👋")
	}
}
