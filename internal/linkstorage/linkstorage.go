package linkstorage

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/BrokenHByte/linkshort/internal/logs"
	"golang.org/x/exp/rand"
)

type LinkStorage struct {
	links sync.Map
}

type externalDataStorage struct {
	Data []externalData `json:"data"`
}

type externalData struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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

func (t *LinkStorage) Load(dataReader io.Reader) {
	// Стираем всё
	t.links.Range(func(key, _ interface{}) bool {
		t.links.Delete(key)
		return true
	})

	bytesJSON, err := io.ReadAll(dataReader)
	if len(bytesJSON) == 0 || err != nil {
		return
	}
	var externalStorage externalDataStorage
	err = json.Unmarshal(bytesJSON, &externalStorage)
	if err != nil {
		logs.Logs().Infoln("Error marshal - ", bytesJSON)
		return
	}

	for _, value := range externalStorage.Data {
		t.links.Store(value.ShortURL, value.OriginalURL)
	}
}

func (t *LinkStorage) Save(dataWriter io.Writer) {

	length := 0
	t.links.Range(func(_, _ interface{}) bool {
		length++
		return true
	})

	storageData := externalDataStorage{Data: make([]externalData, length)}
	index := 0
	t.links.Range(func(k, v interface{}) bool {
		storageData.Data[index] = externalData{UUID: index, ShortURL: k.(string), OriginalURL: v.(string)}
		index++
		return true
	})

	bytesJSON, err := json.Marshal(storageData)
	if err != nil {
		logs.Logs().Infoln("Error marshal - ", bytesJSON)
		return
	}
	dataWriter.Write(bytesJSON)
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
