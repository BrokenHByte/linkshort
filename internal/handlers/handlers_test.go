package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/BrokenHByte/linkshort/internal/config"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockStorage struct {
	mock.Mock
}

func (t *MockStorage) AddLink(originalLink string) (string, bool) {
	if originalLink == "" {
		return "", false
	}
	t.On("GetLink", "/test123").Return(originalLink)
	return "test123", true
}

func (t *MockStorage) GetLink(shortLink string) (string, bool) {
	args := t.Called(shortLink)
	return args.String(0), true
}

func TestHandlersCreateAndGet(t *testing.T) {
	tests := []struct {
		name           string
		request        string
		link           string
		wantFailShort  bool
		wantShortLink  string
		wantFullLink   string
		wantStatusCode int
	}{
		{
			name:    "Test 1",
			request: "http://localhost:8080",
			link:    "https://yandex.ru",

			wantShortLink:  "http://localhost:8080/test123",
			wantFullLink:   "https://yandex.ru",
			wantStatusCode: 201,
		},

		{
			name:    "Negative Test 1",
			request: "http://localhost:8080",
			link:    "strange text",

			wantShortLink:  "http://localhost:8080/test123",
			wantFullLink:   "strange text",
			wantStatusCode: 400,
		},

		{
			name:    "Negative Test 2",
			request: "http://localhost:8080",
			link:    "",

			wantShortLink:  "",
			wantFullLink:   "",
			wantStatusCode: 400,
		},
	}

	config := &config.ServerConfig{ServerAddr: ":8080", BaseURL: "http://localhost:8080"}
	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			storage := new(MockStorage)

			// Отправляем на сохранение ссылку
			reader := strings.NewReader(test.link)
			requestPost := httptest.NewRequest(http.MethodPost, test.request, reader)
			wSend := httptest.NewRecorder()

			h := NewHandlers(config, storage)
			h.HandleCreateShortLink(wSend, requestPost)
			result := wSend.Result()

			shortURL := ""
			require.Equal(t, test.wantStatusCode, result.StatusCode)
			if result.StatusCode == 201 {
				bodyData, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)

				shortURL = string(bodyData)
				assert.Equal(t, test.wantShortLink, shortURL)
			} else {
				return
			}

			parseURL, err := url.Parse(shortURL)
			if err != nil {
				return
			}

			wGet := httptest.NewRecorder()
			requestGet := httptest.NewRequest(http.MethodGet, "/{shortLink}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shortLink", parseURL.Path)
			requestGet = requestGet.WithContext(context.WithValue(requestGet.Context(), chi.RouteCtxKey, rctx))

			h.HandleGetFullLink(wGet, requestGet)
			assert.Equal(t, wGet.Header().Get("Location"), test.wantFullLink)

			storage.AssertExpectations(t)
		})
	}
}
