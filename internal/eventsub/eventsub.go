package eventsub

import (
	"log"
	"strings"

	"github.com/zneix/zneixbot-twitch/pkg/utils"
)

var (
	subscriptionsPending []string

	// secret EventSub secret used while creating subscriptions
	// can be whatever but must be between 10 and 100 characters
	secret string

	// listenPrefix ...
	listenPrefix string

	// baseURL
	baseURL string

	// bindAddress on which the HTTP server will listen on
	bindAddress string
)

func init() {
	log.Println("Initializing EventSub...")

	secret, _ = utils.GetEnv("TWITCH_EVENTSUB_SECRET", true)
	baseURL, _ = utils.GetEnv("BASE_URL", true)

	var bindAddressExists bool
	bindAddress, bindAddressExists = utils.GetEnv("BIND_ADDRESS", false)
	if !bindAddressExists {
		bindAddress = ":2557"
	}
}

func getCallbackString() string {
	return strings.TrimSuffix(baseURL, "/") + "/eventsubcallback"
}
