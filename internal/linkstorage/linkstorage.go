package linkstorage

import "golang.org/x/exp/rand"

const LengthShortLink = 10

type LinkStorage struct {
	links map[string]string
}

func NewLinkStorage() *LinkStorage {
	var g LinkStorage
	g.links = make(map[string]string)
	return &g
}

func generateShortLink() string {
	bytes := make([]byte, LengthShortLink)
	for i := 0; i < LengthShortLink; i++ {
		bytes[i] = byte(65 + rand.Intn(90-65))
	}
	return string(bytes)
}

func (t *LinkStorage) AddLink(originalLink string) string {
	for {
		shortLink := generateShortLink()
		_, ok := t.links[shortLink]
		if !ok {
			t.links[shortLink] = originalLink
			return shortLink
		}
	}
}

func (t *LinkStorage) GetLink(shortLink string) (string, bool) {
	originalLink, ok := t.links[shortLink]
	for !ok {
		return "", false
	}
	return originalLink, true
}
