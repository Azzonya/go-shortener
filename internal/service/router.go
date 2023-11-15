package service

import (
	"context"
	storage2 "github.com/Azzonya/go-shortener/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Rest struct {
	server  *http.Server
	storage *storage2.Storage
	baseURL string
	//ch     chan error
}

func New(baseURL string) *Rest {
	return &Rest{
		baseURL: baseURL,
		storage: storage2.NewStorage(),
	}
}

func (o *Rest) Start(lAddr string) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/", o.Shorten)
	r.GET("/:id", o.Redirect)

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
