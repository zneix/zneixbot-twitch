package bot

import (
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
	Name        string
	LastMsg     string
	LastMsgTime time.Time
	MsgQueue    []*QueuedMessage
	Cooldowns   map[string]time.Time
}

type Command struct {
	Name        string
	Permissions int
	Cooldown    time.Duration
	Run         func(msg twitch.PrivateMessage, args []string)
}

type QueuedMessage struct {
	Message string
	Channel string
	SendAt  string
}

var (
	Zniksbot   *Bot
	tmiTimeout = 1300 * time.Millisecond
)

func SendTwitchMessage(target string, message string) {
	if len(message) == 0 {
		return
	}

	if message[0] == '.' || message[0] == '/' {
		message = ". " + message
	}

	// limitting message length to 300
	if len(message) > 300 {
		message = message[0:297] + "..."

	}

	if Zniksbot.Channels[target].LastMsg == message {
		message += " \U000E0000"
	}

	// TODO: Finish queue system
	if time.Since(Zniksbot.Channels[target].LastMsgTime) < tmiTimeout {
		Zniksbot.Channels[target].MsgQueue = append(Zniksbot.Channels[target].MsgQueue, &QueuedMessage{
			Channel: target,
			Message: message,
		})

		time.AfterFunc(tmiTimeout-time.Now().Sub(Zniksbot.Channels[target].LastMsgTime), func() {
			//
		})
	}

	// print the message
	Zniksbot.Client.Say(target, message)

	Zniksbot.Channels[target].LastMsg = message
	Zniksbot.Channels[target].LastMsgTime = time.Now()
}
