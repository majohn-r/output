package output

import (
	"testing"
)

func Test_canonicalFormat(t *testing.T) {
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
			args: args{format: "\n\n\n\n\n\n\n\n\n\n\n\n\n\n"},
			want: "\n",
		},
		"nothing but terminal punctuation": {
			args: args{format: "!!?.!?."},
			want: ".\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := canonicalFormat(tt.args.format, tt.args.a...); got != tt.want {
				t.Errorf("canonicalFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_renderCanonical(t *testing.T) {
	tests := map[string]struct {
		s    string
		want string
	}{
		"broad test": {
			s:    "test format 25 \"foo\" 1.245..?!..?\n\n\n\n",
			want: "Test format 25 \"foo\" 1.245?\n",
		},
		"narrow test": {
			s:    "1. test format 25 \"foo\" 1.245",
			want: "1. test format 25 \"foo\" 1.245.\n",
		},
		"nothing but newlines": {
			s:    "\n\n\n\n\n\n\n\n\n\n\n\n\n\n",
			want: "\n",
		},
		"nothing but terminal punctuation": {
			s:    "!!?.!?.",
			want: ".\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := renderCanonical(tt.s); got != tt.want {
				t.Errorf("renderCanonical() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fixTerminalPunctuation(t *testing.T) {
	tests := map[string]struct {
		s    string
		want string
	}{
		"broad test": {
			s:    "test format 25 \"foo\" 1.245..?!..?",
			want: "test format 25 \"foo\" 1.245?",
		},
		"narrow test": {
			s:    "1. test format 25 \"foo\" 1.245",
			want: "1. test format 25 \"foo\" 1.245.",
		},
		"nothing": {
			s:    "",
			want: "",
		},
		"nothing but terminal punctuation": {
			s:    "!!?.!?.",
			want: ".",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := fixTerminalPunctuation(tt.s); got != tt.want {
				t.Errorf("fixTerminalPunctuation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stripTrailingNewlines(t *testing.T) {
	tests := map[string]struct {
		s    string
		want string
	}{
		"broad test": {
			s:    "test format 25 \"foo\" 1.245..?!..?\n\n\n\n",
			want: "test format 25 \"foo\" 1.245..?!..?",
		},
		"narrow test": {
			s:    "1. test format 25 \"foo\" 1.245",
			want: "1. test format 25 \"foo\" 1.245",
		},
		"nothing but newlines": {
			s:    "\n\n\n\n\n\n\n\n\n\n\n\n\n\n",
			want: "",
		},
		"nothing": {
			s:    "",
			want: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := stripTrailingNewlines(tt.s); got != tt.want {
				t.Errorf("stripTrailingNewlines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_capitalize(t *testing.T) {
	tests := map[string]struct {
		s    string
		want string
	}{
		"broad test": {
			s:    "test format 25 \"foo\" 1.245..?!..?\n\n\n\n",
			want: "Test format 25 \"foo\" 1.245..?!..?\n\n\n\n",
		},
		"narrow test": {
			s:    "1. test format 25 \"foo\" 1.245",
			want: "1. test format 25 \"foo\" 1.245",
		},
		"nothing but newlines": {
			s:    "\n\n\n\n\n\n\n\n\n\n\n\n\n\n",
			want: "\n\n\n\n\n\n\n\n\n\n\n\n\n\n",
		},
		"nothing": {
			s:    "",
			want: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := capitalize(tt.s); got != tt.want {
				t.Errorf("capitalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSentenceTerminatingPunctuation(t *testing.T) {
	tests := map[string]struct {
		s    string
		want bool
	}{
		"period":               {s: ".", want: true},
		"question mark":        {s: "?", want: true},
		"exclamation point":    {s: "!", want: true},
		"comma":                {s: ",", want: false},
		"semicolon":            {s: ";", want: false},
		"colon":                {s: ":", want: false},
		"en dash":              {s: "\u2013", want: false},
		"em dash":              {s: "—", want: false},
		"hyphen":               {s: "-", want: false},
		"left parenthesis":     {s: "(", want: false},
		"right parenthesis":    {s: ")", want: false},
		"left bracket":         {s: "[", want: false},
		"right bracket":        {s: "]", want: false},
		"left brace":           {s: "{", want: false},
		"right brace":          {s: "}", want: false},
		"apostrophe":           {s: "'", want: false},
		"plain quotation mark": {s: "\"", want: false},
		"ellipsis":             {s: "…", want: false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := isSentenceTerminatingPunctuation(tt.s); got != tt.want {
				t.Errorf("isSentenceTerminatingPunctuation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lastCharacter(t *testing.T) {
	tests := map[string]struct {
		s            string
		wantResultS  string
		wantLastChar string
	}{
		"empty": {
			s:            "",
			wantResultS:  "",
			wantLastChar: "",
		},
		"nearly empty": {
			s:            "a",
			wantResultS:  "",
			wantLastChar: "",
		},
		"few": {
			s:            "ab",
			wantResultS:  "a",
			wantLastChar: "a",
		},
		"lots": {
			s:            "gjdfhgdsfhgiudsfhgiuldshgiu",
			wantResultS:  "gjdfhgdsfhgiudsfhgiuldshgi",
			wantLastChar: "i",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResultS, gotLastChar := lastCharacter(tt.s)
			if gotResultS != tt.wantResultS {
				t.Errorf("lastCharacter() gotResultS = %v, want %v", gotResultS, tt.wantResultS)
			}
			if gotLastChar != tt.wantLastChar {
				t.Errorf("lastCharacter() gotLastChar = %v, want %v", gotLastChar, tt.wantLastChar)
			}
		})
	}
}
