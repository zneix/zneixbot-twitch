package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	. "github.com/zneix/zniksbot/bot"
	"github.com/zneix/zniksbot/utils"
)

func initCommands() map[string]*Command {
	commands := make(map[string]*Command)

	commands["ping"] = &Command{
		Name:        "ping",
		Permissions: 0,
		Cooldown:    5000 * time.Millisecond,
		Run: func(msg twitch.PrivateMessage) {
			SendTwitchMessage(msg.Channel, fmt.Sprintf("hi KKona 👋 I woke up %s ago", utils.TimeSince(Zniksbot.StartTime)))

		},
	}
	commands["help"] = &Command{
		Name:        "help",
		Permissions: 0,
		Cooldown:    5000 * time.Millisecond,
		Run: func(msg twitch.PrivateMessage) {
			SendTwitchMessage(msg.Channel, fmt.Sprintf("@%s, list of commands: ping, help", msg.User.Name))
		},
	}

	return commands
}

func handleCommands(msg twitch.PrivateMessage, command string, args []string) {

	// finding the command
	cmd := Zniksbot.Commands[command]
	log.Println(cmd)

	if cmd == nil {
		return
	}

	// handling cooldowns
	log.Println(time.Since(Zniksbot.Channels[msg.Channel].Cooldowns[msg.User.ID]))
	if time.Since(Zniksbot.Channels[msg.Channel].Cooldowns[msg.User.ID]) < cmd.Cooldown {
		return
	}

	cmd.Run(msg)

	// apply cooldown
	Zniksbot.Channels[msg.Channel].Cooldowns[msg.User.ID] = time.Now()
}
