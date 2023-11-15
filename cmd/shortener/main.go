package main

import (
	"github.com/Azzonya/go-shortener/internal/cfg"
	"github.com/Azzonya/go-shortener/internal/service"
)

func main() {
	conf := cfg.InitConfig()

	srv := service.New(conf.BaseURL)

	srv.Start(conf.HTTPListen)
}
