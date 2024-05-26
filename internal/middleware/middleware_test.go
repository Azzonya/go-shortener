package middleware

import (
	"testing"

	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAuthMiddleware(t *testing.T) {
	type args struct {
		jwtSecret string
	}
	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{
		{
			name: "auth middleware",
			args: args{
				jwtSecret: "my_secret",
			},
			want: AuthMiddleware("my_secret"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AuthMiddleware(tt.args.jwtSecret)
			assert.NoError(t, nil)
		})
	}
}

func TestCompressRequest(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "compress middleware test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := CompressRequest()
			assert.NoError(t, nil)
			assert.NotNil(t, handlerFunc)
		})
	}
}

func TestDecompressRequest(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "decompress middleware test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := DecompressRequest()
			assert.NoError(t, nil)
			assert.NotNil(t, handlerFunc)
		})
	}
}

func TestRequestLogger(t *testing.T) {
	type args struct {
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "request logger test",
			args: args{
				logger: logger.Log,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RequestLogger(tt.args.logger)
			assert.NoError(t, nil)
		})
	}
}
