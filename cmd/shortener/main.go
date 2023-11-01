package main

import (
	"github.com/Azzonya/go-shortener/internal/adapters/server"
)

func main() {
	initConfig()

	app := struct {
		srv *server.Rest
	}{}

	app.srv = server.New()

	app.srv.Start(":8080")
}
