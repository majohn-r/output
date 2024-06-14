package output

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"
)

type (
	// Recorder is an implementation of Bus that simply records its inputs; it's
	// intended for unit tests, where you can provide the code under test (that
	// needs a Bus) with an instance of Recorder and then verify that the code
	// produces the expected console, error, and log output.
	Recorder struct {
		consoleWriter *bytes.Buffer
		errorWriter   *bytes.Buffer
		logger        *RecordingLogger
		tab           uint8
	}

	// WantedRecording is intended to be used in unit tests as part of the test
	// structure; it allows the test writer to capture what the test wants the
	// console, error, and log output to contain.
	WantedRecording struct {
		Console string
		Error   string
		Log     string
	}

	// RecordingLogger is a simple logger intended for use in unit tests; it
	// records the output given to it.
	//
	// Caveats:
	//
	// Your production log may not actually do anything with some calls into it
	// - for instance, many logging frameworks allow you to limit the severity
	// of what is logged, e.g., only warnings or worse; RecordingLogger will
	// record every call made into it.
	//
	// The output recorded cannot be guaranteed to match exactly what your
	// logging code records - but it will include the log level, the message,
	// and all field-value pairs.
	//
	// The RecordingLogger will probably behave differently to a logging
	// mechanism that supports panic and fatal logs, in that a production logger
	// will probably call panic in processing a panic log, and will probably
	// exit the program on a fatal log. RecordingLogger does neither of those.
	RecordingLogger struct {
		writer *bytes.Buffer
	}
)

// NewRecorder returns a recording implementation of Bus.
func NewRecorder() *Recorder {
	return &Recorder{
		consoleWriter: &bytes.Buffer{},
		errorWriter:   &bytes.Buffer{},
		logger:        NewRecordingLogger(),
		tab:           0,
	}
}

// Log records a message and map of fields at a specified log level.
func (r *Recorder) Log(l Level, msg string, fields map[string]any) {
	switch l {
	case Trace:
		r.logger.Trace(msg, fields)
	case Debug:
		r.logger.Debug(msg, fields)
	case Info:
		r.logger.Info(msg, fields)
	case Warning:
		r.logger.Warning(msg, fields)
	case Error:
		r.logger.Error(msg, fields)
	case Panic:
		r.logger.Panic(msg, fields)
	case Fatal:
		r.logger.Fatal(msg, fields)
	default:
		r.WriteCanonicalError("programming error: call to Recorder.Log() with invalid level value %d; message: '%s', args: '%v", l, msg, fields)
	}
}

// ConsoleWriter returns the internal console writer.
func (r *Recorder) ConsoleWriter() io.Writer {
	return r.consoleWriter
}

// ErrorWriter returns the internal error writer.
func (r *Recorder) ErrorWriter() io.Writer {
	return r.errorWriter
}

// WriteCanonicalError records data written as an error.
func (r *Recorder) WriteCanonicalError(format string, a ...any) {
	fmt.Fprint(r.errorWriter, canonicalFormat(format, a...))
}

// WriteError records un-edited data written as an error.
func (r *Recorder) WriteError(format string, a ...any) {
	fmt.Fprintf(r.errorWriter, format, a...)
}

// WriteCanonicalConsole records data written to the console.
func (r *Recorder) WriteCanonicalConsole(format string, a ...any) {
	writeTabbedContent(r.consoleWriter, r.tab, canonicalFormat(format, a...))
}

// WriteConsole records data written to the console.
func (r *Recorder) WriteConsole(format string, a ...any) {
	writeTabbedContent(r.consoleWriter, r.tab, fmt.Sprintf(format, a...))
}

// IncrementTab increments the tab setting by the specified number of spaces
func (r *Recorder) IncrementTab(t uint8) {
	r.tab = addTab(r.tab, t)
}

// DecrementTab decrements the tab setting by the specified number of spaces
func (r *Recorder) DecrementTab(t uint8) {
	r.tab = subtractTab(r.tab, t)
}

