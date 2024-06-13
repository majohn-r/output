package output

import (
	"testing"
)

func TestNewNilBus(t *testing.T) {
	tests := map[string]struct{}{
		"simple test": {},
	}
	for name := range tests {
		t.Run(name, func(t *testing.T) {
			o := NewNilBus()
			// none of these calls do anything ... here for pedantic test coverage
			o.WriteCanonicalConsole("%s %d %t", "foo", 42, true)
			o.WriteConsole("%s %d %t", "foo", 42, true)
			o.WriteCanonicalError("%s %d %t", "foo", 42, true)
			o.WriteError("%s %d %t", "foo", 42, true)
		})
	}
}

func TestNilWriter_Write(t *testing.T) {
	tests := map[string]struct {
		nw      NilWriter
		p       []byte
		wantN   int
		wantErr bool
	}{
		"a few bytes": {
			nw:      NilWriter{},
			p:       []byte{0, 1, 2},
			wantN:   3,
			wantErr: false,
		},
		"nil": {
			nw:      NilWriter{},
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
		nl   NilLogger
		args args
	}{"default": {nl: NilLogger{}, args: args{}}}
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
