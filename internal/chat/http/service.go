package ws

import (
	"net/http"
	"redis-ws/pkg/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	mux chi.Router

	ps *redis.Client
}

func New(redis *redis.Client) *Service {
	s := Service{
		mux: chi.NewRouter(),
		ps:  redis,
	}
	s.routes()
	return &s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) routes() {
	// connect to websocket here
	s.mux.HandleFunc("/", websocket.NewClient(0, s.ps).ServeHTTP)
}
