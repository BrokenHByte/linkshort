package main

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/BrokenHByte/linkshort/internal/config"
	"github.com/BrokenHByte/linkshort/internal/flags"
	"github.com/BrokenHByte/linkshort/internal/linkstorage"
	"github.com/go-chi/chi"
	"golang.org/x/exp/rand"
)

var lStorage = linkstorage.NewLinkStorage()
var configServer *config.ServerConfig

func HandleCreateShortLink(rw http.ResponseWriter, r *http.Request) {
	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		http.Error(rw, "body link error", http.StatusBadRequest)
		return
	}

	newURL := string(link)
	_, err = url.Parse(newURL)
	if err != nil {
		http.Error(rw, "url error", http.StatusBadRequest)
		return
	}

	shortLink := lStorage.AddLink(newURL, "")
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(configServer.PrefixLink + "/" + shortLink))
}

func HandleGetFullLink(rw http.ResponseWriter, r *http.Request) {
	// Получаем короткую ссылку, если валидна, возвращаем редирект
	shortLink := chi.URLParam(r, "shortLink")
	originalLink, ok := lStorage.GetLink(shortLink)
	if !ok {
		http.Error(rw, "invalid short link", http.StatusBadRequest)
		return
	}
	rw.Header().Set("Location", originalLink)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))
	configServer = config.NewServerConfig()
	flags.ParseFlags(configServer)

	r := chi.NewRouter()
	r.Post("/", HandleCreateShortLink)
	r.Get("/{shortLink}", HandleGetFullLink)

	err := http.ListenAndServe(configServer.ServerAddr, r)
	if err != nil {
		panic("The server address is inaccessible or not valid")
	}
}
