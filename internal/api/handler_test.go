package api

import (
	"github.com/Azzonya/go-shortener/internal/inmemory"
	"github.com/Azzonya/go-shortener/internal/repo/pg"
	shortener_service "github.com/Azzonya/go-shortener/internal/shortener"
	"github.com/Azzonya/go-shortener/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				testURL:     "www.example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stor, err := inmemory.NewStorage("/tmp/short-url-repo.json", false)
			require.NoError(t, err)

			db, err := pkg.InitDatabasePg("postgresql://postgres:postgres@localhost:5432/postgres")
			require.NoError(t, err)

			repo := pg.New(db)

			tt.rest.shortener = shortener_service.New("http://localhost:8080", stor, repo, "", false)

			r := gin.Default()
			r.POST(tt.request, tt.rest.Shorten)

			request := httptest.NewRequest(tt.requestMethod, tt.request, strings.NewReader(tt.want.testURL))
			request.AddCookie(&http.Cookie{
				Name:     "userID",
				Value:    "1",
				HttpOnly: false,
			})

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			r.GET("/:id", tt.rest.Redirect)

			testShortURL := "Abcdefgh"

			stor, err := inmemory.NewStorage("/tmp/short-url-repo.json", false)
			require.NoError(t, err)

			db, err := pkg.InitDatabasePg("postgresql://postgres:postgres@localhost:5432/postgres")
			require.NoError(t, err)

			repo := pg.New(db)

			tt.rest.shortener = shortener_service.New("http://localhost:8080", stor, repo, "", false)

			err = stor.Add(testShortURL, tt.want.location)
			require.NoError(t, err)

			request := httptest.NewRequest(tt.requestMethod, "/"+testShortURL, nil)
			request.AddCookie(&http.Cookie{
				Name:     "userID",
				Value:    "1",
				HttpOnly: false,
			})

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
