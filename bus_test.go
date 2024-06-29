package output_test

import (
	"github.com/majohn-r/output"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestNewDefaultBus(t *testing.T) {
	tests := map[string]struct {
		want              output.Bus
		wantConsoleWriter io.Writer
		wantErrorWriter   io.Writer
	}{
		"normal": {
			want:              output.NewCustomBus(os.Stdout, os.Stderr, output.NilLogger{}),
			wantConsoleWriter: os.Stdout,
			wantErrorWriter:   os.Stderr,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := output.NewDefaultBus(output.NilLogger{})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultBus() = %v, want %v", got, tt.want)
			}
			if w := got.ConsoleWriter(); w != tt.wantConsoleWriter {
				t.Errorf("NewDefaultBus() got console writer %v, want %v", w, tt.wantConsoleWriter)
			}
			if w := got.ErrorWriter(); w != tt.wantErrorWriter {
				t.Errorf("NewDefaultBus() got error writer %v, want %v", w, tt.wantErrorWriter)
			}
		})
	}
}
