// Package cfg provides functionality to initialize and parse configuration for the URL shortener application.
package cfg

import (
	"crypto/tls"
	"flag"

	"github.com/caarlos0/env/v6"
)

// Conf represents the application configuration.
type Conf struct {
	HTTPListen      string           `env:"SERVER_ADDRESS"`       // HTTPListen represents the address and port to run the HTTP server.
	HTTPPprof       string           `env:"PPROF_SERVER_ADDRESS"` // HTTPPprof represents the address and port to run the pprof server.
	BaseURL         string           `env:"BASE_URL"`             // BaseURL represents the base address of the resulting shortened URL.
	LogLevel        string           `env:"LOG_LEVEL"`            // LogLevel represents the log level for logging.
	FileStoragePath string           `env:"FILE_STORAGE_PATH"`    // FileStoragePath represents the file path for storing data.
	PgDsn           string           `env:"DATABASE_DSN"`         // PgDsn represents the database connection line for PostgreSQL.
	JWTSecret       string           `env:"JWT_SECRET"`           // JWTSecret represents the JWT cookie secret for authentication.
	EnableHTTPS     bool             `env:"ENABLE_HTTPS"`         // EnableHTTPS specifies whether to enable HTTPS for the server.
	TLSCertificate  *tls.Certificate `env:"TLS_CERTIFICATE"`      // TLSCertificate represents the TLS certificate.
}

// InitConfig initializes the application configuration from environment variables and command-line flags.
func InitConfig() Conf {
	conf := Conf{}
	flag.StringVar(&conf.HTTPListen, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.HTTPPprof, "p", "localhost:9595", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "base address of the resulting shortened URL")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.StringVar(&conf.FileStoragePath, "f", "/tmp/short-url-repo.json", "file path")
	flag.StringVar(&conf.PgDsn, "d", "", "database connection line")
	flag.StringVar(&conf.JWTSecret, "jwt_secret", "supersecret", "jwt cookie secret")
	flag.BoolVar(&conf.EnableHTTPS, "s", false, "http or https")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		panic(err)
	}

	if conf.EnableHTTPS {
		tlsCertificate, err := generateSelfSignedCert()
		if err != nil {
			panic(err)
		}
		conf.TLSCertificate = &tlsCertificate
	}

	return conf
}

// UseDatabase checks if the application is configured to use a database.
func (c *Conf) UseDatabase() bool {
	return len(c.PgDsn) > 0
}
