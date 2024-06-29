package output

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func Test_bus_Log(t *testing.T) {
	type args struct {
		l    Level
		msg  string
		args map[string]any
	}
	tests := map[string]struct {
		args
		Log   string
		Error string
	}{
		"trace": {
			args: args{
				l:    Trace,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='trace' f='v' msg='hello'\n",
		},
		"debug": {
			args: args{
				l:    Debug,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='debug' f='v' msg='hello'\n",
		},
		"info": {
			args: args{
				l:    Info,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='info' f='v' msg='hello'\n",
		},
		"warning": {
			args: args{
				l:    Warning,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='warning' f='v' msg='hello'\n",
		},
		"error": {
			args: args{
				l:    Error,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='error' f='v' msg='hello'\n",
		},
		"panic": {
			args: args{
				l:    Panic,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='panic' f='v' msg='hello'\n",
		},
		"fatal": {
			args: args{
				l:    Fatal,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='fatal' f='v' msg='hello'\n",
		},
		"illegal": {
			args: args{
				l:    Trace + 1,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Error: "Programming error: call to bus.Log() with invalid level value 7; message: 'hello', args: 'map[f:v].\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eW := &bytes.Buffer{}
			l := NewRecordingLogger()
			o := NewCustomBus(nil, eW, l)
			o.Log(tt.args.l, tt.args.msg, tt.args.args)
			if got := l.String(); got != tt.Log {
				t.Errorf("bus.Log() got log %q want %q", got, tt.Log)
			}
			if got := eW.String(); got != tt.Error {
				t.Errorf("bus.Log() got error %q want %q", got, tt.Error)
			}
		})
	}
}

func Test_bus_WriteCanonicalError(t *testing.T) {
	type args struct {
		format string
		a      []any
	}
	tests := map[string]struct {
		args
		want string
	}{
		"broad test": {
			args: args{
				format: "test format %d %q %v..?!..?\n\n\n\n",
				a:      []any{25, "foo", 1.245},
			},
			want: "Test format 25 \"foo\" 1.245?\n",
		},
		"narrow test": {
			args: args{
				format: "1. test format %d %q %v",
				a:      []any{25, "foo", 1.245},
			},
			want: "1. test format 25 \"foo\" 1.245.\n",
		},
		"nothing but newlines": {
			args: args{
				format: "\n\n\n\n\n\n\n\n\n\n\n\n\n\n",
			},
			want: "\n",
		},
		"nothing but terminal punctuation": {
			args: args{
				format: "!!?.!?.",
			},
			want: ".\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			w := &bytes.Buffer{}
			o := &bus{errorWriter: w, performWrites: true}
			o.WriteCanonicalError(tt.args.format, tt.args.a...)
			if got := w.String(); got != tt.want {
				t.Errorf("bus.WriteCanonicalError() got %q want %q", got, tt.want)
			}
		})
	}
}

func Test_bus_WriteError(t *testing.T) {
	type args struct {
		format string
		a      []any
	}
	tests := map[string]struct {
		args
		want string
	}{
		"broad test": {
			args: args{
				format: "test format %d %q %v..?!..?\n\n\n\n",
				a:      []any{25, "foo", 1.245},
			},
			want: "test format 25 \"foo\" 1.245..?!..?\n\n\n\n",
		},
		"narrow test": {
			args: args{
				format: "1. test format %d %q %v",
				a:      []any{25, "foo", 1.245},
			},
			want: "1. test format 25 \"foo\" 1.245",
		},
		"nothing but newlines": {
			args: args{
				format: "\n\n\n\n\n\n\n\n\n\n\n\n\n\n",
			},
			want: "\n\n\n\n\n\n\n\n\n\n\n\n\n\n",
		},
		"nothing but terminal punctuation": {
			args: args{
				format: "!!?.!?.",
			},
			want: "!!?.!?.",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			w := &bytes.Buffer{}
			o := &bus{errorWriter: w, performWrites: true}
			o.WriteError(tt.args.format, tt.args.a...)
			if got := w.String(); got != tt.want {
				t.Errorf("bus.WriteError() got %q want %q", got, tt.want)
			}
		})
	}
}

func Test_bus_WriteCanonicalConsole(t *testing.T) {
	type args struct {
		format string
		a      []any
	}
	tests := map[string]struct {
		w   *bytes.Buffer
		tab uint8
		args
		want string
	}{
		"strict rules": {
			w:   &bytes.Buffer{},
			tab: 0,
			args: args{
				format: "test %s...\n\n",
				a:      []any{"foo."},
			},
			want: "Test foo.\n",
		},
		"strict rules with tab": {
			w:   &bytes.Buffer{},
			tab: 10,
			args: args{
				format: "test %s...\n\n",
				a:      []any{"foo."},
			},
			want: "          Test foo.\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := &bus{consoleWriter: tt.w, performWrites: true, tab: tt.tab}
			o.WriteCanonicalConsole(tt.args.format, tt.args.a...)
			if got := tt.w.String(); got != tt.want {
				t.Errorf("bus.WriteCanonicalConsole() got %q want %q", got, tt.want)
			}
		})
	}
}

func Test_bus_WriteConsole(t *testing.T) {
	type args struct {
		format string
		a      []any
	}
	tests := map[string]struct {
		w   *bytes.Buffer
		tab uint8
		args
		want string
	}{
		"lax rules": {
			w:   &bytes.Buffer{},
			tab: 0,
			args: args{
				format: "test %s...\n\n",
				a:      []any{"foo."},
			},
			want: "test foo....\n\n",
		},
		"lax rules with non-zero tab": {
			w:   &bytes.Buffer{},
			tab: 5,
			args: args{
				format: "test %s...\n\n",
				a:      []any{"foo."},
			},
			want: "     test foo....\n\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := &bus{consoleWriter: tt.w, performWrites: true, tab: tt.tab}
			o.WriteConsole(tt.args.format, tt.args.a...)
			if got := tt.w.String(); got != tt.want {
				t.Errorf("bus.WriteConsole() got %q want %q", got, tt.want)
			}
		})
	}
}

func Test_bus_IsConsoleTTY(t *testing.T) {
	tests := map[string]struct {
		b    Bus
		want bool
	}{
		"simple": {b: NewDefaultBus(NewRecordingLogger()), want: false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.b.IsConsoleTTY(); got != tt.want {
				t.Errorf("bus.IsConsoleTTY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bus_IsErrorTTY(t *testing.T) {
	tests := map[string]struct {
		b    Bus
		want bool
	}{
		"simple": {b: NewDefaultBus(NewRecordingLogger()), want: false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.b.IsErrorTTY(); got != tt.want {
				t.Errorf("bus.IsErrorTTY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTTY(t *testing.T) {
	oldIsTerminal := isTerminal
	oldIsCygwinTerminal := isCygwinTerminal
	defer func() {
		isTerminal = oldIsTerminal
		isCygwinTerminal = oldIsCygwinTerminal
	}()
	tests := map[string]struct {
		terminalFunc       func(uintptr) bool
		cygwinTerminalFunc func(uintptr) bool
		w                  io.Writer
		want               bool
	}{
		"non-file": {w: &bytes.Buffer{}},
		"non-tty": {
			terminalFunc: func(_ uintptr) bool {
				return false
			},
			cygwinTerminalFunc: func(_ uintptr) bool {
				return false
			},
			w: os.Stdout,
		},
		"plain terminal": {
			terminalFunc: func(_ uintptr) bool {
				return true
			},
			cygwinTerminalFunc: func(_ uintptr) bool {
				return false
			},
			w:    os.Stdout,
			want: true,
		},
		"cygwin terminal": {
			terminalFunc: func(_ uintptr) bool {
				return false
			},
			cygwinTerminalFunc: func(_ uintptr) bool {
				return true
			},
			w:    os.Stdout,
			want: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			isTerminal = tt.terminalFunc
			isCygwinTerminal = tt.cygwinTerminalFunc
			if got := isTTY(tt.w); got != tt.want {
				t.Errorf("isTTY() = %t, want %t", got, tt.want)
			}
		})
	}
}

func Test_bus_IncrementTab(t *testing.T) {
	tests := map[string]struct {
		initialTab uint8
		t          uint8
		want       uint8
	}{
		"typical":       {initialTab: 0, t: 2, want: 2},
		"overflow":      {initialTab: 64, t: 192, want: 64},
		"near-overflow": {initialTab: 64, t: 191, want: 255},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &bus{tab: tt.initialTab}
			b.IncrementTab(tt.t)
			if got := b.Tab(); got != tt.want {
				t.Errorf("bus.IncrementTab got %d want %d", got, tt.want)
			}
		})
	}
}

func Test_bus_DecrementTab(t *testing.T) {
	tests := map[string]struct {
		initialTab uint8
		t          uint8
		want       uint8
	}{
		"typical":   {initialTab: 2, t: 2, want: 0},
		"underflow": {initialTab: 2, t: 3, want: 2},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &bus{tab: tt.initialTab}
			b.DecrementTab(tt.t)
			if got := b.Tab(); got != tt.want {
				t.Errorf("bus.DecrementTab got %d want %d", got, tt.want)
			}
		})
	}
}
