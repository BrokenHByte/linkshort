package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/BrokenHByte/linkshort/internal/config"
	"github.com/go-chi/chi"
)

// Интерфейс для использования хранилища
type ActionsStorage interface {
	AddLink(originalLink string) (string, bool)
	GetLink(shortLink string) (string, bool)
}

type Handlers struct {
	config  *config.ServerConfig
	storage ActionsStorage
}

func NewHandlers(config *config.ServerConfig,
	storage ActionsStorage) *Handlers {
	return &Handlers{config, storage}
}

func (t *Handlers) HandleCreateShortLink(rw http.ResponseWriter, r *http.Request) {
	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		http.Error(rw, "body link error", http.StatusBadRequest)
		return
	}

	newURL := string(link)
	_, err = url.ParseRequestURI(newURL)
	if err != nil {
		http.Error(rw, "url error", http.StatusBadRequest)
		return
	}

	shortLink, ok := t.storage.AddLink(newURL)
	if !ok {
		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(t.config.BaseURL + "/" + shortLink))
}

func (t *Handlers) HandleGetFullLink(rw http.ResponseWriter, r *http.Request) {
	// Получаем короткую ссылку, если валидна, возвращаем редирект
	shortLink := chi.URLParam(r, "shortLink")
	originalLink, ok := t.storage.GetLink(shortLink)
	if !ok {
		http.Error(rw, "invalid short link", http.StatusBadRequest)
		return
	}
	rw.Header().Set("Location", originalLink)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
