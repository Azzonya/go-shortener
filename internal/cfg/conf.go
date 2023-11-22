package cfg

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Conf struct {
	HTTPListen string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
	LogLevel   string `env:"LOG_LEVEL"`
}

func InitConfig() Conf {
	conf := Conf{}
	flag.StringVar(&conf.HTTPListen, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "base address of the resulting shortened URL")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		panic(err)
	}

	return conf
}
