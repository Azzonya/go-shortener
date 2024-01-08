package app

import (
	"context"
	"github.com/Azzonya/go-shortener/internal/api"
	"github.com/Azzonya/go-shortener/internal/cfg"
	"github.com/Azzonya/go-shortener/internal/inmemory"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/repo/pg"
	shortener_service "github.com/Azzonya/go-shortener/internal/shortener"
	"github.com/Azzonya/go-shortener/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"os/signal"
	"syscall"
)

type appSt struct {
	conf      *cfg.Conf
	api       *api.Rest
	storage   *inmemory.Storage
	shortener *shortener_service.Shortener
	db        *pgxpool.Pool
	repo      *pg.St
}

func StopSignal() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	return ch
}

func (a *appSt) Init(conf *cfg.Conf) {
	var err error

	a.conf = conf

	useDB := false
	if conf.UseDatabase() {
		a.db, err = pkg.InitDatabasePg(conf.PgDsn)
		if err != nil {
			panic(err)
		}

		a.repo = pg.New(a.db)

		useDB = true
	}

	if err = logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}

	a.storage, err = inmemory.NewStorage(conf.FileStoragePath, useDB)
	if err != nil {
		panic(err)
	}

	a.shortener = shortener_service.New(conf.BaseURL, a.storage, a.repo, useDB)

	a.api = api.New(a.shortener, conf.JWTSecret)
}

func (a *appSt) Start() {
	a.api.Start(a.conf.HTTPListen)
}

func (a *appSt) Listen() {
	select {
	case <-StopSignal():
	case <-a.api.ErrorChan:
	}
}

func (a *appSt) Stop() {
	if !a.shortener.UseDB {
		a.storage.SyncData()
		a.db.Close()
	}

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
