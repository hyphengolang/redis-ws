package ws

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed pages
var pages embed.FS

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
	fs := http.FileServer(http.FS(pages))
	s.mux.Handle("/*", fs)
}
