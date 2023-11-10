package main

import (
	"github.com/Azzonya/go-shortener/internal/service"
)

func main() {
	conf := initConfig()

	srv := service.New(conf.BaseURL)

	srv.Start(conf.HTTPListen)
}
