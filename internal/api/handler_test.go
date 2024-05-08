package api

import (
	"bytes"
	"encoding/json"
	"github.com/Azzonya/go-shortener/internal/entities"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Azzonya/go-shortener/internal/repo/inmemory"
	shortener_service "github.com/Azzonya/go-shortener/internal/shortener"
)

func TestRest_ShortenJSON(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		testURL     string
	}

	tests := []struct {
		name          string
		rest          Rest
		request       string
		requestMethod string
		want          want
	}{
		{
			name:          "test ShortenJSON",
			request:       "/",
			requestMethod: http.MethodPost,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				testURL:     "www.example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := inmemory.New("/tmp/short-url-repo.json")
			require.NoError(t, err)

			tt.rest.shortener = shortener_service.New("http://localhost:8080", repo)

			r := gin.Default()
			r.POST(tt.request, tt.rest.ShortenJSON)

			reqBody := &Request{
				URL: tt.want.testURL,
			}

			requestBody, err := json.Marshal(reqBody)
			if err != nil {
				panic(err)
			}

			request := httptest.NewRequest(tt.requestMethod, tt.request, bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			//tt.rest.HShortenerURL(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			responseBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			repObj := Response{}

			err = json.Unmarshal(responseBody, &repObj)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			shortURL := repObj.Result

			parts := strings.Split(shortURL, "/")

			originalURL, exist := tt.rest.shortener.GetOneByShortURL(parts[len(parts)-1])
			if !exist {
				require.Fail(t, "Expected short URL in urlMap")
			}

			assert.Equal(t, tt.want.testURL, originalURL)
		})
	}
}

func TestRest_Shorten(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		testURL     string
	}

	tests := []struct {
		name          string
		rest          Rest
		request       string
		requestMethod string
		want          want
	}{
		{
			name:          "1st test",
			request:       "/",
			requestMethod: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				testURL:     "www.example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := inmemory.New("/tmp/short-url-repo.json")
			require.NoError(t, err)

			tt.rest.shortener = shortener_service.New("http://localhost:8080", repo)

			r := gin.Default()
			r.POST(tt.request, tt.rest.Shorten)

			request := httptest.NewRequest(tt.requestMethod, tt.request, strings.NewReader(tt.want.testURL))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			//tt.rest.HShortenerURL(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			reqObj, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			shortURL := string(reqObj)

			parts := strings.Split(shortURL, "/")

			originalURL, exist := tt.rest.shortener.GetOneByShortURL(parts[len(parts)-1])
			if !exist {
				require.Fail(t, "Expected short URL in urlMap")
			}

			assert.Equal(t, tt.want.testURL, originalURL)
		})
	} //
}

func TestRest_Redirect(t *testing.T) {
	type want struct {
		location   string
		statusCode int
	}

	tests := []struct {
		name          string
		rest          Rest
		requestMethod string
		want          want
	}{
		{
			name:          "1st test",
			rest:          Rest{},
			requestMethod: http.MethodGet,
			want: want{
				location:   "https://www.youtube.com",
				statusCode: http.StatusTemporaryRedirect,
			},
		},
		{
			name:          "2nd test",
			rest:          Rest{},
			requestMethod: http.MethodGet,
			want: want{
				location:   "https://www.google.com",
				statusCode: http.StatusTemporaryRedirect,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			r.GET("/:id", tt.rest.Redirect)

			testShortURL := "Abcdefgh"

			repo, err := inmemory.New("/tmp/short-url-repo.json")
			require.NoError(t, err)

			tt.rest.shortener = shortener_service.New("http://localhost:8080", repo)

			err = repo.Add(tt.want.location, testShortURL, "")
			require.NoError(t, err)

			request := httptest.NewRequest(tt.requestMethod, "/"+testShortURL, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

func TestRest_ShortenURLs(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name          string
		rest          Rest
		request       string
		requestMethod string
		want          want
	}{
		{
			name:          "test ShortenURLs",
			request:       "/api/shorten/batch",
			requestMethod: http.MethodPost,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := inmemory.New("/tmp/short-url-repo.json")
			require.NoError(t, err)

			tt.rest.shortener = shortener_service.New("http://localhost:8080", repo)

			r := gin.Default()
			r.POST(tt.request, tt.rest.ShortenURLs)

			reqBody := []*entities.ReqURL{
				{
					OriginalURL: "www.example.com",
				},
				{
					OriginalURL: "www.example2.com",
				},
			}

			requestBody, err := json.Marshal(reqBody)
			if err != nil {
				panic(err)
			}

			request := httptest.NewRequest(tt.requestMethod, tt.request, bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			responseBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			repObj := []*entities.ReqURL{}

			err = json.Unmarshal(responseBody, &repObj)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestRest_ListAll(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name          string
		rest          Rest
		request       string
		requestMethod string
		want          want
	}{
		{
			name:          "test ListAll",
			request:       "/api/user/urls",
			requestMethod: http.MethodGet,
			want: want{
				contentType: "application/json; charset=utf-8",
				statusCode:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := inmemory.New("/tmp/short-url-repo.json")
			require.NoError(t, err)

			tt.rest.shortener = shortener_service.New("http://localhost:8080", repo)

			r := gin.Default()
			r.GET(tt.request, tt.rest.ListAll)

			request := httptest.NewRequest(tt.requestMethod, tt.request, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

		})
	}
}

func TestRest_DeleteURLs(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name          string
		rest          Rest
		request       string
		requestMethod string
		want          want
	}{
		{
			name:          "test DeleteURLs",
			request:       "/api/user/urls",
			requestMethod: http.MethodDelete,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := inmemory.New("/tmp/short-url-repo.json")
			require.NoError(t, err)

			tt.rest.shortener = shortener_service.New("http://localhost:8080", repo)

			r := gin.Default()
			r.POST("/api/shorten/batch", tt.rest.ShortenURLs)
			r.DELETE(tt.request, tt.rest.DeleteURLs)

			reqBody := []*entities.ReqURL{
				{
					OriginalURL: "www.example.com",
				},
				{
					OriginalURL: "www.example2.com",
				},
			}

			requestBody, err := json.Marshal(reqBody)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, http.StatusCreated, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			responseBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			repObj := []*entities.ReqURL{}

			err = json.Unmarshal(responseBody, &repObj)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			URLs := []string{"example.com"}

			responseBodyDelete, err := json.Marshal(URLs)
			require.NoError(t, err)

			requestDelete := httptest.NewRequest(tt.requestMethod, tt.request, bytes.NewBuffer(responseBodyDelete))

			v := httptest.NewRecorder()
			r.ServeHTTP(v, requestDelete)

			resultDelete := v.Result()
			err = resultDelete.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resultDelete.StatusCode)

		})
	}
}

func TestRest_Ping(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name          string
		request       string
		requestMethod string
		rest          Rest
		want          want
	}{
		{
			name:          "test Ping",
			request:       "/ping",
			requestMethod: "GET",
			want: want{
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := inmemory.New("/tmp/short-url-repo.json")
			require.NoError(t, err)

			tt.rest.shortener = shortener_service.New("http://localhost:8080", repo)

			r := gin.Default()
			r.GET(tt.request, tt.rest.Ping)

			request := httptest.NewRequest(tt.requestMethod, tt.request, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
