package main

import "flag"

type conf struct {
	HTTPListen string `mapstructure:"HTTP_LISTEN"`
	BaseURL    string `mapstructure:"Base_URL"`
}

func initConfig() conf {
	conf := conf{}
	flag.StringVar(&conf.HTTPListen, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "base address of the resulting shortened URL")
	flag.Parse()

	return conf
}
