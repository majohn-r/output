package output

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

type (
	// Level is used to specify log levels for Bus.Log().
	Level uint32

	// Bus defines a set of functions for writing console messages and error
	// messages, and for providing access to the console writer, the error
	// writer, and a Logger instance; its primary use is to simplify how
	// application code handles console, error, and logged output, and its
	// secondary use is to make it easy to test output writing.
	Bus interface {
		Log(Level, string, map[string]any)
		WriteCanonicalConsole(string, ...any)
		WriteConsole(string, ...any)
		WriteCanonicalError(string, ...any)
		WriteError(string, ...any)
		ConsoleWriter() io.Writer
		ErrorWriter() io.Writer
		IsConsoleTTY() bool
		IsErrorTTY() bool
		// Tab returns the current tab setting (number of spaces)
		Tab() uint8
		// IncrementTab increases the current tab setting up to the max uint8 value
		IncrementTab(uint8)
		// DecrementTab decreases the current tab setting; will not go below 0
		DecrementTab(uint8)
	}

	// Logger defines a set of functions for writing to a log at various log
	// levels
	Logger interface {
		Trace(msg string, fields map[string]any)
		Debug(msg string, fields map[string]any)
		Info(msg string, fields map[string]any)
		Warning(msg string, fields map[string]any)
		Error(msg string, fields map[string]any)
		Panic(msg string, fields map[string]any)
		Fatal(msg string, fields map[string]any)
	}

	bus struct {
		consoleWriter io.Writer
		errorWriter   io.Writer
		logger        Logger
		performWrites bool
		consoleTTY    bool
		errorTTY      bool
		tab           uint8
	}
)

// These are the different logging levels.
const (
	Fatal Level = iota
	Panic
	Error
	Warning
	Info
	Debug
	Trace
)

// NewDefaultBus returns an implementation of Bus that writes console messages
// to stdout and error messages to stderr.
func NewDefaultBus(l Logger) Bus {
	return NewCustomBus(os.Stdout, os.Stderr, l)
}

// vars so testing can replace
var (
	isTerminal       = isatty.IsTerminal
	isCygwinTerminal = isatty.IsCygwinTerminal
)

func isTTY(w io.Writer) (b bool) {
	if f, ok := w.(*os.File); ok {
		fd := f.Fd()
		b = isTerminal(fd) || isCygwinTerminal(fd)
	}
	return
}

// NewCustomBus returns an implementation of Bus that lets the caller specify
// the console and error writers and the Logger.
func NewCustomBus(c, e io.Writer, l Logger) Bus {
	return &bus{
		consoleWriter: c,
		errorWriter:   e,
		logger:        l,
		performWrites: true,
		consoleTTY:    isTTY(c),
		errorTTY:      isTTY(e),
		tab:           0,
	}
}

// Log logs a message and map of fields at a specified log level.
func (b *bus) Log(l Level, msg string, args map[string]any) {
	if b.performWrites {
		switch l {
		case Trace:
			b.logger.Trace(msg, args)
		case Debug:
			b.logger.Debug(msg, args)
		case Info:
			b.logger.Info(msg, args)
		case Warning:
			b.logger.Warning(msg, args)
		case Error:
			b.logger.Error(msg, args)
		case Panic:
			b.logger.Panic(msg, args)
		case Fatal:
			b.logger.Fatal(msg, args)
		default:
			b.WriteCanonicalError("programming error: call to bus.Log() with invalid level value %d; message: '%s', args: '%v", l, msg, args)
		}
	}
}

// ConsoleWriter returns a writer for console output.
func (b *bus) ConsoleWriter() io.Writer {
	return b.consoleWriter
}

// ErrorWriter returns a writer for error output.
func (b *bus) ErrorWriter() io.Writer {
	return b.errorWriter
}

// WriteCanonicalError writes error output in a canonical format.
func (b *bus) WriteCanonicalError(format string, a ...any) {
	if b.performWrites {
		_, _ = fmt.Fprint(b.errorWriter, canonicalFormat(format, a...))
	}
}

// WriteError writes unedited error output.
func (b *bus) WriteError(format string, a ...any) {
	if b.performWrites {
		fmt.Fprintf(b.errorWriter, format, a...)
	}
}

// WriteCanonicalConsole writes output to a console in a canonical format.
func (b *bus) WriteCanonicalConsole(format string, a ...any) {
	if b.performWrites {
		writeTabbedContent(b.consoleWriter, b.tab, canonicalFormat(format, a...))
	}
}

// WriteConsole writes output to a console.
func (b *bus) WriteConsole(format string, a ...any) {
	if b.performWrites {
		writeTabbedContent(b.consoleWriter, b.tab, fmt.Sprintf(format, a...))
	}
}

func writeTabbedContent(w io.Writer, tab uint8, content string) {
	fmt.Fprintf(w, "%*s%s", tab, "", content)
}

// IncrementTab increments the tab setting by the specified number of spaces
func (b *bus) IncrementTab(t uint8) {
	b.tab = addTab(b.tab, t)
}

func addTab(originalTab, increment uint8) uint8 {
	if originalTab+increment > originalTab {
		return originalTab + increment
	}
	return originalTab
}

// DecrementTab decrements the tab setting by the specified number of spaces
func (b *bus) DecrementTab(t uint8) {
	b.tab = subtractTab(b.tab, t)
}

func subtractTab(originalTab, decrement uint8) uint8 {
	if originalTab >= decrement {
		return originalTab - decrement
	}
	return originalTab
}

// Tab returns the current tab setting
func (b *bus) Tab() uint8 {
	return b.tab
}

// IsConsoleTTY returns whether the console writer is a TTY
func (b *bus) IsConsoleTTY() bool {
	return b.consoleTTY
}

// IsErrorTTY returns whether the error writer is a TTY
func (b *bus) IsErrorTTY() bool {
	return b.errorTTY
}
