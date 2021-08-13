package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/zneixbot-twitch/pkg/bot"
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

func initCommands() (commands map[string]*bot.Command) {
	commands = make(map[string]*bot.Command)

	commands["ping"] = &bot.Command{
		Name:        "ping",
		Description: "Pings the bot see if it's alive.",
		Usage:       "",
		Permission:  bot.PermissionNone,
		Cooldown:    5000 * time.Millisecond,
		Run: func(msg twitch.PrivateMessage, args []string) (err bot.CommandError) {
			channel := zb.Channels[msg.RoomID]
			channel.Send(fmt.Sprintf("hi KKona 👋 I woke up %s ago", utils.TimeSince(zb.StartTime)))
			return
		},
	}
	commands["help"] = &bot.Command{
		Name:        "help",
		Description: "Displays a list of commands.",
		Usage:       "",
		Permission:  bot.PermissionNone,
		Cooldown:    5000 * time.Millisecond,
		Run: func(msg twitch.PrivateMessage, args []string) (err bot.CommandError) {
			channel := zb.Channels[msg.RoomID]
			channel.Send(fmt.Sprintf("@%s, list of commands: ping, help", msg.User.Name))
			return
		},
	}
	commands["chatdelay"] = &bot.Command{
		Name:        "chatdelay",
		Description: "Checks the chat delay in the target channel.",
		Usage:       "<channel name>",
		Permission:  bot.PermissionNone,
		Cooldown:    5000 * time.Millisecond,
		Run: func(msg twitch.PrivateMessage, args []string) (cmdErr bot.CommandError) {
			if len(args) < 1 {
				return bot.CommandErrorInvalidArguments
			}

			channel := zb.Channels[msg.RoomID]

			req, err := http.NewRequest("GET", fmt.Sprintf("%s/twitch/chatdelay/%s", ivrAPI, args[0]), nil)
			if err != nil {
				log.Println("Error in chatdelay command in http.NewRequest: " + err.Error())
				return bot.CommandErrorInternal
			}

			req.Header.Add("User-Agent", USER_AGENT)

			resp, err := httpClient.Do(req)
			if err != nil {
				log.Println("Error in chatdelay command in httpClient.Do: " + err.Error())
				return bot.CommandErrorHTTPRequest
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)

			var jsonResponse ivrAPIChatDelayResponse
			if err := json.Unmarshal(body, &jsonResponse); err != nil {
				log.Println("Error in chatdelay command in json.Unmarshal: " + err.Error())
				return bot.CommandErrorHTTPRequest
			}

			fmt.Println(jsonResponse)
			if jsonResponse.Status != 200 || jsonResponse.Error != "" {
				channel.Send("Something went wrong, perhaps the channel name you've given is invalid FeelsBadMan")
				return bot.CommandErrorReply
			}

			channel.Send(fmt.Sprintf("The delay in %s's channel is set to %d miliseconds OMGScoots", jsonResponse.Username, jsonResponse.Delay))

			return
		},
	}
	commands["set"] = &bot.Command{
		Name:        "set",
		Description: "Sets various bot-related flags in the database",
		Usage:       "<flag> <target> <value>",
		Permission:  bot.PermissionAdmin,
		Cooldown:    0 * time.Millisecond,
		Run: func(msg twitch.PrivateMessage, args []string) (cmdErr bot.CommandError) {
			if len(args) < 3 {
				return bot.CommandErrorInvalidArguments
			}

			channel := zb.Channels[msg.RoomID]
			flag, found := bot.FlagMap[args[0]]

			// If an invalid flag was provided, early out and provide a list of available flags for you to use
			if !found {
				keys := make([]string, 0, len(bot.FlagMap))

				for k := range bot.FlagMap {
					keys = append(keys, k)
				}

				channel.Send(fmt.Sprintf("Invalid flag provided, available flags: %s", strings.Join(keys, ", ")))
				return bot.CommandErrorReply
			}

			// Since we already early-out right above, it's safe to assume that
			// default: will never execute, provided we handle all the bot.Flag's
			switch flag {
			case bot.FlagChannelMode:
				targetChannel := zb.Channels[args[1]]
				if targetChannel == nil {
					channel.Send("Channel with provided ID doesn't exist in the database FeelsDankMan")
					return
				}

				modeNum, err := strconv.Atoi(args[2])
				if err != nil {
					channel.Send("Failed to parse channel mode to a number FeelsDankMan")
					return
				}

				err, mode := bot.ParseChannelMode(modeNum)
				if err != nil {
					channel.Send("Provided channel mode is out of bounds FeelsDankMan")
					return
				}

				// Update the mode in the database
				err = targetChannel.ChangeMode(zb.Mongo, mode)
				if err != nil {
					//channel.Send("Failed to update the channel mode in database monkaS")
					log.Println("Error in set command in Channel.ChangeMode: " + err.Error())
					return bot.CommandErrorMongo
				} else {
					channel.Send(fmt.Sprintf("Updated channel mode for target channel to %s", channel.Mode.String()))
					return
				}
			default:
				return
			}
		},
	}

	return
}

func handleCommands(msg twitch.PrivateMessage, command string, args []string) {
	// find the command
	cmd := zb.Commands[command]

	if cmd == nil {
		return
	}

	// check for permission
	if !utils.HasBits(uint64(zb.Users[msg.User.ID].Permission), uint64(cmd.Permission)) {
		return
	}

	// handle cooldown
	if time.Since(zb.Channels[msg.RoomID].Cooldowns[msg.User.ID]) < cmd.Cooldown {
		return
	}

	cmdErr := cmd.Run(msg, args)
	channel := zb.Channels[msg.RoomID]

	// Handle error returned in the command execution
	switch cmdErr {
	case bot.CommandErrorNoError:
		// No error was returned, command was executed successfully
		// Do nothing

	case bot.CommandErrorInvalidArguments:
		// Show default usage
		// TODO: Explain better what happened
		// TODO: Consider dynamic cooldown
		channel.Send(fmt.Sprintf("Invalid arguments provided! Usage: %s%s %s", prefix, cmd.Name, cmd.Usage))

	case bot.CommandErrorInternal:
		// Return a generic message in case this ever happens
		channel.Send("Unexpected error occured monkaS @zneix")

	case bot.CommandErrorHTTPRequest:
		// Inform the user about the http error (possibly return status code?)
		channel.Send("Request to a third-party API failed FutureMan")

	case bot.CommandErrorReply:
		// An error occured, but the user was already informed about it
		// Do nothing

	case bot.CommandErrorMongo:
		// A MongoDB error occured
		// TODO: Add some better description for what exactly happened
		channel.Send("Database operation failed monkaS @zneix")

	default:
		log.Printf("Unhandled command error %v in command %s", cmdErr, cmd.Name)
	}

	// apply cooldown
	if !utils.HasBits(uint64(zb.Users[msg.User.ID].Permission), uint64(bot.PermissionSkipCooldown)) {
		log.Printf("Applying cooldown to %s", msg.User.ID)
		channel.Cooldowns[msg.User.ID] = time.Now()
	}
}
