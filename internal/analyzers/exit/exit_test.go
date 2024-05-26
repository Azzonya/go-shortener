package exit

import (
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *analysis.Analyzer
	}{
		{
			name: "new error check",
			want: &analysis.Analyzer{
				Name: "osexitcheck",
				Doc:  "checks for direct os.Exit calls inside main-like functions",
				Run:  run,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got.Name, tt.want.Name) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
