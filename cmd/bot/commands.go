package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	. "github.com/zneix/zneixbot-twitch/pkg/bot"
	"github.com/zneix/zneixbot-twitch/pkg/utils"
)

const (
	ivrAPI = "https://api.ivr.fi"
)

var (
	httpClient = &http.Client{
		Timeout: 15 * time.Second,
	}
)

func initCommands() map[string]*Command {
	commands := make(map[string]*Command)

	commands["ping"] = &Command{
		Name:        "ping",
		Permissions: 0,
		Cooldown:    5000 * time.Millisecond,
		Run: func(msg twitch.PrivateMessage, args []string) {
			SendTwitchMessage(msg.RoomID, fmt.Sprintf("hi KKona 👋 I woke up %s ago", utils.TimeSince(Zniksbot.StartTime)))
		},
	}
	commands["help"] = &Command{
		Name:        "help",
		Permissions: 0,
		Cooldown:    5000 * time.Millisecond,
		Run: func(msg twitch.PrivateMessage, args []string) {
			SendTwitchMessage(msg.RoomID, fmt.Sprintf("@%s, list of commands: ping, help", msg.User.Name))
		},
	}
	commands["chatdelay"] = &Command{
		Name:        "chatdelay",
		Permissions: 0,
		Cooldown:    5000 * time.Millisecond,
		Run: func(msg twitch.PrivateMessage, args []string) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/twitch/chatdelay/%s", ivrAPI, args[0]), nil)
			if err != nil {
				//
			}

			req.Header.Add("User-Agent", USER_AGENT)

			resp, err := httpClient.Do(req)
			if err != nil {
				//
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)

			var jsonResponse ivrAPIChatDelayResponse
			if err := json.Unmarshal(body, &jsonResponse); err != nil {
				//
			}

			fmt.Println(jsonResponse)
			if jsonResponse.Status != 200 || jsonResponse.Error != "" {
				SendTwitchMessage(msg.RoomID, "Something went wrong, perhaps the channel name you've given is invalid FeelsBadMan")
				return
			}

			SendTwitchMessage(msg.RoomID, fmt.Sprintf("The delay in %s's channel is set to %d miliseconds OMGScoots", jsonResponse.Username, jsonResponse.Delay))

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
	log.Println(time.Since(Zniksbot.Channels[msg.RoomID].Cooldowns[msg.User.ID]))
	if time.Since(Zniksbot.Channels[msg.RoomID].Cooldowns[msg.User.ID]) < cmd.Cooldown {
		return
	}

	cmd.Run(msg, args)

	// apply cooldown
	if msg.User.ID != "99631238" {
		Zniksbot.Channels[msg.RoomID].Cooldowns[msg.User.ID] = time.Now()
	}
}
