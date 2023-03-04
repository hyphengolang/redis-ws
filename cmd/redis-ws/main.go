package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	chatHTTP "redis-ws/internal/chat/http"
	uiHTTP "redis-ws/internal/ui/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

var port int

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// port as flag
	flag.IntVar(&port, "port", 8080, "port to listen on")

	flag.Parse()
}

var ctx = context.Background()

func run() error {
	r, err := newRedisConn(ctx, ":6379")
	if err != nil {
		return err
	}
	defer r.Close()

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Mount("/", newUIService())
	mux.Mount("/v1/chats", newChatService(r))

	addr := fmt.Sprintf(":%d", port)

	log.Printf("listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func newRedisConn(ctx context.Context, addr string) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return r, r.Ping(ctx).Err()
}

// baseUrl is needed here
func newUIService() http.Handler {
	return uiHTTP.New()
}

func newChatService(r *redis.Client) http.Handler {
	return chatHTTP.New(r)
}
