package cfg

import (
	"crypto/tls"
	"os"
	"reflect"
	"testing"
)

func TestUseDatabase(t *testing.T) {
	conf := Conf{}

	if conf.UseDatabase() {
		t.Error("Expected UseDatabase() to return false when PgDsn is not set")
	}

	conf.PgDsn = "testdbdsn"

	if !conf.UseDatabase() {
		t.Error("Expected UseDatabase() to return true when PgDsn is set")
	}
}

func TestInitConfig(t *testing.T) {
	os.Setenv("SERVER_ADDRESS", "localhost:9000")
	os.Setenv("PPROF_SERVER_ADDRESS", "localhost:9595")
	os.Setenv("BASE_URL", "http://example.com")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/short-url-repo-test.json")
	os.Setenv("DATABASE_DSN", "postgres://user:password@localhost:5432/dbname")
	os.Setenv("JWT_SECRET", "testsecret")

	// Подготавливаем аргументы командной строки
	os.Args = []string{"cmd", "-a", "localhost:8000", "-p", "localhost:9090", "-b", "http://localhost:8000", "-l", "info", "-f", "/tmp/short-url-repo.json", "-d", "testdbdsn", "-jwt_secret", "testjwtsecret"}

	tests := []struct {
		name string
		want Conf
	}{
		{
			name: "Init config",
			want: Conf{
				HTTPListen:      "localhost:9000",
				HTTPPprof:       "localhost:9595",
				BaseURL:         "http://example.com",
				LogLevel:        "debug",
				FileStoragePath: "/tmp/short-url-repo-test.json",
				PgDsn:           "postgres://user:password@localhost:5432/dbname",
				JWTSecret:       "testsecret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConf_LoadFromFile(t *testing.T) {
	type fields struct {
		TLSCertificate  *tls.Certificate
		HTTPListen      string
		HTTPPprof       string
		BaseURL         string
		LogLevel        string
		FileStoragePath string
		PgDsn           string
		JWTSecret       string
		ConfigFilePath  string
		EnableHTTPS     bool
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test load from file",
			fields: fields{},
			args: args{
				filePath: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				TLSCertificate:  tt.fields.TLSCertificate,
				HTTPListen:      tt.fields.HTTPListen,
				HTTPPprof:       tt.fields.HTTPPprof,
				BaseURL:         tt.fields.BaseURL,
				LogLevel:        tt.fields.LogLevel,
				FileStoragePath: tt.fields.FileStoragePath,
				PgDsn:           tt.fields.PgDsn,
				JWTSecret:       tt.fields.JWTSecret,
				ConfigFilePath:  tt.fields.ConfigFilePath,
				EnableHTTPS:     tt.fields.EnableHTTPS,
			}
			if err := c.LoadFromFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConf_applyConfig(t *testing.T) {
	type fields struct {
		TLSCertificate  *tls.Certificate
		HTTPListen      string
		HTTPPprof       string
		BaseURL         string
		LogLevel        string
		FileStoragePath string
		PgDsn           string
		JWTSecret       string
		ConfigFilePath  string
		EnableHTTPS     bool
	}
	type args struct {
		newConf Conf
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "test applyConfig",
			fields: fields{},
			args: args{
				newConf: Conf{
					JWTSecret: "jwt secret",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				TLSCertificate:  tt.fields.TLSCertificate,
				HTTPListen:      tt.fields.HTTPListen,
				HTTPPprof:       tt.fields.HTTPPprof,
				BaseURL:         tt.fields.BaseURL,
				LogLevel:        tt.fields.LogLevel,
				FileStoragePath: tt.fields.FileStoragePath,
				PgDsn:           tt.fields.PgDsn,
				JWTSecret:       tt.fields.JWTSecret,
				ConfigFilePath:  tt.fields.ConfigFilePath,
				EnableHTTPS:     tt.fields.EnableHTTPS,
			}
			c.applyConfig(tt.args.newConf)

			if c.JWTSecret != tt.args.newConf.JWTSecret {
				t.Errorf("not replaced")
			}
		})
	}
}

func TestConf_OverrideEnv(t *testing.T) {
	type fields struct {
		TLSCertificate  *tls.Certificate
		HTTPListen      string
		HTTPPprof       string
		BaseURL         string
		LogLevel        string
		FileStoragePath string
		PgDsn           string
		JWTSecret       string
		ConfigFilePath  string
		EnableHTTPS     bool
	}
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test override env",
			fields: fields{},
			args: args{
				name:  "JWT_SECRET",
				value: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				TLSCertificate:  tt.fields.TLSCertificate,
				HTTPListen:      tt.fields.HTTPListen,
				HTTPPprof:       tt.fields.HTTPPprof,
				BaseURL:         tt.fields.BaseURL,
				LogLevel:        tt.fields.LogLevel,
				FileStoragePath: tt.fields.FileStoragePath,
				PgDsn:           tt.fields.PgDsn,
				JWTSecret:       tt.fields.JWTSecret,
				ConfigFilePath:  tt.fields.ConfigFilePath,
				EnableHTTPS:     tt.fields.EnableHTTPS,
			}

			prevValue := os.Getenv(tt.args.name)

			if err := c.OverrideEnv(tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("OverrideEnv() error = %v, wantErr %v", err, tt.wantErr)
			}

			newValue := os.Getenv(tt.args.name)
			if newValue != tt.args.value {
				t.Errorf("not oveeride")
			}

			if err := c.OverrideEnv(tt.args.name, prevValue); (err != nil) != tt.wantErr {
				t.Errorf("OverrideEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
