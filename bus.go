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
		// Deprecated: use ConsolePrintf or ConsolePrintln
		WriteCanonicalConsole(string, ...any)
		// Deprecated: use ConsolePrintf or ConsolePrintln
		WriteConsole(string, ...any)
		// Deprecated: use ErrorPrintf or ErrorPrintln
		WriteCanonicalError(string, ...any)
		// Deprecated: use ErrorPrintf or ErrorPrintln
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
		BeginConsoleList(bool)
		EndConsoleList()
		ConsoleListDecorator() *ListDecorator
		BeginErrorList(bool)
		EndErrorList()
		ErrorListDecorator() *ListDecorator
		ConsolePrintf(string, ...any)
		ConsolePrintln(string)
		ErrorPrintf(string, ...any)
		ErrorPrintln(string)
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
		consoleWriter        io.Writer
		errorWriter          io.Writer
		logger               Logger
		performWrites        bool
		consoleTTY           bool
		errorTTY             bool
		tab                  uint8
		consoleListDecorator *ListDecorator
		errorListDecorator   *ListDecorator
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
		consoleWriter:        c,
		errorWriter:          e,
		logger:               l,
		performWrites:        true,
		consoleTTY:           isTTY(c),
		errorTTY:             isTTY(e),
		tab:                  0,
		consoleListDecorator: newListDecorator(false, false),
		errorListDecorator:   newListDecorator(false, false),
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
			b.WriteCanonicalError(
				"programming error: call to bus.Log() with invalid level value %d; message: '%s', args: '%v",
				l,
				msg,
				args,
			)
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

// ErrorPrintln prints a message to the error channel, terminated by a newline
func (b *bus) ErrorPrintln(msg string) {
	if b.performWrites {
		doPrintln(b.errorWriter, b.errorListDecorator, msg)
	}
}

// ErrorPrintf prints a message with arguments to the error channel
func (b *bus) ErrorPrintf(format string, args ...any) {
	if b.performWrites {
		doPrintf(b.errorWriter, b.errorListDecorator, format, args...)
	}
}

// ConsolePrintln prints a message to the error channel, terminated by a newline
func (b *bus) ConsolePrintln(msg string) {
	if b.performWrites {
		writeTabbedContent(b.consoleWriter, b.tab, doSprintln(b.consoleListDecorator, msg))
	}
}

// ConsolePrintf prints a message with arguments to the error channel
func (b *bus) ConsolePrintf(format string, args ...any) {
	if b.performWrites {
		writeTabbedContent(b.consoleWriter, b.tab, doSprintf(b.consoleListDecorator, format, args...))
	}
}

func doPrintf(w io.Writer, decorator *ListDecorator, format string, args ...any) {
	_, _ = fmt.Fprint(w, doSprintf(decorator, format, args...))
}

func doSprintf(decorator *ListDecorator, format string, args ...any) string {
	return fmt.Sprintf("%s%s", decorator.Decorator(), fmt.Sprintf(format, args...))
}

func doPrintln(w io.Writer, decorator *ListDecorator, msg string) {
	_, _ = fmt.Fprint(w, doSprintln(decorator, msg))
}

func doSprintln(decorator *ListDecorator, msg string) string {
	return fmt.Sprintf("%s%s\n", decorator.Decorator(), msg)
}

// WriteCanonicalError writes error output in a canonical format.
func (b *bus) WriteCanonicalError(format string, a ...any) {
	if b.performWrites {
		_, _ = fmt.Fprint(b.errorWriter, b.errorListDecorator.Decorator()+canonicalFormat(format, a...))
	}
}

// WriteError writes unedited error output.
func (b *bus) WriteError(format string, a ...any) {
	if b.performWrites {
		fmt.Fprintf(b.errorWriter, b.errorListDecorator.Decorator()+format, a...)
	}
}

// WriteCanonicalConsole writes output to a console in a canonical format.
func (b *bus) WriteCanonicalConsole(format string, a ...any) {
	if b.performWrites {
		writeTabbedContent(b.consoleWriter, b.tab, b.consoleListDecorator.Decorator()+canonicalFormat(format, a...))
	}
}

// WriteConsole writes output to a console.
func (b *bus) WriteConsole(format string, a ...any) {
	if b.performWrites {
		writeTabbedContent(b.consoleWriter, b.tab, fmt.Sprintf(b.consoleListDecorator.Decorator()+format, a...))
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

// BeginConsoleList initiates console listing
func (b *bus) BeginConsoleList(numeric bool) {
	b.consoleListDecorator = newListDecorator(true, numeric)
}

// EndConsoleList terminates console listing
func (b *bus) EndConsoleList() {
	b.consoleListDecorator = newListDecorator(false, false)
}

// ConsoleListDecorator makes the console list decorator available
func (b *bus) ConsoleListDecorator() *ListDecorator {
	return b.consoleListDecorator
}

// BeginErrorList initiates error listing
func (b *bus) BeginErrorList(numeric bool) {
	b.errorListDecorator = newListDecorator(true, numeric)
}

// EndErrorList terminates error listing
func (b *bus) EndErrorList() {
	b.errorListDecorator = newListDecorator(false, false)
}

// ErrorListDecorator makes the error list decorator available
func (b *bus) ErrorListDecorator() *ListDecorator {
	return b.errorListDecorator
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
