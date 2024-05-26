package cfg

import (
	"testing"
)

func Test_generateSelfSignedCert(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test TLS certificate generation",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := generateSelfSignedCert()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSelfSignedCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
