package output_test

import (
	"fmt"
	"github.com/majohn-r/output"
	"reflect"
	"strings"
	"testing"
)

func TestNewRecorderPrinting(t *testing.T) {
	tests := map[string]struct {
		usePrintf   bool
		consoleFmt  string
		consoleArgs []any
		consoleMsg  string
		errorFmt    string
		errorArgs   []any
		errorMsg    string
		logMessage  string
		logArgs     map[string]any
		tab         uint8
		enableList  bool
		numericList bool
		output.WantedRecording
	}{
		"println test": {
			consoleMsg:  "hello 42 true",
			errorMsg:    "24 false bye",
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         0,
			enableList:  false,
			numericList: false,
			WantedRecording: output.WantedRecording{
				Console: "hello 42 true\n",
				Error:   "24 false bye\n",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"println test with bulleted list": {
			consoleMsg:  "hello 42 true",
			errorMsg:    "24 false bye",
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         0,
			enableList:  true,
			numericList: false,
			WantedRecording: output.WantedRecording{
				Console: "● hello 42 true\n",
				Error:   "● 24 false bye\n",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"println test with numeric list": {
			consoleMsg:  "hello 42 true",
			errorMsg:    "24 false bye",
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         0,
			enableList:  true,
			numericList: true,
			WantedRecording: output.WantedRecording{
				Console: " 1. hello 42 true\n",
				Error:   " 1. 24 false bye\n",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"printf test": {
			usePrintf:   true,
			consoleFmt:  "%s %d %t",
			consoleArgs: []any{"hello", 42, true},
			errorFmt:    "%d %t %s",
			errorArgs:   []any{24, false, "bye"},
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         0,
			enableList:  false,
			numericList: false,
			WantedRecording: output.WantedRecording{
				Console: "hello 42 true",
				Error:   "24 false bye",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"printf test with bulleted list": {
			usePrintf:   true,
			consoleFmt:  "%s %d %t",
			consoleArgs: []any{"hello", 42, true},
			errorFmt:    "%d %t %s",
			errorArgs:   []any{24, false, "bye"},
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         0,
			enableList:  true,
			numericList: false,
			WantedRecording: output.WantedRecording{
				Console: "● hello 42 true",
				Error:   "● 24 false bye",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"printf test with numeric list": {
			usePrintf:   true,
			consoleFmt:  "%s %d %t",
			consoleArgs: []any{"hello", 42, true},
			errorFmt:    "%d %t %s",
			errorArgs:   []any{24, false, "bye"},
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         0,
			enableList:  true,
			numericList: true,
			WantedRecording: output.WantedRecording{
				Console: " 1. hello 42 true",
				Error:   " 1. 24 false bye",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"println with tab": {
			consoleMsg:  "hello 42 true",
			errorMsg:    "24 false bye",
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         6,
			enableList:  false,
			numericList: false,
			WantedRecording: output.WantedRecording{
				Console: "      hello 42 true\n",
				Error:   "24 false bye\n",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"println test with tab and bulleted list": {
			consoleMsg:  "hello 42 true",
			errorMsg:    "24 false bye",
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         6,
			enableList:  true,
			numericList: false,
			WantedRecording: output.WantedRecording{
				Console: "      ● hello 42 true\n",
				Error:   "● 24 false bye\n",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"println test with tab and numeric list": {
			consoleMsg:  "hello 42 true",
			errorMsg:    "24 false bye",
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         6,
			enableList:  true,
			numericList: true,
			WantedRecording: output.WantedRecording{
				Console: "       1. hello 42 true\n",
				Error:   " 1. 24 false bye\n",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"printf test with tab": {
			usePrintf:   true,
			consoleFmt:  "%s %d %t",
			consoleArgs: []any{"hello", 42, true},
			errorFmt:    "%d %t %s",
			errorArgs:   []any{24, false, "bye"},
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         4,
			enableList:  false,
			numericList: false,
			WantedRecording: output.WantedRecording{
				Console: "    hello 42 true",
				Error:   "24 false bye",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"printf test with tab and bulleted list": {
			usePrintf:   true,
			consoleFmt:  "%s %d %t",
			consoleArgs: []any{"hello", 42, true},
			errorFmt:    "%d %t %s",
			errorArgs:   []any{24, false, "bye"},
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         4,
			enableList:  true,
			numericList: false,
			WantedRecording: output.WantedRecording{
				Console: "    ● hello 42 true",
				Error:   "● 24 false bye",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
		"printf test with tab and numeric list": {
			usePrintf:   true,
			consoleFmt:  "%s %d %t",
			consoleArgs: []any{"hello", 42, true},
			errorFmt:    "%d %t %s",
			errorArgs:   []any{24, false, "bye"},
			logMessage:  "hello!",
			logArgs:     map[string]any{"field": "value"},
			tab:         4,
			enableList:  true,
			numericList: true,
			WantedRecording: output.WantedRecording{
				Console: "     1. hello 42 true",
				Error:   " 1. 24 false bye",
				Log:     "level='error' field='value' msg='hello!'\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			o.IncrementTab(tt.tab)
			var i any = o
			if _, ok := i.(output.Bus); !ok {
				t.Errorf("NewRecorder() Recorder does not implement Bus")
			}
			if o.ConsoleWriter() == nil {
				t.Errorf("NewRecorder() console writer is nil")
			}
			if o.ErrorWriter() == nil {
				t.Errorf("NewRecorder() error writer is nil")
			}
			if tt.enableList {
				o.BeginConsoleList(tt.numericList)
				o.BeginErrorList(tt.numericList)
			}
			if tt.usePrintf {
				o.ConsolePrintf(tt.consoleFmt, tt.consoleArgs...)
				o.ErrorPrintf(tt.errorFmt, tt.errorArgs...)
			} else {
				o.ConsolePrintln(tt.consoleMsg)
				o.ErrorPrintln(tt.errorMsg)
			}
			o.EndConsoleList()
			o.EndErrorList()
			o.Log(output.Error, tt.logMessage, tt.logArgs)
			if got := o.ConsoleOutput(); got != tt.WantedRecording.Console {
				t.Errorf("NewRecorder().ConsoleOutput() = %v, want %v", got, tt.WantedRecording.Console)
			}
			if got := o.ErrorOutput(); got != tt.WantedRecording.Error {
				t.Errorf("NewRecorder().ErrorOutput() = %v, want %v", got, tt.WantedRecording.Error)
			}
			if got := o.LogOutput(); got != tt.WantedRecording.Log {
				t.Errorf("NewRecorder().LogOutput() = %v, want %v", got, tt.WantedRecording.Log)
			}
			if issues, verified := o.Verify(tt.WantedRecording); !verified {
				for _, issue := range issues {
					t.Errorf("NewRecorder() %s", issue)
				}
			}
		})
	}
}

func TestRecorder_Verify(t *testing.T) {
	type args struct {
		o *output.Recorder
		w output.WantedRecording
	}
	tests := map[string]struct {
		args
		wantDifferences []string
		wantVerified    bool
	}{
		"normal": {
			args: args{
				o: output.NewRecorder(),
				w: output.WantedRecording{},
			},
			wantVerified: true,
		},
		"errors": {
			args: args{
				o: output.NewRecorder(),
				w: output.WantedRecording{
					Console: "unexpected console output",
					Error:   "unexpected error output",
					Log:     "unexpected log output",
				},
			},
			wantDifferences: []string{
				"console output = \"\", want \"unexpected console output\"",
				"error output = \"\", want \"unexpected error output\"",
				"log output = \"\", want \"unexpected log output\"",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotDifferences, gotVerified := tt.args.o.Verify(tt.args.w)
			if !reflect.DeepEqual(gotDifferences, tt.wantDifferences) {
				t.Errorf("Recorder.Verify() gotIssues = %v, want %v", gotDifferences, tt.wantDifferences)
			}
			if gotVerified != tt.wantVerified {
				t.Errorf("Recorder.Verify() gotVerified = %v, want %v", gotVerified, tt.wantVerified)
			}
			vr := newVerificationReporter()
			tt.args.o.Report(vr, "Recorder.Verify()", tt.args.w)
			for i, line := range vr.buffer {
				wanted := "Recorder.Verify() " + tt.wantDifferences[i]
				if !strings.HasSuffix(line, wanted) {
					t.Errorf("Recorder.Verify() recorded %q wanted %q", line, wanted)
				}
			}
		})
		tType := reflect.TypeOf(t)
		interfaceType := reflect.TypeOf((*output.TestingReporter)(nil)).Elem()
		if !tType.Implements(interfaceType) {
			t.Errorf("Recorder.Verify() *testing.T does not implement TestingReporter")
		}
	}
}

type verificationReporter struct {
	buffer []string
}

func newVerificationReporter() *verificationReporter {
	return &verificationReporter{buffer: []string{}}
}

func (vr *verificationReporter) Errorf(format string, data ...any) {
	vr.buffer = append(vr.buffer, fmt.Sprintf(format, data...))
}

func TestRecorder_Log(t *testing.T) {
	type args struct {
		l    output.Level
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
				l:    output.Trace,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='trace' f='v' msg='hello'\n",
		},
		"debug": {
			args: args{
				l:    output.Debug,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='debug' f='v' msg='hello'\n",
		},
		"info": {
			args: args{
				l:    output.Info,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='info' f='v' msg='hello'\n",
		},
		"warning": {
			args: args{
				l:    output.Warning,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='warning' f='v' msg='hello'\n",
		},
		"error": {
			args: args{
				l:    output.Error,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='error' f='v' msg='hello'\n",
		},
		"panic": {
			args: args{
				l:    output.Panic,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='panic' f='v' msg='hello'\n",
		},
		"fatal": {
			args: args{
				l:    output.Fatal,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Log: "level='fatal' f='v' msg='hello'\n",
		},
		"illegal": {
			args: args{
				l:    output.Trace + 1,
				msg:  "hello",
				args: map[string]any{"f": "v"},
			},
			Error: "" +
				"Programming error: call to Recorder.Log() with invalid level value 7; " +
				"message: 'hello', " +
				"args: 'map[f:v]'.\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := output.NewRecorder()
			r.Log(tt.args.l, tt.args.msg, tt.args.args)
			if got := r.LogOutput(); got != tt.Log {
				t.Errorf("Recorder.Log() got log %q want %q", got, tt.Log)
			}
			if got := r.ErrorOutput(); got != tt.Error {
				t.Errorf("Recorder.Log() got error %q want %q", got, tt.Error)
			}
		})
	}
}

func TestRecordingLogger_Trace(t *testing.T) {
	type args struct {
		msg    string
		fields map[string]any
	}
	tests := map[string]struct {
		name string
		tl   *output.RecordingLogger
		args
		want string
	}{
		"simple test": {
			tl: output.NewRecordingLogger(),
			args: args{
				msg:    "simple message",
				fields: map[string]any{"f1": 1, "f2": true, "f3": "v"},
			},
			want: "level='trace' f1='1' f2='true' f3='v' msg='simple message'\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.tl.Trace(tt.args.msg, tt.args.fields)
			if got := tt.tl.String(); got != tt.want {
				t.Errorf("RecordingLogger.Trace() got %q want %q", got, tt.want)
			}
		})
	}
}

func TestRecordingLogger_Debug(t *testing.T) {
	type args struct {
		msg    string
		fields map[string]any
	}
	tests := map[string]struct {
		tl *output.RecordingLogger
		args
		want string
	}{
		"simple test": {
			tl: output.NewRecordingLogger(),
			args: args{
				msg:    "simple message",
				fields: map[string]any{"f1": 1, "f2": true, "f3": "v"},
			},
			want: "level='debug' f1='1' f2='true' f3='v' msg='simple message'\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.tl.Debug(tt.args.msg, tt.args.fields)
			if got := tt.tl.String(); got != tt.want {
				t.Errorf("RecordingLogger.Debug() got %q want %q", got, tt.want)
			}
		})
	}
}

func TestRecordingLogger_Info(t *testing.T) {
	type args struct {
		msg    string
		fields map[string]any
	}
	tests := map[string]struct {
		tl *output.RecordingLogger
		args
		want string
	}{
		"simple test": {
			tl: output.NewRecordingLogger(),
			args: args{
				msg:    "simple message",
				fields: map[string]any{"f1": 1, "f2": true, "f3": "v"},
			},
			want: "level='info' f1='1' f2='true' f3='v' msg='simple message'\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.tl.Info(tt.args.msg, tt.args.fields)
			if got := tt.tl.String(); got != tt.want {
				t.Errorf("RecordingLogger.Info() got %q want %q", got, tt.want)
			}
		})
	}
}

func TestRecordingLogger_Warning(t *testing.T) {
	type args struct {
		msg    string
		fields map[string]any
	}
	tests := map[string]struct {
		tl *output.RecordingLogger
		args
		want string
	}{
		"simple test": {
			tl: output.NewRecordingLogger(),
			args: args{
				msg:    "simple message",
				fields: map[string]any{"f1": 1, "f2": true, "f3": "v"},
			},
			want: "level='warning' f1='1' f2='true' f3='v' msg='simple message'\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.tl.Warning(tt.args.msg, tt.args.fields)
			if got := tt.tl.String(); got != tt.want {
				t.Errorf("RecordingLogger.Warning() got %q want %q", got, tt.want)
			}
		})
	}
}

func TestRecordingLogger_Error(t *testing.T) {
	type args struct {
		msg    string
		fields map[string]any
	}
	tests := map[string]struct {
		tl *output.RecordingLogger
		args
		want string
	}{
		"simple test": {
			tl: output.NewRecordingLogger(),
			args: args{
				msg:    "simple message",
				fields: map[string]any{"f1": 1, "f2": true, "f3": "v"},
			},
			want: "level='error' f1='1' f2='true' f3='v' msg='simple message'\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.tl.Error(tt.args.msg, tt.args.fields)
			if got := tt.tl.String(); got != tt.want {
				t.Errorf("RecordingLogger.Error() got %q want %q", got, tt.want)
			}
		})
	}
}

func TestRecordingLogger_Panic(t *testing.T) {
	type args struct {
		msg    string
		fields map[string]any
	}
	tests := map[string]struct {
		tl *output.RecordingLogger
		args
		want string
	}{
		"simple test": {
			tl: output.NewRecordingLogger(),
			args: args{
				msg:    "simple message",
				fields: map[string]any{"f1": 1, "f2": true, "f3": "v"},
			},
			want: "level='panic' f1='1' f2='true' f3='v' msg='simple message'\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.tl.Panic(tt.args.msg, tt.args.fields)
			if got := tt.tl.String(); got != tt.want {
				t.Errorf("RecordingLogger.Panic() got %q want %q", got, tt.want)
			}
		})
	}
}

func TestRecordingLogger_Fatal(t *testing.T) {
	type args struct {
		msg    string
		fields map[string]any
	}
	tests := map[string]struct {
		tl *output.RecordingLogger
		args
		want string
	}{
		"simple test": {
			tl: output.NewRecordingLogger(),
			args: args{
				msg:    "simple message",
				fields: map[string]any{"f1": 1, "f2": true, "f3": "v"},
			},
			want: "level='fatal' f1='1' f2='true' f3='v' msg='simple message'\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.tl.Fatal(tt.args.msg, tt.args.fields)
			if got := tt.tl.String(); got != tt.want {
				t.Errorf("RecordingLogger.Fatal() got %q want %q", got, tt.want)
			}
		})
	}
}

func TestRecorder_IsConsoleTTY(t *testing.T) {
	tests := map[string]struct {
		r    *output.Recorder
		want bool
	}{"simple": {r: output.NewRecorder(), want: false}}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.r.IsConsoleTTY(); got != tt.want {
				t.Errorf("Recorder.IsConsoleTTY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecorder_IsErrorTTY(t *testing.T) {
	tests := map[string]struct {
		r    *output.Recorder
		want bool
	}{"simple": {r: output.NewRecorder(), want: false}}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.r.IsErrorTTY(); got != tt.want {
				t.Errorf("Recorder.IsErrorTTY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Recorder_IncrementTab(t *testing.T) {
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
			r := &output.Recorder{}
			r.IncrementTab(tt.initialTab)
			r.IncrementTab(tt.t)
			if got := r.Tab(); got != tt.want {
				t.Errorf("Recorder.IncrementTab got %d want %d", got, tt.want)
			}
		})
	}
}

func Test_Recorder_DecrementTab(t *testing.T) {
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
			r := &output.Recorder{}
			r.IncrementTab(tt.initialTab)
			r.DecrementTab(tt.t)
			if got := r.Tab(); got != tt.want {
				t.Errorf("Recorder.DecrementTab got %d want %d", got, tt.want)
			}
		})
	}
}

func TestRecorder_BeginConsoleList(t *testing.T) {
	tests := map[string]struct {
		numeric bool
		want    string
	}{
		"bullet": {
			numeric: false,
			want:    "● ",
		},
		"numeric": {
			numeric: true,
			want:    " 1. ",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := output.NewRecorder()
			r.BeginConsoleList(tt.numeric)
			if got := r.ConsoleListDecorator().Decorator(); got != tt.want {
				t.Errorf("BeginConsoleList() = %v, want %v", got, tt.want)
			}
			r.EndConsoleList()
			if got := r.ConsoleListDecorator().Decorator(); got != "" {
				t.Errorf("EndConsoleList() = %v, want %v", got, "")
			}
		})
	}
}

func TestRecorder_BeginErrorList(t *testing.T) {
	tests := map[string]struct {
		numeric bool
		want    string
	}{
		"bullet": {
			numeric: false,
			want:    "● ",
		},
		"numeric": {
			numeric: true,
			want:    " 1. ",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := output.NewRecorder()
			r.BeginErrorList(tt.numeric)
			if got := r.ErrorListDecorator().Decorator(); got != tt.want {
				t.Errorf("BeginErrorList() = %v, want %v", got, tt.want)
			}
			r.EndErrorList()
			if got := r.ErrorListDecorator().Decorator(); got != "" {
				t.Errorf("EndErrorList() = %v, want %v", got, "")
			}
		})
	}
}
