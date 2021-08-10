package bot

import (
	"log"
	"time"
)

func SendToChannel(zb *Bot, channelID string) {
	log.Printf("Starting routine for %s\n", channelID)
	var channel = zb.Channels[channelID]

	for message := range channel.QueueChannel {
		// Actually send the message to the chat
		zb.Client.Say(channel.Login, message.Message)

		// Update last sent message
		channel.LastMsg = message.Message

		// Wait for the cooldown
		time.Sleep(channel.Ratelimit)
	}
	log.Println("Done with routine for " + channelID)
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
