package tests

import (
	"testing"

	generate "hera/core"

	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
	"github.com/stretchr/testify/assert"
)

func TestDefaultTheme_ImplementsTint(t *testing.T) {
	theme := &generate.DefaultTheme{}
	var _ tint.Tint = theme
}

func TestDefaultTheme_ID(t *testing.T) {
	theme := &generate.DefaultTheme{}
	assert.Equal(t, "default", theme.ID())
}

func TestDefaultTheme_DisplayName(t *testing.T) {
	theme := &generate.DefaultTheme{}
	assert.Equal(t, "Default", theme.DisplayName())
}

func TestDefaultTheme_About(t *testing.T) {
	theme := &generate.DefaultTheme{}
	assert.Equal(t, "Original Hera colors", theme.About())
}

func TestDefaultTheme_BasicColors(t *testing.T) {
	theme := &generate.DefaultTheme{}

	assert.Equal(t, lipgloss.Color("#AAAAAA"), theme.Fg())
	assert.Equal(t, lipgloss.Color("#171717"), theme.Bg())
	assert.Equal(t, lipgloss.Color("#444444"), theme.SelectionBg())
	assert.Equal(t, lipgloss.Color("#222222"), theme.Cursor())
}

func TestDefaultTheme_BrightColors(t *testing.T) {
	theme := &generate.DefaultTheme{}

	assert.Equal(t, lipgloss.Color("#555555"), theme.BrightBlack())
	assert.Equal(t, lipgloss.Color("#75B7FF"), theme.BrightBlue())
	assert.Equal(t, lipgloss.Color("#88AACC"), theme.BrightCyan())
	assert.Equal(t, lipgloss.Color("#75FFBA"), theme.BrightGreen())
	assert.Equal(t, lipgloss.Color("#FF75B7"), theme.BrightPurple())
	assert.Equal(t, lipgloss.Color("#FF4444"), theme.BrightRed())
	assert.Equal(t, lipgloss.Color("#FFFFFF"), theme.BrightWhite())
	assert.Equal(t, lipgloss.Color("#FFD700"), theme.BrightYellow())
}

func TestDefaultTheme_StandardColors(t *testing.T) {
	theme := &generate.DefaultTheme{}

	assert.Equal(t, lipgloss.Color("#000000"), theme.Black())
	assert.Equal(t, lipgloss.Color("#146fba"), theme.Blue())
	assert.Equal(t, lipgloss.Color("#5DA9E9"), theme.Cyan())
	assert.Equal(t, lipgloss.Color("#a2ffd0"), theme.Green())
	assert.Equal(t, lipgloss.Color("#874BFD"), theme.Purple())
	assert.Equal(t, lipgloss.Color("#FF4400"), theme.Red())
	assert.Equal(t, lipgloss.Color("#FFFFFF"), theme.White())
	assert.Equal(t, lipgloss.Color("#ffdb1e"), theme.Yellow())
}

func TestRegisterDefaults(t *testing.T) {
	tint.NewDefaultRegistry()
	generate.RegisterDefaults()

	ok := tint.DefaultRegistry.SetTintID("default")
	assert.True(t, ok, "default theme should be registered")
}
