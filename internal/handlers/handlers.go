package handlers

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/BrokenHByte/linkshort/internal/config"
	"github.com/BrokenHByte/linkshort/internal/logs"
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

func (t *Handlers) HandleCreateShortLink() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		link, err := io.ReadAll(r.Body)
		if err != nil || len(link) == 0 {
			http.Error(w, "body link error", http.StatusBadRequest)
			return
		}

		newURL := string(link)
		_, err = url.ParseRequestURI(newURL)
		if err != nil {
			http.Error(w, "url error", http.StatusBadRequest)
			return
		}

		shortLink, ok := t.storage.AddLink(newURL)
		if !ok {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(t.config.BaseURL + "/" + shortLink))
	})
}

func (t *Handlers) HandleGetFullLink() http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Получаем короткую ссылку, если валидна, возвращаем редирект
		shortLink := chi.URLParam(r, "shortLink")
		originalLink, ok := t.storage.GetLink(shortLink)
		if !ok {
			http.Error(rw, "invalid short link", http.StatusBadRequest)
			return
		}
		rw.Header().Set("Location", originalLink)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	})
}

type RequestShortenJSON struct {
	URL string `json:"url"`
}

type ResponseShortenJSON struct {
	Result string `json:"result"`
}

func (t *Handlers) HandleShortenJSON() http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			logs.Logs().Infoln("Content-Type not equal application/json")
		}

		jsonBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			logs.Logs().Infoln("Invalid json request body")
			return
		}

		var requestJSONObject RequestShortenJSON
		if err = json.Unmarshal(jsonBytes, &requestJSONObject); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			logs.Logs().Infoln("Error unmarshal - ", jsonBytes)
			return
		}

		shortLink, ok := t.storage.AddLink(requestJSONObject.URL)
		if !ok {
			http.Error(rw, "", http.StatusBadRequest)
			logs.Logs().Infoln("Error AddLink - ", requestJSONObject.URL)
			return
		}

		responseJSONObject := ResponseShortenJSON{t.config.BaseURL + "/" + shortLink}
		responseJSONBytes, err := json.Marshal(responseJSONObject)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			logs.Logs().Infoln("Error marshal - ", responseJSONObject)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(responseJSONBytes)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (t *Handlers) MiddlewareReadGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Распаковываем вывод если есть метка в header
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") || r.Header.Get("Content-Type") != "application/x-gzip" {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			logs.Logs().Error(err.Error())
			return
		}
		defer gz.Close()

		r.Body = gz
		next.ServeHTTP(w, r)
	})
}

func (t *Handlers) MiddlewareWriteGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Сжимаем вывод только если есть поддержка у клиента, и на выход идёт текст или json
		content := r.Header.Get("Content-Type")
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") || !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") ||
			(content != "application/json" && content != "text/html") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
