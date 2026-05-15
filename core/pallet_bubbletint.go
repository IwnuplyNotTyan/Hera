//go:build !notint

package generate

import (
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

type DefaultTheme struct{}

func (t *DefaultTheme) DisplayName() string                  { return "Default" }
func (t *DefaultTheme) ID() string                           { return "default" }
func (t *DefaultTheme) About() string                        { return "Original Hera colors" }
func (t *DefaultTheme) Fg() lipgloss.TerminalColor           { return lipgloss.Color("#AAAAAA") }
func (t *DefaultTheme) Bg() lipgloss.TerminalColor           { return lipgloss.Color("#171717") }
func (t *DefaultTheme) SelectionBg() lipgloss.TerminalColor  { return lipgloss.Color("#444444") }
func (t *DefaultTheme) Cursor() lipgloss.TerminalColor       { return lipgloss.Color("#333333") }
func (t *DefaultTheme) BrightBlack() lipgloss.TerminalColor  { return lipgloss.Color("#282828") }
func (t *DefaultTheme) BrightBlue() lipgloss.TerminalColor   { return lipgloss.Color("#75B7FF") }
func (t *DefaultTheme) BrightCyan() lipgloss.TerminalColor   { return lipgloss.Color("#88AACC") }
func (t *DefaultTheme) BrightGreen() lipgloss.TerminalColor  { return lipgloss.Color("#75FFBA") }
func (t *DefaultTheme) BrightPurple() lipgloss.TerminalColor { return lipgloss.Color("#FF75B7") }
func (t *DefaultTheme) BrightRed() lipgloss.TerminalColor    { return lipgloss.Color("#FF4444") }
func (t *DefaultTheme) BrightWhite() lipgloss.TerminalColor  { return lipgloss.Color("#FFFFFF") }
func (t *DefaultTheme) BrightYellow() lipgloss.TerminalColor { return lipgloss.Color("#FFD700") }
func (t *DefaultTheme) Black() lipgloss.TerminalColor        { return lipgloss.Color("#000000") }
func (t *DefaultTheme) Blue() lipgloss.TerminalColor         { return lipgloss.Color("#146fba") }
func (t *DefaultTheme) Cyan() lipgloss.TerminalColor         { return lipgloss.Color("#5DA9E9") }
func (t *DefaultTheme) Green() lipgloss.TerminalColor        { return lipgloss.Color("#a2ffd0") }
func (t *DefaultTheme) Purple() lipgloss.TerminalColor       { return lipgloss.Color("#874BFD") }
func (t *DefaultTheme) Red() lipgloss.TerminalColor          { return lipgloss.Color("#FF4400") }
func (t *DefaultTheme) White() lipgloss.TerminalColor        { return lipgloss.Color("#FFFFFF") }
func (t *DefaultTheme) Yellow() lipgloss.TerminalColor       { return lipgloss.Color("#ffdb1e") }

func RegisterDefaults() {
	tint.Register(&DefaultTheme{})
}
