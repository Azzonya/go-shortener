package main

import (
	"github.com/Azzonya/go-shortener/internal/adapters/server"
)

func main() {
	conf := initConfig()

	app := struct {
		srv *server.Rest
	}{}

	app.srv = server.New(conf.BaseURL) //

	app.srv.Start(conf.HTTPListen)
}
