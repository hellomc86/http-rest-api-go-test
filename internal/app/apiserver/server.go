package apiserver

import (
	"log/slog"

	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	logger *slog.Logger
}

func newServer(router *mux.Router, logger *slog.Logger) *server {
	s := &server{
		router: router,
		logger: logger,
	}

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
