package output

type (

	// NilWriter is a writer that does nothing at all; its intended use is for
	// test code where the side effect of writing to the console or writing
	// error output is of no interest whatsoever.
	NilWriter struct{}

	// NilLogger is a logger that does nothing at all; its intended use is for
	// test code where the side effect of logging is of no interest whatsoever.
	NilLogger struct{}
)

// NewNilBus returns an implementation of Bus that does nothing.
func NewNilBus() Bus {
	nw := NilWriter{}
	return &bus{
		consoleWriter: nw,
		errorWriter:   nw,
		logger:        NilLogger{},
		performWrites: false,
	}
}

// Trace does nothing.
func (nl NilLogger) Trace(_ string, _ map[string]any) {
}

// Debug does nothing.
func (nl NilLogger) Debug(_ string, _ map[string]any) {
}

// Info does nothing.
func (nl NilLogger) Info(_ string, _ map[string]any) {
}

// Warning does nothing.
func (nl NilLogger) Warning(_ string, _ map[string]any) {
}

// Error does nothing.
func (nl NilLogger) Error(_ string, _ map[string]any) {
}

// Panic does nothing.
func (nl NilLogger) Panic(_ string, _ map[string]any) {
}

// Fatal does nothing.
func (nl NilLogger) Fatal(_ string, _ map[string]any) {
}

// Write does nothing except return the expected values
func (nw NilWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
