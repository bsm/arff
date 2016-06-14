package arff

import (
	"strings"
	"unicode/utf8"
)

const quoteRune = '\''

// Inspired by https://golang.org/src/strconv/quote.go
// Copyright 2009 The Go Authors. All rights reserved.
func quote(s string) string {
	s = strings.TrimSpace(s)
	if !stringNeedsQuotes(s) {
		return s
	}

	buf := make([]rune, 0, 3*len(s)/2) // avoid reallocations
	buf = append(buf, quoteRune)

	for _, r := range s {
		switch r {
		case quoteRune, '\\':
			buf = append(buf, '\\', r)
		case '\a':
			buf = append(buf, '\\', 'a')
		case '\b':
			buf = append(buf, '\\', 'b')
		case '\f':
			buf = append(buf, '\\', 'f')
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		case '\v':
			buf = append(buf, '\\', 'v')
		default:
			buf = append(buf, r)
		}
	}
	buf = append(buf, quoteRune)
	return string(buf)
}

func unquoteAll(sl []string) []string {
	for i, s := range sl {
		sl[i] = unquote(s)
	}
	return sl
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	l := len(s)
	if l < 2 || s[0] != quoteRune || s[l-1] != quoteRune {
		return s
	}

	buf := make([]rune, 0, l-2)

	for pos := 1; pos < l-1; {
		r1, w1 := utf8.DecodeRuneInString(s[pos:])
		pos += w1

		switch r1 {
		case '\\':
			r2, w2 := utf8.DecodeRuneInString(s[pos:])
			pos += w2

			switch r2 {
			case quoteRune, '\\':
				buf = append(buf, r2)
			case 'a':
				buf = append(buf, '\a')
			case 'b':
				buf = append(buf, '\b')
			case 'f':
				buf = append(buf, '\f')
			case 'n':
				buf = append(buf, '\n')
			case 'r':
				buf = append(buf, '\r')
			case 't':
				buf = append(buf, '\t')
			case 'v':
				buf = append(buf, '\v')
			default:
				buf = append(buf, r1, r2)
			}
		default:
			buf = append(buf, r1)
		}
	}
	return strings.TrimSpace(string(buf))
}

func stringNeedsQuotes(s string) bool {
	if s == "" {
		return false
	}
	if s == "?" {
		return true
	}

	for _, r := range s {
		switch r {
		case ' ', '"', '\'', '%', '{', '}', '\a', '\b', '\f', '\n', '\r', '\t', '\v':
			return true
		}
	}
	return false
}
