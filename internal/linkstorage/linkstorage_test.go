package linkstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {

	tests := []struct {
		name        string
		addLink     string
		wantAddOK   bool
		wantGetOK   bool
		wantGetLink string
	}{
		{
			name:        "Test 1",
			addLink:     "http://ya.ru",
			wantAddOK:   true,
			wantGetOK:   true,
			wantGetLink: "http://ya.ru",
		},

		{
			name:        "Test 2",
			addLink:     "/////", // Классная валидная сслыка
			wantAddOK:   true,
			wantGetOK:   true,
			wantGetLink: "/////",
		},

		{
			name:        "Negative 1",
			addLink:     "123 GO",
			wantAddOK:   false,
			wantGetOK:   false,
			wantGetLink: "",
		},

		{
			name:        "Negative 2",
			addLink:     "",
			wantAddOK:   false,
			wantGetOK:   false,
			wantGetLink: "",
		},
	}

	storage := new(LinkStorage)
	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			short, ok1 := storage.AddLink(test.addLink)
			assert.Equal(t, ok1, test.wantAddOK)
			link, ok2 := storage.GetLink(short)
			assert.Equal(t, ok2, test.wantGetOK)
			assert.Equal(t, link, test.wantGetLink)
		})
	}
}
