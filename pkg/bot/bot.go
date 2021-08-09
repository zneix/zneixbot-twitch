package bot

import (
	"fmt"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bot struct {
	Client    *twitch.Client
	Mongo     *mongo.Client
	Channels  map[string]*Channel
	Commands  map[string]*Command
	StartTime time.Time
}

type Channel struct {
	Login        string
	LastMsg      string
	QueueChannel chan *QueuedMessage
	Cooldowns    map[string]time.Time
}

type Command struct {
	Name        string
	Permissions int
	Cooldown    time.Duration
	Run         func(msg twitch.PrivateMessage, args []string)
}

type QueuedMessage struct {
	Message   string
	ChannelID string
}

var (
	tmiTimeout = 1300 * time.Millisecond
)

func SendToChannel(zb *Bot, channelID string) {
	fmt.Printf("Starting routine for %s\n", channelID)
	var channel = zb.Channels[channelID]

	for message := range channel.QueueChannel {
		// Actually send the message to the chat
		fmt.Printf("%# v\n", message) // debug
		zb.Client.Say(channel.Login, message.Message)
		// Update last sent message
		channel.LastMsg = message.Message

		// Wait for the pleb cooldown
		time.Sleep(tmiTimeout)
		fmt.Println("Unlocked " + channelID)
	}
	fmt.Println("Done with routine for " + channelID)
}

func SendTwitchMessage(channel *Channel, message string) {
	// Don't attempt to send an empty message
	if len(message) == 0 {
		return
	}

	// Escape commands
	// TODO: Allow some commands to go through, e.g. /me
	if message[0] == '.' || message[0] == '/' {
		message = ". " + message
	}

	// limitting message length to 300
	// TODO: Investigate changing the limit based on bot's state in the channel and other settings
	if len(message) > 300 {
		message = message[0:297] + "..."

	}

	// Append magic character at the end of the message if it's a duplicate
	if channel.LastMsg == message {
		message += " \U000E0000"
	}

	// Send message object to the message processing queue
	channel.QueueChannel <- &QueuedMessage{
		Message: message,
	}
}