// Tab returns the current tab setting
func (r *Recorder) Tab() uint8 {
	return r.tab
}

// ConsoleOutput returns the data written as console output.
func (r *Recorder) ConsoleOutput() string {
	return r.consoleWriter.String()
}

// ErrorOutput returns the data written as error output.
func (r *Recorder) ErrorOutput() string {
	return r.errorWriter.String()
}

// LogOutput returns the data written to a log.
func (r *Recorder) LogOutput() string {
	return r.logger.writer.String()
}

// IsConsoleTTY returns whether the console writer is a TTY
func (r *Recorder) IsConsoleTTY() bool {
	return false
}

// IsErrorTTY returns whether the error writer is a TTY
func (r *Recorder) IsErrorTTY() bool {
	return false
}

// Verify verifies the recorded output against the expected output and returns
// any differences found.
func (r *Recorder) Verify(w WantedRecording) (differences []string, verified bool) {
	verified = true
	if got := r.ConsoleOutput(); got != w.Console {
		differences = append(differences, fmt.Sprintf("console output = %q, want %q", got, w.Console))
		verified = false
	}
	if got := r.ErrorOutput(); got != w.Error {
		differences = append(differences, fmt.Sprintf("error output = %q, want %q", got, w.Error))
		verified = false
	}
	if got := r.LogOutput(); got != w.Log {
		differences = append(differences, fmt.Sprintf("log output = %q, want %q", got, w.Log))
		verified = false
	}
	return
}

// TestingReporter is an interface that requires the one *testing.T function
// that we care about: Errorf
type TestingReporter interface {
	Errorf(string, ...any)
}

// Report handles the common use case for using a Recorder: detecting whether
// any differences were recorded, and reporting them if there were any
// differences.
func (r *Recorder) Report(t TestingReporter, header string, w WantedRecording) {
	if differences, verified := r.Verify(w); !verified {
		var location string
		if _, file, line, ok := runtime.Caller(1); ok {
			canonicalFile := strings.ReplaceAll(file, "/", "\\")
			location = fmt.Sprintf("called from %s:%d: ", canonicalFile, line)
		}
		for _, difference := range differences {
			t.Errorf("%s%s %s", location, header, difference)
			location = ""
		}
	}
}

// NewRecordingLogger returns a recording implementation of Logger.
func NewRecordingLogger() *RecordingLogger {
	return &RecordingLogger{writer: &bytes.Buffer{}}
}

// Trace records a trace log message.
func (rl *RecordingLogger) Trace(msg string, fields map[string]any) {
	rl.log("trace", msg, fields)
}

// Debug records a debug log message.
func (rl *RecordingLogger) Debug(msg string, fields map[string]any) {
	rl.log("debug", msg, fields)
}

// Info records an info log message.
func (rl *RecordingLogger) Info(msg string, fields map[string]any) {
	rl.log("info", msg, fields)
}

// Warning records a warning log message.
func (rl *RecordingLogger) Warning(msg string, fields map[string]any) {
	rl.log("warning", msg, fields)
}

// Error records an error log message.
func (rl *RecordingLogger) Error(msg string, fields map[string]any) {
	rl.log("error", msg, fields)
}

// Panic records a panic log message and does not call panic().
func (rl *RecordingLogger) Panic(msg string, fields map[string]any) {
	rl.log("panic", msg, fields)
}

// Fatal records a fatal log message and does not terminate the program.
func (rl *RecordingLogger) Fatal(msg string, fields map[string]any) {
	rl.log("fatal", msg, fields)
}

func (rl *RecordingLogger) log(level, msg string, fields map[string]any) {
	parts := make([]string, 0, len(fields))
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s='%v'", k, v))
	}
	sort.Strings(parts)
	fmt.Fprintf(rl.writer, "level='%s' %s msg='%s'\n", level, strings.Join(parts, " "), msg)
}
