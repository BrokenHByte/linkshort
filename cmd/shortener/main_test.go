package main

type want struct {
	statusCode int
}

/*
func TestHandleCreateShortLink(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		request string
		link    string
		want    want
	}{
		{
			name:    "Test link 1",
			method:  http.MethodPost,
			request: "http://localhost:8080",
			link:    "https://yandex.ru",
			want: want{
				statusCode: 201,
			},
		},

		{
			name:    "Test link 2",
			method:  http.MethodPost,
			request: "http://localhost:8080",
			link:    "Secret word",
			want: want{
				statusCode: 201,
			},
		},

		{
			name:    "Negative Method 1",
			method:  http.MethodGet,
			request: "http://localhost:8080/abs",
			link:    "https://yandex.ru",
			want: want{
				statusCode: 400,
			},
		},

		{
			name:    "Negative Method 2",
			method:  http.MethodPost,
			request: "http://localhost:8080",
			link:    "",
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			// Отправляем на сохранение ссылку
			reader := strings.NewReader(test.link)
			request := httptest.NewRequest(test.method, test.request, reader)
			wSend := httptest.NewRecorder()
			HandleCreateShortLink(wSend, request)
			result := wSend.Result()

			smallURL := ""
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			if result.StatusCode == 201 {
				userResult, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				smallURL = string(userResult)
				if assert.Greater(t, len(smallURL), 1) {
					_, err = url.Parse(smallURL)
					assert.Nil(t, err)
				}
			}
		})
	}
}

func TestHandleConvertToFullLink(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		request string
		link    string
		want    want
	}{
		{
			name:    "Test link 1",
			method:  http.MethodGet,
			request: "http://localhost:8080",
			link:    "https://yandex.ru",
			want: want{
				statusCode: 307,
			},
		},

		{
			name:    "Test link 2",
			method:  http.MethodGet,
			request: "http://localhost:8080",
			link:    "Secret word",
			want: want{
				statusCode: 307,
			},
		},

		{
			name:    "Negative Method 2",
			method:  http.MethodGet,
			request: "http://localhost:8080",
			link:    "",
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			// Отправляем на сохранение ссылку
			reader := strings.NewReader(test.link)
			request := httptest.NewRequest(http.MethodPost, test.request, reader)
			wSend := httptest.NewRecorder()
			HandleCreateShortLink(wSend, request)
			result := wSend.Result()

			smallURL := "http://localhost:8080"
			if result.StatusCode == 201 {
				userResult, err := io.ReadAll(result.Body)
				if err == nil {
					smallURL = string(userResult)
					result.Body.Close()
				}
			}

			// Отправляем сохранённую сжатую ссылку
			reader2 := strings.NewReader(smallURL)
			request2 := httptest.NewRequest(test.method, smallURL, reader2)
			wGet := httptest.NewRecorder()
			HandleGetFullLink(wGet, request2)
			result = wGet.Result()
			result.Body.Close()
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			if result.StatusCode == 307 {
				location := result.Header.Get("Location")
				assert.Equal(t, test.link, location)
			}
		})
	}
}*/
