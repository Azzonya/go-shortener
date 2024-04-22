package cfg

import (
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
