// Package cfg provides functionality to initialize and parse configuration for the URL shortener application.
package cfg

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env/v6"
)

// Conf represents the application configuration.
type Conf struct {
	TLSCertificate  *tls.Certificate `env:"TLS_CERTIFICATE"`      // TLSCertificate represents the TLS certificate.
	HTTPListen      string           `env:"SERVER_ADDRESS"`       // HTTPListen represents the address and port to run the HTTP server.
	HTTPPprof       string           `env:"PPROF_SERVER_ADDRESS"` // HTTPPprof represents the address and port to run the pprof server.
	BaseURL         string           `env:"BASE_URL"`             // BaseURL represents the base address of the resulting shortened URL.
	LogLevel        string           `env:"LOG_LEVEL"`            // LogLevel represents the log level for logging.
	FileStoragePath string           `env:"FILE_STORAGE_PATH"`    // FileStoragePath represents the file path for storing data.
	PgDsn           string           `env:"DATABASE_DSN"`         // PgDsn represents the database connection line for PostgreSQL.
	JWTSecret       string           `env:"JWT_SECRET"`           // JWTSecret represents the JWT cookie secret for authentication.
	ConfigFilePath  string           `env:"config_file_path"`     // ConfigFilePath represents the path to config file
	EnableHTTPS     bool             `env:"ENABLE_HTTPS"`         // EnableHTTPS specifies whether to enable HTTPS for the server.
	TrustedSubnet   string           `env:"TRUSTED_SUBNET"`       // TrustedSubnet representation of classless addressing (CIDR)
}

// InitConfig initializes the application configuration from environment variables and command-line flags.
func InitConfig() Conf {
	conf := Conf{}
	flag.StringVar(&conf.ConfigFilePath, "c", "", "path to config file")
	flag.StringVar(&conf.HTTPListen, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.HTTPPprof, "p", "localhost:9595", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "base address of the resulting shortened URL")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.StringVar(&conf.FileStoragePath, "f", "/tmp/short-url-repo.json", "file path")
	flag.StringVar(&conf.PgDsn, "d", "", "database connection line")
	flag.StringVar(&conf.JWTSecret, "jwt_secret", "supersecret", "jwt cookie secret")
	flag.StringVar(&conf.TrustedSubnet, "t", "", "trusted subnet")
	flag.BoolVar(&conf.EnableHTTPS, "s", false, "http or https")
	flag.Parse()

	var err error

	if conf.ConfigFilePath != "" {
		if err = conf.OverrideEnv("CONFIG", conf.ConfigFilePath); err != nil {
			panic(err)
		}

		err = conf.LoadFromFile(conf.FileStoragePath)
		if err != nil {
			panic(err)
		}
	}

	err = env.Parse(&conf)
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

// LoadFromFile loads config from file
func (c *Conf) LoadFromFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", filePath)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open config file: %w", err)
	}
	defer f.Close()

	cfgBytes, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}

	var cfg Conf
	if err := json.Unmarshal(cfgBytes, &cfg); err != nil {
		return fmt.Errorf("cannot unmarshal config: %w", err)
	}

	c.applyConfig(cfg)

	return nil
}

// applyConfig replace new values from file
func (c *Conf) applyConfig(newConf Conf) {
	if newConf.HTTPListen != "" {
		c.HTTPListen = newConf.HTTPListen
	}
	if newConf.HTTPPprof != "" {
		c.HTTPPprof = newConf.HTTPPprof
	}
	if newConf.BaseURL != "" {
		c.BaseURL = newConf.BaseURL
	}
	if newConf.LogLevel != "" {
		c.LogLevel = newConf.LogLevel
	}
	if newConf.FileStoragePath != "" {
		c.FileStoragePath = newConf.FileStoragePath
	}
	if newConf.PgDsn != "" {
		c.PgDsn = newConf.PgDsn
	}
	if newConf.JWTSecret != "" {
		c.JWTSecret = newConf.JWTSecret
	}
	if newConf.EnableHTTPS {
		c.EnableHTTPS = newConf.EnableHTTPS
	}
	if newConf.TLSCertificate != nil {
		c.TLSCertificate = newConf.TLSCertificate
	}

	if newConf.TrustedSubnet != "" {
		c.TrustedSubnet = newConf.TrustedSubnet
	}
}

// OverrideEnv override env value
func (c *Conf) OverrideEnv(name string, value string) error {
	if value == "" {
		return nil
	}

	err := os.Setenv(name, value)
	if err != nil {
		return err
	}

	return nil
}

// UseDatabase checks if the application is configured to use a database.
func (c *Conf) UseDatabase() bool {
	return len(c.PgDsn) > 0
}
