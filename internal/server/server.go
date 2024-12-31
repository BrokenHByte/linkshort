package server

import (
	"net/http"

	"github.com/BrokenHByte/linkshort/internal/config"
	"github.com/BrokenHByte/linkshort/internal/handlers"
	"github.com/go-chi/chi"
)

type Server struct {
	handlers *handlers.Handlers
	config   *config.ServerConfig
}

func RunServer(handlers *handlers.Handlers, config *config.ServerConfig) *Server {
	server := Server{handlers, config}
	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) { handlers.HandleCreateShortLink(w, r) })
	r.Get("/{shortLink}", handlers.HandleGetFullLink)

	err := http.ListenAndServe(config.ServerAddr, r)
	if err != nil {
		panic("The server address is inaccessible or not valid")
	}
	return &server
}
