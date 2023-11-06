package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type conf struct {
	HTTPListen string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func initConfig() conf {
	conf := conf{}
	flag.StringVar(&conf.HTTPListen, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "base address of the resulting shortened URL")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		panic(err)
	}

	return conf
}
