package linkstorage

import (
	"net/url"
	"sync"

	"golang.org/x/exp/rand"
)

type LinkStorage struct {
	links sync.Map
}

const LengthShortLink = 10
const NumberCharSet = 52

var charSet *[NumberCharSet]uint8 // 26 + 26

func generateChars() {
	charSet = new([NumberCharSet]uint8)
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

func generateShortLink() string {
	if charSet == nil {
		generateChars()
	}

	bytes := make([]byte, LengthShortLink)
	for i := 0; i < LengthShortLink; i++ {
		bytes[i] = charSet[byte(rand.Intn(NumberCharSet))]
	}
	return string(bytes)
}

func (t *LinkStorage) AddLink(originalLink string) (string, bool) {
	if originalLink == "" {
		return "", false
	}
	_, okParse := url.ParseRequestURI(originalLink)
	if okParse != nil {
		return "", false
	}

	for {
		shortLink := generateShortLink()
		_, ok := t.links.Load(shortLink)
		if !ok {
			t.links.Store(shortLink, originalLink)
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
