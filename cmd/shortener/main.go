package main

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/exp/rand"
)

const LengthShortLink = 10

type LinkStorage struct {
	links map[string]string
}

func NewLinkStorage() *LinkStorage {
	var g LinkStorage
	g.links = make(map[string]string)
	return &g
}

var linkStorage = NewLinkStorage()

func generateShortLink() string {
	bytes := make([]byte, LengthShortLink)
	for i := 0; i < LengthShortLink; i++ {
		bytes[i] = byte(65 + rand.Intn(90-65))
	}
	return string(bytes)
}

func (t *LinkStorage) addLink(originalLink string) string {
	for {
		shortLink := "/" + generateShortLink()
		_, ok := t.links[shortLink]
		if !ok {
			t.links[shortLink] = originalLink
			return shortLink
		}
	}
}

func (t *LinkStorage) getLink(shortLink string) (string, bool) {
	originalLink, ok := t.links[shortLink]
	for !ok {
		return "", false
	}
	return originalLink, true
}

func HandleRoute(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		HandleCreateShortLink(response, request)
	} else if request.Method == http.MethodGet {
		HandleConvertToFullLink(response, request)
	} else {
		http.Error(response, "", http.StatusBadRequest)
	}
}

func HandleCreateShortLink(response http.ResponseWriter, request *http.Request) {
	var link []byte
	link, err := io.ReadAll(request.Body)
	if err != nil || len(link) == 0 {
		http.Error(response, "body error", http.StatusBadRequest)
		return
	}

	newURL := string(link)
	_, err = url.Parse(newURL)
	if err != nil {
		http.Error(response, "url error", http.StatusBadRequest)
		return
	}

	shortLink := linkStorage.addLink(newURL)
	response.Header().Set("Content-Type", "text/plain")
	response.WriteHeader(201)
	response.Write([]byte("http://localhost:8080" + shortLink))
}

func HandleConvertToFullLink(response http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		http.Error(response, "", http.StatusBadRequest)
		return
	}

	id := request.URL.Path
	originalLink, ok := linkStorage.getLink(id)
	if !ok {
		http.Error(response, "", http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "text/plain")
	response.Header().Set("Location", originalLink)
	response.WriteHeader(307)
}

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, HandleRoute)
	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
