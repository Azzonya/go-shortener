// Package app provides functionality to initialize, start, and stop the URL shortener application.
// It sets up the API server, repository, and database connections based on the configuration.
package app

import (
	"context"
	"github.com/Azzonya/go-shortener/internal/interceptor"
	pb "github.com/Azzonya/go-shortener/pkg/proto/shortener"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
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
	grpcSrv   *grpc.Server
	shortener *shortener.Shortener
	db        *pgxpool.Pool
	repo      repo.Repo
}

// StopSignal returns a channel for receiving OS signals to stop the application.
func StopSignal() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
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
		a.conf.TrustedSubnet,
		a.conf.EnableHTTPS,
		a.conf.TLSCertificate,
	)

	interceptors := make([]grpc.UnaryServerInterceptor, 0, 3)

	interceptors = append(interceptors, interceptor.GrpcInterceptorLogger())
	interceptors = append(interceptors, interceptor.GrpcInterceptorAuth())

	a.grpcSrv = grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors...,
	))
	grpcHandlers := api.NewGrpcHandlers(a.shortener)
	pb.RegisterShortenerServer(a.grpcSrv, grpcHandlers)

	reflection.Register(a.grpcSrv)
}

// Start starts the application, initializing and running the API server.
func (a *appSt) Start() {
	a.api.Start(a.conf.HTTPListen, a.conf.HTTPPprof)
	logger.Log.Info("Rest api started " + a.conf.HTTPListen)

	lis, err := net.Listen("tcp", ":"+a.conf.GrpcPort)
	if err != nil {
		panic(err)
	}
	go func() {
		err = a.grpcSrv.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	logger.Log.Info("grpc-server started " + lis.Addr().String())
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

	a.grpcSrv.GracefulStop()

	if err := a.api.Stop(context.Background()); err != nil {
		panic(err)
	}

	<-a.api.IdleConnsClosed
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
