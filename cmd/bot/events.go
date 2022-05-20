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
	zb.TwitchIRC.OnConnect(func() {
		log.Println("connected to IRC")
		joinChannels()
	})

	// PRIVMSG
	zb.TwitchIRC.OnPrivateMessage(func(message twitch.PrivateMessage) {
		channel := zb.Channels[zb.Logins[message.Channel]]

		// Ignore inactive channels
		if channel.Mode == bot.ChannelModeInactive {
			return
		}

		// Handle non-commands
		if !strings.HasPrefix(message.Message, prefix) {
			handlePajbotAnnounceChain(&message)
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
	zb.TwitchIRC.OnUserStateMessage(func(message twitch.UserStateMessage) {
		channel := zb.Channels[zb.Logins[message.Channel]]

		// Ignore inactive channels
		if channel.Mode == bot.ChannelModeInactive {
			return
		}

		// Check if user type changed and update ChannelMode in the current channel accordingly
		newMode := bot.ChannelModeNormal
		for key := range message.User.Badges {
			if key == "vip" {
				newMode = bot.ChannelModeVIP
				break
			}
			if key == "moderator" || key == "broadcaster" {
				newMode = bot.ChannelModeModerator
				break
			}
		}

		if newMode != channel.Mode {
			channel.ChangeMode(zb.Mongo, newMode)
			//log.Printf("Changing mode in %s from %v to %v", channel.VerboseName(), channel.Mode.String(), newMode.String())
			//channel.Mode = newMode
		}
	})
}

const (
	pajladaUserID             = "11148817" // channel where custom events are taking place
	pajbotUserID              = "82008718"
	pajbotAnnounceChainUserID = "147677870" // icecreamdatabase replies with letter "q"
	zneixUserID               = "99631238"  // bot creator's ID for testing purposes
)

func handlePajbotAnnounceChain(message *twitch.PrivateMessage) {
	if message.RoomID != pajladaUserID || (message.User.ID != pajbotAnnounceChainUserID && message.User.ID != zneixUserID) {
		return
	}

	if message.Message != "/announce q" {
		return
	}

	zb.Channels[message.RoomID].Send(".me /announce rrrree pajaR 💢")
}
