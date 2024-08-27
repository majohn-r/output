package output_test

import (
	"github.com/majohn-r/output"
	"testing"
)

func TestNewNilBus(t *testing.T) {
	t.Run("basic check", func(t *testing.T) {
		o := output.NewNilBus()
		if o == nil {
			t.Error("NewNilBus() should not return nil")
		}
	})
}

func TestNilWriter_Write(t *testing.T) {
	tests := map[string]struct {
		nw      output.NilWriter
		p       []byte
		wantN   int
		wantErr bool
	}{
		"a few bytes": {
			nw:      output.NilWriter{},
			p:       []byte{0, 1, 2},
			wantN:   3,
			wantErr: false,
		},
		"nil": {
			nw:      output.NilWriter{},
			p:       nil,
			wantN:   0,
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotN, err := tt.nw.Write(tt.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("NilWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("NilWriter.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestNilLogger_All(t *testing.T) {
	type args struct {
		msg    string
		fields map[string]any
	}
	tests := map[string]struct {
		nl   output.NilLogger
		args args
	}{"default": {nl: output.NilLogger{}, args: args{}}}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.nl.Debug(tt.args.msg, tt.args.fields)
			tt.nl.Error(tt.args.msg, tt.args.fields)
			tt.nl.Fatal(tt.args.msg, tt.args.fields)
			tt.nl.Info(tt.args.msg, tt.args.fields)
			tt.nl.Panic(tt.args.msg, tt.args.fields)
			tt.nl.Trace(tt.args.msg, tt.args.fields)
			tt.nl.Warning(tt.args.msg, tt.args.fields)
		})
	}
}
