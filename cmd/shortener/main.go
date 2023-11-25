package main

import (
	"github.com/Azzonya/go-shortener/internal/api"
	"github.com/Azzonya/go-shortener/internal/cfg"
	"github.com/Azzonya/go-shortener/internal/logger"
)

func main() {
	conf := cfg.InitConfig()

	if err := logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}
	srv := api.New(conf.BaseURL)

	srv.Start(conf.HTTPListen)
}
