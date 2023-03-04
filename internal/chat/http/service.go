package ws

import (
	"net/http"
	"redis-ws/pkg/websocket"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	mux chi.Router
}

func New() *Service {
	s := Service{
		mux: chi.NewRouter(),
	}
	s.routes()
	return &s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) routes() {
	// connect to websocket here
	s.mux.HandleFunc("/", websocket.NewClient(0).ServeHTTP)
}
