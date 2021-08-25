package eventsub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/nicklaw5/helix"
	"github.com/zneix/zneixbot-twitch/pkg/bot"
	"github.com/zneix/zneixbot-twitch/pkg/utils"
)

var (
	zb *bot.Bot
)

type eventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is the public zneixbot's API, but most of the endpoints are (and will be) undocumented ThreeLetterAPI TeaTime\nMore information on the GitHub repo: https://github.com/zneix/zneixbot-twitch"))
}

func health(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memory := fmt.Sprintf("Alloc=%v MiB, TotalAlloc=%v MiB, Sys=%v MiB, NumGC=%v",
		m.Alloc/1024/1024,
		m.TotalAlloc/1024/1024,
		m.Sys/1024/1024,
		m.NumGC)

	w.Write([]byte(fmt.Sprintf("Uptime: %s\nMemory: %s", utils.TimeSince(zb.StartTime), memory)))
}

func eventSubCallback(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body in eventSubCallback: " + err.Error())
		return
	}
	defer r.Body.Close()

	// First of all, check if the message really came from Twitch by verifying the signature
	if !helix.VerifyEventSubNotification(secret, r.Header, string(body)) {
		log.Println("Received a notification, but the signature was invalid")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Read data sent in the request
	var notification eventSubNotification
	//err = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
	err = json.Unmarshal(body, &notification)
	if err != nil {
		log.Printf("Error unmarshaling incoming eventsub message: %s, request body: %s\n", err, string(body))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If a challenge is specified in request, respond to it
	if notification.Challenge != "" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(notification.Challenge))
		return
	}

	handleIncomingNotification(notification)
	w.WriteHeader(http.StatusOK)
}

func handleMainRoutes(router *chi.Mux, initializedBot *bot.Bot) {
	zb = initializedBot

	router.Get("/", index)
	router.Get("/health", health)
	router.Post("/eventsubcallback", eventSubCallback)
}
