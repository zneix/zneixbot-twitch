package utils

import (
	"log"
	"os"

	. "github.com/zneix/zniksbot/bot"
)

const (
	envPrefix = "ZNIKSBOT_"
)

func GetEnv(envName string, isRequired bool) (value string, exists bool) {
	value, envExists := os.LookupEnv(envPrefix + envName)

	if !envExists && isRequired {
		log.Fatalf("Missing required %s environment variable", envPrefix+envName)
	}

	return value, envExists
}

func SendMessage(target string, message string) {
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

	if Zniksbot.LastMsgs[target] == message {
		message += " \U000E0000"
	}

	Zniksbot.Client.Say(target, message)
	Zniksbot.LastMsgs[target] = message
}
