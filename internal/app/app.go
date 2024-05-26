// Package app provides functionality to initialize, start, and stop the URL shortener application.
// It sets up the API server, repository, and database connections based on the configuration.
package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Azzonya/go-shortener/internal/api"
	"github.com/Azzonya/go-shortener/internal/cfg"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/repo"
	"github.com/Azzonya/go-shortener/internal/repo/inmemory"
	"github.com/Azzonya/go-shortener/internal/repo/pg"
	"github.com/Azzonya/go-shortener/internal/shortener"
	"github.com/Azzonya/go-shortener/pkg"
)

// appSt represents the application state containing configuration, API server, shortener, database connection, and repository.
type appSt struct {
	conf      *cfg.Conf
	api       *api.Rest
	shortener *shortener.Shortener
	db        *pgxpool.Pool
	repo      repo.Repo
}

// StopSignal returns a channel for receiving OS signals to stop the application.
func StopSignal() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	return ch
}

// Init initializes the application with the provided configuration.
func (a *appSt) Init(conf *cfg.Conf) {
	var err error

	a.conf = conf

	if conf.UseDatabase() {
		a.db, err = pkg.InitDatabasePg(conf.PgDsn)
		if err != nil {
			panic(err)
		}

		a.repo = pg.New(a.db)
	} else {
		a.repo, err = inmemory.New(conf.FileStoragePath)
		if err != nil {
			panic(err)
		}
	}

	if err = logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}

	a.shortener = shortener.New(conf.BaseURL, a.repo)

	a.api = api.New(
		a.shortener,
		conf.JWTSecret,
		a.conf.EnableHTTPS,
		a.conf.TLSCertificate,
	)
}

// Start starts the application, initializing and running the API server.
func (a *appSt) Start() {
	a.api.Start(a.conf.HTTPListen, a.conf.HTTPPprof)
}

// Listen listens for signals to stop the application.
func (a *appSt) Listen() {
	select {
	case <-StopSignal():
	case <-a.api.ErrorChan:
	}
}

// Stop stops the application, closing database connections and shutting down the API server.
func (a *appSt) Stop() {
	if !a.conf.UseDatabase() {
		a.repo.SyncData()
		a.db.Close()
	}

	if err := a.api.Stop(context.Background()); err != nil {
		panic(err)
	}
}

// Start initializes and starts the URL shortener application.
func Start() {
	conf := cfg.InitConfig()

	app := &appSt{}

	app.Init(&conf)
	app.Start()
	app.Listen()
	app.Stop()
}
