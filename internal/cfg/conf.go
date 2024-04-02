package cfg

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Conf struct {
	HTTPListen      string `env:"SERVER_ADDRESS"`
	HTTPPprof       string `env:"PPROF_SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	PgDsn           string `env:"DATABASE_DSN"`
	JWTSecret       string `env:"JWT_SECRET"`
}

func InitConfig() Conf {
	conf := Conf{}
	flag.StringVar(&conf.HTTPListen, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.HTTPPprof, "p", "localhost:9595", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "base address of the resulting shortened URL")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.StringVar(&conf.FileStoragePath, "f", "/tmp/short-url-repo.json", "file path")
	flag.StringVar(&conf.PgDsn, "d", "", "database connection line")
	flag.StringVar(&conf.JWTSecret, "jwt_secret", "supersecret", "jwt cookie secret")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		panic(err)
	}

	return conf
}

func (c *Conf) UseDatabase() bool {
	return len(c.PgDsn) > 0
}
