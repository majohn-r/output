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
	terminalPunctuation := "."
	if isSentenceTerminatingPunctuation(lastChar) {
		terminalPunctuation = lastChar
		// remove the punctuation at the end
		s, lastChar = lastCharacter(s)
		// look for more and remove them too
		for s != "" && isSentenceTerminatingPunctuation(lastChar) {
			s, lastChar = lastCharacter(s)
		}
	}
	s += terminalPunctuation
	return s
}

func lastCharacter(s string) (resultS, lastChar string) {
	if s != "" {
		resultS = s[:len(s)-1]
		if resultS != "" {
			lastChar = resultS[len(resultS)-1:]
		}
	}
	return
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
