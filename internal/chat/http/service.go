package ws

import (
	"net/http"
	"redis-ws/internal/chat"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Option func(*Service)

func WithRedis(r *redis.Client) Option {
	return func(s *Service) {
		s.br = chat.NewBroker(chat.WithRedis(r))
	}
}

type Service struct {
	mux chi.Router

	br chat.Broker
}

func New(opts ...Option) *Service {
	s := Service{
		mux: chi.NewRouter(),
	}

	for _, opt := range opts {
		opt(&s)
	}

	if s.br == nil {
		s.br = chat.NewBroker()
	}

	s.routes()
	return &s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) routes() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ID query param

		found, _ := s.br.Insert(r.Context(), "CHAT")
		found.ServeHTTP(w, r)
	})
}
