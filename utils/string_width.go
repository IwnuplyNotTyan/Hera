package utils

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/text/width"
)

func StringWidth(s string) int {
	w := 0
	for _, r := range s {
		p := width.LookupRune(r)
		switch p.Kind() {
		case width.EastAsianWide, width.EastAsianFullwidth, width.EastAsianAmbiguous:
			w += 2
		default:
			w += 1
		}
	}
	return w
}

func RuneWidth(r rune) int {
	p := width.LookupRune(r)
	switch p.Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth, width.EastAsianAmbiguous:
		return 2
	default:
		return 1
	}
}

func PadString(s string, width int) string {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}
	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}

func TruncateString(s string, maxWidth int) string {
	w := 0
	result := make([]rune, 0, len(s))
	for _, r := range s {
		rw := RuneWidth(r)
		if w+rw > maxWidth {
			break
		}
		w += rw
		result = append(result, r)
	}
	return string(result)
}

func AlignCenter(s string, width int) string {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}
	leftPad := (width - currentWidth) / 2
	rightPad := width - currentWidth - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

func AlignRight(s string, width int) string {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}
	padding := width - currentWidth
	return strings.Repeat(" ", padding) + s
}

func ContainsWideChars(s string) bool {
	for _, r := range s {
		if RuneWidth(r) == 2 {
			return true
		}
	}
	return false
}

func ValidUTF8(s string) bool {
	return utf8.ValidString(s)
}
