package bot

import (
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bot struct {
	Client    *twitch.Client
	Mongo     *mongo.Client
	Logins    map[string]string
	Channels  map[string]*Channel
	Commands  map[string]*Command
	BotType   BotTypeEnum
	StartTime time.Time
}

type Channel struct {
	Login        string
	LastMsg      string
	Ratelimit    time.Duration
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
	Message string
}

// BotTypeEnum represents which kind of global rate-limit the bot has
type BotTypeEnum int32

const (
	BotTypeNormal BotTypeEnum = iota
	BotTypeKnown
	BotTypeVerified
)

// ratelimit values for sending messages.
// RatelimitMsgElevated can be used when the bot is a Moderator, VIP or Broadcaster.
// RatelimitMsgNormal should be used in all other cases
const (
	RatelimitMsgNormal   = 1250 * time.Millisecond
	RatelimitMsgElevated = 100 * time.Millisecond
)
