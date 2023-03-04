package main

import (
	"log"
	"net/http"

	chatHTTP "redis-ws/internal/chat/http"
	uiHTTP "redis-ws/internal/ui/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func run() error {
	// 1. run html frontend

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Mount("/", newUIService())
	mux.Mount("/v1/chats", newChatService())

	log.Println("redis-ws starting")
	return http.ListenAndServe(":8080", mux)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func newUIService() http.Handler {
	return uiHTTP.New()
}

func newChatService() http.Handler {
	return chatHTTP.New()
}
