package output

import (
	"fmt"
	"strings"
	"unicode"
)

const sentenceTerminatingCharacters = ".!?:"

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
		s = strings.TrimRight(s, sentenceTerminatingCharacters)
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
	return strings.TrimRight(s, "\n")
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
	for _, r := range s {
		if strings.ContainsRune(sentenceTerminatingCharacters, r) {
			return true
		}
	}
	return false
}
