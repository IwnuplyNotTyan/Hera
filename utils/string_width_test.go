package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringWidth_ASCII(t *testing.T) {
	assert.Equal(t, 5, StringWidth("hello"))
	assert.Equal(t, 0, StringWidth(""))
	assert.Equal(t, 10, StringWidth("0123456789"))
}

func TestStringWidth_WithWideChars(t *testing.T) {
	assert.Equal(t, 5, StringWidth("hello"))
	assert.Equal(t, 12, StringWidth("Привет"))
	assert.Equal(t, 15, StringWidth("Привет hi"))
	assert.Equal(t, 2, StringWidth("⽕"))
	assert.Equal(t, 5, StringWidth("⽕ U "))
}

func TestStringWidth_WithEmoji(t *testing.T) {
	assert.Equal(t, 2, StringWidth("♥"))
	assert.Equal(t, 2, StringWidth("♡"))
	assert.Equal(t, 2, StringWidth("●"))
	assert.Equal(t, 2, StringWidth("■"))
	assert.Equal(t, 8, StringWidth("♥ ♥ ♥"))
}

func TestRuneWidth(t *testing.T) {
	assert.Equal(t, 1, RuneWidth('a'))
	assert.Equal(t, 1, RuneWidth('1'))
	assert.Equal(t, 1, RuneWidth(' '))
	assert.Equal(t, 2, RuneWidth('⽕'))
	assert.Equal(t, 2, RuneWidth('♥'))
	assert.Equal(t, 2, RuneWidth('●'))
}

func TestPadString_ASCII(t *testing.T) {
	result := PadString("hello", 10)
	assert.Equal(t, "hello     ", result)
	assert.Equal(t, 10, StringWidth(result))
}

func TestPadString_WithWideChars(t *testing.T) {
	result := PadString("Привет", 12)
	assert.Equal(t, 12, StringWidth(result))
}

func TestPadString_WithEmoji(t *testing.T) {
	result := PadString("❤️", 6)
	assert.Equal(t, 6, StringWidth(result))
}

func TestPadString_AlreadyWideEnough(t *testing.T) {
	result := PadString("hello world", 5)
	assert.Equal(t, "hello world", result)
}

func TestPadString_Empty(t *testing.T) {
	result := PadString("", 5)
	assert.Equal(t, 5, StringWidth(result))
}

func TestTruncateString(t *testing.T) {
	result := TruncateString("hello world", 5)
	assert.Equal(t, "hello", result)
	assert.Equal(t, 5, StringWidth(result))
}

func TestTruncateString_WithWideChars(t *testing.T) {
	result := TruncateString("Привет мир", 8)
	assert.Equal(t, 8, StringWidth(result))
}

func TestAlignCenter_ASCII(t *testing.T) {
	result := AlignCenter("hi", 10)
	assert.Equal(t, "    hi    ", result)
	assert.Equal(t, 10, StringWidth(result))
}

func TestAlignCenter_WithWideChars(t *testing.T) {
	result := AlignCenter("⽕", 6)
	assert.Equal(t, 6, StringWidth(result))
}

func TestAlignRight_ASCII(t *testing.T) {
	result := AlignRight("hi", 10)
	assert.Equal(t, "        hi", result)
	assert.Equal(t, 10, StringWidth(result))
}

func TestAlignRight_WithWideChars(t *testing.T) {
	result := AlignRight("Привет", 12)
	assert.Equal(t, 12, StringWidth(result))
}

func TestContainsWideChars_True(t *testing.T) {
	assert.True(t, ContainsWideChars("⽕ U "))
	assert.True(t, ContainsWideChars("hello⽕"))
}

func TestContainsWideChars_False(t *testing.T) {
	assert.False(t, ContainsWideChars("hello"))
	assert.False(t, ContainsWideChars("12345"))
	assert.False(t, ContainsWideChars("abc"))
}

func TestValidUTF8(t *testing.T) {
	assert.True(t, ValidUTF8("hello"))
	assert.True(t, ValidUTF8("Привет"))
	assert.True(t, ValidUTF8("⽕ U "))
	assert.False(t, ValidUTF8(string([]byte{0xff, 0xfe})))
}

func TestPadString_MatchesEnglishAndRussian(t *testing.T) {
	en := "Player 1 ♥ ♥ ♥"
	ru := "Игрок 1 ♥ ♥ ♥"

	paddedEn := PadString(en, 30)
	paddedRu := PadString(ru, 30)

	assert.Equal(t, 30, StringWidth(paddedEn))
	assert.Equal(t, 30, StringWidth(paddedRu))
}
