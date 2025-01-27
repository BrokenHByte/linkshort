package server

import (
	"net/http"
	"sync"

	"github.com/BrokenHByte/linkshort/internal/config"
	"github.com/BrokenHByte/linkshort/internal/handlers"
	"github.com/BrokenHByte/linkshort/internal/logs"
	"github.com/go-chi/chi"
)

func RunServer(wg *sync.WaitGroup, handlers *handlers.Handlers, config *config.ServerConfig) *http.Server {
	r := chi.NewRouter()
	r.Use(handlers.MiddlewareReadGzip)
	r.Use(handlers.MiddlewareWriteGzip)
	r.Use(logs.LoggingRequest)
	r.Post("/", handlers.HandleCreateShortLink())
	r.Post("/api/shorten", handlers.HandleShortenJSON())
	r.Get("/{shortLink}", handlers.HandleGetFullLink())

	srv := &http.Server{Addr: config.ServerAddr, Handler: r}

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logs.Logs().Fatalf("ListenAndServe(): %v", err)
		}
	}()
	return srv
}
