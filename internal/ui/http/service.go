package ws

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed pages/index.html
var indexPage string

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
	tmpl := template.Must(template.New("index").Parse(indexPage))

	s.mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Set the WebSocket URL based on the server's address and port
		webSocketURL := "ws://" + r.Host + "/v1/chats"
		tmpl.Execute(w, struct{ WebSocketURL string }{webSocketURL})
	})
}
