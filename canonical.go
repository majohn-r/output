package output

import (
	"fmt"
	"unicode"
)

func canonicalFormat(format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	return renderCanonical(s)
}

func renderCanonical(s string) string {
	s = stripTrailingNewlines(s)
	s = fixTerminalPunctuation(s)
	s = capitalize(s)
	return s + "\n"
}

func fixTerminalPunctuation(s string) string {
	if s == "" {
		return s
	}
	lastChar := s[len(s)-1:]
	terminalPunctuation := lastChar
	if !isSentenceTerminatingPunctuation(lastChar) {
		terminalPunctuation = "."
	} else {
		// remove the punctuation at the end
		s = s[:len(s)-1]
		if s != "" {
			lastChar = s[len(s)-1:]
		}
		// look for more and remove them too
		for s != "" && isSentenceTerminatingPunctuation(lastChar) {
			s = s[:len(s)-1]
			if s != "" {
				lastChar = s[len(s)-1:]
			}
		}
	}
	s += terminalPunctuation
	return s
}

func stripTrailingNewlines(s string) string {
	for s != "" && s[len(s)-1:] == "\n" {
		s = s[:len(s)-1]
	}
	return s
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if unicode.IsLower(r[0]) {
		r[0] = unicode.ToUpper(r[0])
		s = string(r)
	}
	return s
}

func isSentenceTerminatingPunctuation(s string) bool {
	switch s {
	case ".", "!", "?":
		return true
	}
	return false
}
