package main

import (
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/zneixbot-twitch/pkg/bot"
)

const prefix = "z!"

func registerEventHandlers() {
	// Authenticated with IRC
	zb.Client.OnConnect(func() {
		log.Println("connected to IRC")
	})

	// PRIVMSG
	zb.Client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if !strings.HasPrefix(message.Message, prefix) {
			return
		}

		args := strings.Fields(message.Message)
		command := args[0][len(prefix):]
		args = args[1:]

		if len(command) == 0 {
			return
		}

		handleCommands(message, command, args)
	})

	// USERSTATE
	zb.Client.OnUserStateMessage(func(message twitch.UserStateMessage) {
		// Check for user type change that signifies ratelimit in the current channel
		thisRatelimit := bot.RatelimitMsgNormal
		for key := range message.User.Badges {
			if key == "moderator" || key == "vip" || key == "broadcaster" {
				thisRatelimit = bot.RatelimitMsgElevated
				break
			}
		}

		channel := zb.Channels[zb.Logins[message.Channel]]
		if thisRatelimit != channel.Ratelimit {
			log.Printf("Changed ratelimit in #%s from %v to %v", channel.Login, channel.Ratelimit, thisRatelimit)
			channel.Ratelimit = thisRatelimit
		}
	})
}
