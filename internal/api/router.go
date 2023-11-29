package api

import (
	"context"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/middleware"
	shortener_service "github.com/Azzonya/go-shortener/internal/shortener-service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Rest struct {
	server    *http.Server
	shortener *shortener_service.Shortener

	eChan chan error
}

func New(shortener *shortener_service.Shortener) *Rest {
	return &Rest{
		shortener: shortener,

		eChan: make(chan error, 1),
	}
}

func (o *Rest) Start(lAddr string) {
	logger.Log.Info("Running server", zap.String("address", lAddr))

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(middleware.RequestLogger(logger.Log), gin.Recovery())
	r.Use(middleware.CompressRequest(), gin.Recovery())
	r.Use(middleware.DecompressRequest(), gin.Recovery())

	o.SetRouters(r)

	o.server = &http.Server{
		Addr:    lAddr,
		Handler: r,
	}

	go func() {
		err := o.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			o.eChan <- err
		}
	}()
}

func (o *Rest) Wait() <-chan error {
	return o.eChan
}

func (o *Rest) Stop(ctx context.Context) error {
	defer close(o.eChan)

	err := o.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (o *Rest) SetRouters(r *gin.Engine) {
	r.POST("/", o.Shorten)
	r.GET("/:id", o.Redirect)
	r.POST("/api/shorten", o.ShortenJSON)
}
