package linkstorage

import (
	"sync"

	"golang.org/x/exp/rand"
)

type LinkStorage struct {
	links sync.Map
}

const LengthShortLink = 10

var charSet []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateShortLink() string {
	bytes := make([]rune, LengthShortLink)
	for i := 0; i < LengthShortLink; i++ {
		bytes[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(bytes)
}

func (t *LinkStorage) AddLink(originalLink string) (string, bool) {
	if originalLink == "" {
		return "", false
	}

	for {
		shortLink := generateShortLink()
		_, loaded := t.links.LoadOrStore(shortLink, originalLink)
		if !loaded {
			return shortLink, true
		}
	}
}

func (t *LinkStorage) GetLink(shortLink string) (string, bool) {
	if shortLink == "" {
		return "", false
	}
	originalLink, ok := t.links.Load(shortLink)
	for !ok {
		return "", false
	}
	return originalLink.(string), true
}
