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
	TwitchIRC *twitch.Client
	Mongo     *db.Connection
	Logins    map[string]string
	Channels  map[string]*Channel
	Commands  map[string]*Command
	Users     map[string]*User
	Self      *Self
	StartTime time.Time
}

type Channel struct {
	ID    string      `bson:"id"`
	Login string      `bson:"login"`
	Mode  ChannelMode `bson:"mode"`

	LastMsg      string
	QueueChannel chan *QueuedMessage
	Cooldowns    map[string]time.Time
}

type Command struct {
	Name        string
	Description string
	Usage       string
	Permission  Permission
	Cooldown    time.Duration
	Run         func(msg twitch.PrivateMessage, args []string) (err CommandError)
}

type User struct {
	ID         string     `bson:"id"`
	Permission Permission `bson:"permission"`
}

type QueuedMessage struct {
	Message string
}

// BotType represents which kind of global rate-limit the bot has
type BotType int

const (
	BotTypeNormal BotType = iota
	BotTypeKnown
	BotTypeVerified
)

// ChannelMode ...
type ChannelMode int

const (
	ChannelModeNormal ChannelMode = iota
	ChannelModeInactive
	ChannelModeVIP
	ChannelModeModerator

	ChannelModeBoundary
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

type Permission int32

const (
	PermissionNone  Permission = 0
	PermissionAdmin Permission = 1 << (iota - 1)
	PermissionSkipCooldown
)

// Flag flags that can be set with the set command and update corresponding properties in the database
type Flag int

const (
	FlagChannelMode Flag = iota
)

var FlagMap = map[string]Flag{
	"channel-mode": FlagChannelMode,
}

// CommandError is a type for
type CommandError int

const (
	// CommandErrorNoError command returned successfully
	CommandErrorNoError CommandError = iota

	// CommandErrorInvalidArguments user used invalid syntax
	// Could be caused by providing invalid arguments or too little of them
	CommandErrorInvalidArguments

	// CommandErrorInternal something completely unexpected happened, e.g. parsing new http.Req
	CommandErrorInternal

	// CommandErrorHTTPRequest executing an HTTP request failed, or any other HTTP request-related action
	CommandErrorHTTPRequest

	// CommandErrorReply something related to user's input failed
	// Response should be sent to the user in the Command.Run right before the return
	CommandErrorReply

	// CommandErrorMongo emits when a MongoDB operation failed / couldn't be performed
	CommandErrorMongo
)

func (err CommandError) String() (str string) {
	switch err {
	case CommandErrorNoError:
		return "channel-mode"
	default:
		return
	}
}
