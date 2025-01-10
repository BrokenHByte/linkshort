package linkstorage

import (
	"sync"

	"golang.org/x/exp/rand"
)

const LengthShortLink = 10

type LinkStorage struct {
	links sync.Map
}

var charSet [52]uint8 // 26 + 26

func init() {
	// UTF-8, but valid
	index := 0
	for i := 'A'; i <= 'Z'; i++ {
		charSet[index] = uint8(i)
		index++
	}
	for i := 'a'; i <= 'z'; i++ {
		charSet[index] = uint8(i)
		index++
	}
}

func NewLinkStorage() *LinkStorage {
	var g LinkStorage
	return &g
}

func generateShortLink() string {
	bytes := make([]byte, LengthShortLink)
	for i := 0; i < LengthShortLink; i++ {
		bytes[i] = charSet[byte(rand.Intn(52))]
	}
	return string(bytes)
}

func (t *LinkStorage) AddLink(originalLink string, prefix string) string {
	for {
		shortLink := prefix + generateShortLink()
		_, ok := t.links.Load(shortLink)
		if !ok {
			t.links.Store(shortLink, originalLink)
			return shortLink
		}
	}
}

func (t *LinkStorage) GetLink(shortLink string) (string, bool) {
	originalLink, ok := t.links.Load(shortLink)
	for !ok {
		return "", false
	}
	return originalLink.(string), true
}
