package api

import (
	"context"
	"github.com/Azzonya/go-shortener/internal/repo/inmemory"
	shortener_service "github.com/Azzonya/go-shortener/internal/shortener"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	repo, err := inmemory.New("/tmp/short-url-repo.json")
	require.NoError(t, err)

	shortener := shortener_service.New("http://localhost:8080", repo)

	type args struct {
		shortener *shortener_service.Shortener
		jwtSecret string
	}
	tests := []struct {
		name string
		args args
		want *Rest
	}{
		{
			name: "test New router",
			args: args{
				shortener: shortener,
				jwtSecret: "supersecret",
			},
			want: &Rest{
				shortener: shortener,
				jwtSecret: "supersecret",
				ErrorChan: make(chan error, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restTest := New(tt.args.shortener, tt.args.jwtSecret)
			assert.Equalf(t, tt.want.shortener, restTest.shortener, "New(%v, %v)", tt.args.shortener, tt.args.jwtSecret)
			assert.Equalf(t, tt.want.jwtSecret, restTest.jwtSecret, "New(%v, %v)", tt.args.shortener, tt.args.jwtSecret)
		})
	}
}

func TestRest_SetRouters(t *testing.T) {
	r := gin.Default()
	rest := &Rest{} // Замените на создание вашего объекта Rest

	// Вызов метода SetRouters для установки маршрутов
	rest.SetRouters(r)

	// Проверка наличия ожидаемых маршрутов в роутере
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/"},
		{"GET", "/:id"},
		{"POST", "/api/shorten"},
		{"GET", "/ping"},
		{"POST", "/api/shorten/batch"},
		{"GET", "/api/user/urls"},
		{"DELETE", "/api/user/urls"},
	}

	for _, route := range expectedRoutes {
		// Проверка наличия маршрута в роутере по заданному методу и пути
		routeInfo := r.Routes()
		assert.NotNil(t, routeInfo, "Route %s %s not found in router", route.method, route.path)
	}
}

func TestRest_Start(t *testing.T) {
	rest := &Rest{
		jwtSecret: "my_secret",
	}

	testCases := []struct {
		name       string
		listenAddr string
		pprofAddr  string
	}{
		{"Case 1", "localhost:8080", "localhost:6060"},
		{"Case 2", "127.0.0.1:8081", "127.0.0.1:6061"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rest.Start(tc.listenAddr, tc.pprofAddr)
			assert.NotNil(t, rest.pprofServer)
			assert.NotNil(t, rest.server)
		})
	}
}

func TestRest_Stop(t *testing.T) {
	rest := &Rest{
		jwtSecret: "my_secret",
		ErrorChan: make(chan error, 1),
	}

	testCases := []struct {
		name       string
		listenAddr string
		pprofAddr  string
	}{
		{"Case 1", "localhost:8080", "localhost:6060"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rest.Start(tc.listenAddr, tc.pprofAddr)
			assert.NotNil(t, rest.pprofServer)
			assert.NotNil(t, rest.server)

			time.Sleep(3 * time.Second)

			err := rest.Stop(context.Background())
			assert.NoError(t, err)
		})
	}
}
