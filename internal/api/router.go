package api

import (
	"context"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/middleware"
	storage2 "github.com/Azzonya/go-shortener/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Rest struct {
	server  *http.Server
	storage *storage2.Storage
	baseURL string
	logger  zap.Logger
	//ch     chan error
}

func New(baseURL string) *Rest {
	return &Rest{
		baseURL: baseURL,
		storage: storage2.NewStorage(),
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

	err := o.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return
	}
}

func (o *Rest) Stop(ctx context.Context) error {
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
