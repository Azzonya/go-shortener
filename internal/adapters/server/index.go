package server

import (
	"context"
	"math/rand"
	"net/http"
	"strings"
)

type Rest struct {
	server *http.Server
	urlMap map[string]string
	//ch     chan error
}

func New() *Rest {
	return &Rest{
		urlMap: make(map[string]string),
	}
}

func (o *Rest) Start(lAddr string) {
	http.HandleFunc("/", o.HShortenerUrl)

	o.server = &http.Server{
		Addr: lAddr,
	}

	//go func() {
	err := o.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return
	}
	//}()
}

//func (o *Rest) Wait() chan error {
//	return o.ch
//}

func (o *Rest) Stop(ctx context.Context) error {
	err := o.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (o *Rest) generateShortURL() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const shortURLLength = 6
	var result strings.Builder
	for i := 0; i < shortURLLength; i++ {
		result.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}

	shortUrl := result.String()
	if o.urlMap[shortUrl] != "" {
		shortUrl = o.generateShortURL()
	}
	return shortUrl
}
