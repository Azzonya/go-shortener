// Package api provides the HTTP handlers and server setup for the URL shortening REST API.
package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/middleware"
	shortener_service "github.com/Azzonya/go-shortener/internal/shortener"
)

// Rest represents the REST API server.
type Rest struct {
	server      *http.Server
	pprofServer *http.Server
	shortener   *shortener_service.Shortener

	ErrorChan chan error
	jwtSecret string
}

// New creates a new instance of the REST API server.
func New(shortener *shortener_service.Shortener, jwtSecret string) *Rest {
	return &Rest{
		shortener: shortener,

		ErrorChan: make(chan error, 1),
		jwtSecret: jwtSecret,
	}
}

// Start starts the REST API server.
func (o *Rest) Start(lAddr, pAddr string) {
	logger.Log.Info("Running server", zap.String("address", lAddr))

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(
		middleware.RequestLogger(logger.Log),
		middleware.CompressRequest(),
		middleware.DecompressRequest(),
		middleware.AuthMiddleware(o.jwtSecret),
		gin.Recovery())

	o.SetRouters(r)

	o.server = &http.Server{
		Addr:    lAddr,
		Handler: r,
	}

	o.pprofServer = &http.Server{
		Addr: pAddr,
	}

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Log.Error("Recovered from panic: %v")
			}
		}()

		err := o.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			o.ErrorChan <- err
		}
	}()

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Log.Error("Recovered from panic: %v")
			}
		}()

		err := o.pprofServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			o.ErrorChan <- err
		}
	}()
}

// Stop stops the REST API server.
func (o *Rest) Stop(ctx context.Context) error {
	defer close(o.ErrorChan)

	err := o.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

// SetRouters sets up the routes for the REST API server.
func (o *Rest) SetRouters(r *gin.Engine) {
	r.POST("/", o.Shorten)
	r.GET("/:id", o.Redirect)
	r.POST("/api/shorten", o.ShortenJSON)
	r.GET("/ping", o.Ping)
	r.POST("/api/shorten/batch", o.ShortenURLs)
	r.GET("/api/user/urls", o.ListAll)
	r.DELETE("/api/user/urls", o.DeleteURLs)
}
