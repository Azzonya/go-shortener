package server

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRest_HShortenerURL_HShortener(t *testing.T) {
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
			tt.rest.urlMap = make(map[string]string)

			r := gin.Default()
			r.POST(tt.request, tt.rest.HShortener)

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

			originalURL, exist := tt.rest.urlMap[parts[len(parts)-1]]
			if !exist {
				require.Fail(t, "Expected short URL in urlMap", tt.rest.urlMap)
			}

			assert.Equal(t, tt.want.testURL, originalURL)
		})
	} //
}

func TestRest_HShortenerURL_HRedirect(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			r.GET("/:id", tt.rest.HRedirect)

			testShortURL := "Abcdefgh"
			tt.rest.urlMap = make(map[string]string)
			tt.rest.urlMap[testShortURL] = tt.want.location

			request := httptest.NewRequest(tt.requestMethod, "/"+testShortURL, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
