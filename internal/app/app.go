package app

import (
	"context"
	"github.com/Azzonya/go-shortener/internal/api"
	"github.com/Azzonya/go-shortener/internal/cfg"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

type appSt struct {
	conf    *cfg.Conf
	api     *api.Rest
	storage *storage.Storage
}

func StopSignal() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	return ch
}

func (a *appSt) Init(conf *cfg.Conf) {
	var err error

	a.conf = conf

	if err = logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}

	a.storage, err = storage.NewStorage(conf.FileStoragePath)
	if err != nil {
		panic(err)
	}

	a.api = api.New(conf.BaseURL, a.storage)
}

func (a *appSt) Start() {
	a.api.Start(a.conf.HTTPListen)
}

func (a *appSt) Listen() {
	select {
	case <-StopSignal():
	case <-a.api.Wait():
	}
}

func (a *appSt) Stop() {
	a.storage.SyncData()

	if err := a.api.Stop(context.Background()); err != nil {
		panic(err)
	}
}

func Start() {
	conf := cfg.InitConfig()

	app := &appSt{}

	app.Init(&conf)
	app.Start()
	app.Listen()
	app.Stop()
}
