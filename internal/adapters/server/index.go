package server

import (
	"context"
	"math/rand"
	"net/http"
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
	http.HandleFunc("/", o.HShortenerURL)

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
	const shorURLLenth = 8
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, shorURLLenth)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
