package ws

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
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
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			http.Error(w, "could not upgrade", http.StatusBadRequest)
			return
		}

		go func() {
			defer conn.Close()

			for {
				p, err := wsutil.ReadClientText(conn)
				if err != nil {
					return
				}
				err = wsutil.WriteServerText(conn, p)
				if err != nil {
					return
				}
			}
		}()
	})
}
