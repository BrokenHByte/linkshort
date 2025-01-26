package server

import (
	"net/http"

	"github.com/BrokenHByte/linkshort/internal/config"
	"github.com/BrokenHByte/linkshort/internal/handlers"
	"github.com/BrokenHByte/linkshort/internal/logs"
	"github.com/go-chi/chi"
)

type Server struct {
	handlers *handlers.Handlers
	config   *config.ServerConfig
}

func RunServer(handlers *handlers.Handlers, config *config.ServerConfig) *Server {
	server := Server{handlers, config}
	r := chi.NewRouter()
	r.Use(handlers.MiddlewareReadGzip)
	r.Use(handlers.MiddlewareWriteGzip)
	r.Use(logs.LoggingRequest)
	r.Post("/", handlers.HandleCreateShortLink())
	r.Post("/api/shorten", handlers.HandleShortenJSON())
	r.Get("/{shortLink}", handlers.HandleGetFullLink())

	err := http.ListenAndServe(config.ServerAddr, r)
	if err != nil {
		panic("The server address is inaccessible or not valid")
	}
	return &server
}
