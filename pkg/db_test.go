package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDatabasePg(t *testing.T) {
	type args struct {
		pgDsn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "init db test",
			args: args{
				pgDsn: "postgresql://postgres:postgres@localhost:5437/postgresdb",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitDatabasePg(tt.args.pgDsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitDatabasePg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
