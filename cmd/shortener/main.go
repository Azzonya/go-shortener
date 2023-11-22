package main

import (
	"github.com/Azzonya/go-shortener/internal/cfg"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/service"
)

func main() {
	conf := cfg.InitConfig()

	if err := logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}
	srv := service.New(conf.BaseURL)

	srv.Start(conf.HTTPListen)
}
