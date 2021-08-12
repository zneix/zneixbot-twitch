package bot

import (
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	db "github.com/zneix/zneixbot-twitch/pkg/mongo"
)

// Self contains properties related to bot's user account
type Self struct {
	Login   string
	OAuth   string
	BotType BotType
}

type Bot struct {
	Client    *twitch.Client
	Mongo     *db.Connection
	Logins    map[string]string
	Channels  map[string]*Channel
	Commands  map[string]*Command
	Self      *Self
	StartTime time.Time
}

type Channel struct {
	Login string      `bson:"login"`
	ID    string      `bson:"id"`
	Mode  ChannelMode `bson:"mode"`

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
	Message string
}

// BotType represents which kind of global rate-limit the bot has
type BotType int32

const (
	BotTypeNormal BotType = iota
	BotTypeKnown
	BotTypeVerified
)

// ChannelMode ...
type ChannelMode int32

const (
	ChannelModeNormal ChannelMode = iota
	ChannelModeInactive
	ChannelModeVIP
	ChannelModeModerator
)

// String human-readable name of the mode. Meant to be used only in logs and output visible to end-user
func (mode ChannelMode) String() string {
	return []string{
		"Normal",
		"Inactive",
		"VIP",
		"Moderator",
	}[int(mode)]
}

// MessageRatelimit returns underlying value of messageRatelimit that corresponds to mode
func (mode ChannelMode) MessageRatelimit() time.Duration {
	return []time.Duration{
		time.Duration(messageRatelimitNormal),
		time.Duration(messageRatelimitNormal),
		time.Duration(messageRatelimitElevated),
		time.Duration(messageRatelimitElevated),
	}[int(mode)]
}

// messageRatelimit minimum value between which the bot user can send messages in target Channel
type messageRatelimit time.Duration

const (
	// messageRatelimitNormal used in channels where bot has no special permissions
	messageRatelimitNormal = messageRatelimit(1250 * time.Millisecond)
	// messageRatelimitElevated used in channels where bot is a Moderator, VIP or Broadcaster
	messageRatelimitElevated = messageRatelimit(100 * time.Millisecond)
)
