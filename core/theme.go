package generate

import "github.com/charmbracelet/lipgloss"

type ThemeRegistry interface {
	SetTintID(id string) bool
	TintIDs() []string
	Fg() lipgloss.TerminalColor
	Bg() lipgloss.TerminalColor
	SelectionBg() lipgloss.TerminalColor
	Cursor() lipgloss.TerminalColor
	BrightBlack() lipgloss.TerminalColor
	BrightBlue() lipgloss.TerminalColor
	BrightCyan() lipgloss.TerminalColor
	BrightGreen() lipgloss.TerminalColor
	BrightPurple() lipgloss.TerminalColor
	BrightRed() lipgloss.TerminalColor
	BrightWhite() lipgloss.TerminalColor
	BrightYellow() lipgloss.TerminalColor
	Black() lipgloss.TerminalColor
	Blue() lipgloss.TerminalColor
	Cyan() lipgloss.TerminalColor
	Green() lipgloss.TerminalColor
	Purple() lipgloss.TerminalColor
	Red() lipgloss.TerminalColor
	White() lipgloss.TerminalColor
	Yellow() lipgloss.TerminalColor
}
