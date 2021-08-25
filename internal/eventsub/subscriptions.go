package eventsub

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nicklaw5/helix"
	"github.com/zneix/zneixbot-twitch/pkg/bot"
)

func CreateChannelSubscription(zb *bot.Bot, subscription *bot.ChannelEventSubSubscription, channelID string) (err error) {
	resp, err := zb.Helix.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    subscription.Type,
		Version: subscription.Version,
		Condition: helix.EventSubCondition{
			BroadcasterUserID: channelID,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: getCallbackString(),
			Secret:   secret,
		},
	})
	if err != nil {
		return err
	}

	log.Printf("Create subscription response for %# v: %# v\n", subscription, resp.Data)

	// TODO: Properly handle pending status
	//subscriptionsPending = append(subscriptionsPending, sub.ID)

	return nil
}

func handleIncomingNotification(notification eventSubNotification) {
	// TODO: Export these separate types to separate functions if we'll consider handling more notification types
	switch notification.Subscription.Type {
	case helix.EventSubTypeChannelFollow:
		var notificationEvent helix.EventSubChannelFollowEvent
		err := json.Unmarshal(notification.Event, &notificationEvent)
		if err != nil {
			log.Printf("Error unmarshaling received EventSub notification event: %s, data: %s\n", err, string(notification.Event))
			return
		}

		channel := zb.Channels[notificationEvent.BroadcasterUserID]

		channel.Send(fmt.Sprintf("Woah %s, thank you so much FeelsDankMan 👉 <3", notificationEvent.UserLogin))

	default:
		log.Printf("Unhandled EventSub notification: %# v\n", notification)
	}
}
