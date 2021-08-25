package eventsub

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zneix/zneixbot-twitch/pkg/bot"
)

func mountRouter(r *chi.Mux) *chi.Mux {
	if baseURL == "" {
		log.Printf("Listening on %s (Prefix=%s, BaseURL=%s)\n", bindAddress, listenPrefix, baseURL)
		return r
	}

	// figure out prefix from address
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		log.Fatal("Scheme must be included in base url")
	}

	listenPrefix = u.Path
	// Empty prefix can't be passed to chi
	if listenPrefix == "" {
		listenPrefix = "/"
	}
	ur := chi.NewRouter()
	ur.Mount(listenPrefix, r)

	log.Printf("Listening on %s (Prefix=%s, BaseURL=%s)\n", bindAddress, listenPrefix, baseURL)
	return ur
}

func listen(bind string, router *chi.Mux) {
	srv := &http.Server{
		Handler:      router,
		Addr:         bind,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// InitializeWebServer creates new http.Server which will handle and respond to webhook events sent by Twitch
func InitializeWebServer(zb *bot.Bot, serverQuit chan struct{}) {
	defer close(serverQuit)

	router := chi.NewRouter()

	handleMainRoutes(router, zb)
	listen(bindAddress, mountRouter(router))
}
