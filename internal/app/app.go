package app

import (
	"github.com/Azzonya/go-shortener/internal/api"
	"github.com/Azzonya/go-shortener/internal/cfg"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/storage"
)

type appSt struct {
	conf    *cfg.Conf
	api     *api.Rest
	storage *storage.Storage
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

func Start() {
	conf := cfg.InitConfig()

	app := &appSt{}

	app.Init(&conf)
	app.Start()
}
