package helixclient

import (
	"github.com/nicklaw5/helix"
	"github.com/zneix/zneixbot-twitch/pkg/utils"
)

// New returns a helix.Client that has requested an AppAccessToken and will keep it refreshed every 24h
func New() (*helix.Client, error) {
	twitchClientID, _ := utils.GetEnv("TWITCH_CLIENT_ID", true)
	twitchClientSecret, _ := utils.GetEnv("TWITCH_CLIENT_SECRET", true)

	client, err := helix.NewClient(&helix.Options{
		ClientID:     twitchClientID,
		ClientSecret: twitchClientSecret,
	})

	if err != nil {
		return nil, err
	}

	waitForFirstAppAccessToken := make(chan struct{})

	// Initialize methods responsible for refreshing access token
	go initAppAccessToken(client, waitForFirstAppAccessToken)
	<-waitForFirstAppAccessToken

	return client, nil
}
