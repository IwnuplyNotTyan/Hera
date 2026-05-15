//go:build notint

package generate

import "github.com/charmbracelet/lipgloss"

type noTintRegistry struct{}

func newNoTintRegistry() ThemeRegistry {
	return &noTintRegistry{}
}

func (r *noTintRegistry) SetTintID(id string) bool             { return true }
func (r *noTintRegistry) TintIDs() []string                    { return []string{"default"} }
func (r *noTintRegistry) Fg() lipgloss.TerminalColor           { return lipgloss.Color("#AAAAAA") }
func (r *noTintRegistry) Bg() lipgloss.TerminalColor           { return lipgloss.Color("#171717") }
func (r *noTintRegistry) SelectionBg() lipgloss.TerminalColor  { return lipgloss.Color("#444444") }
func (r *noTintRegistry) Cursor() lipgloss.TerminalColor       { return lipgloss.Color("#333333") }
func (r *noTintRegistry) BrightBlack() lipgloss.TerminalColor  { return lipgloss.Color("#555555") }
func (r *noTintRegistry) BrightBlue() lipgloss.TerminalColor   { return lipgloss.Color("#75B7FF") }
func (r *noTintRegistry) BrightCyan() lipgloss.TerminalColor   { return lipgloss.Color("#88AACC") }
func (r *noTintRegistry) BrightGreen() lipgloss.TerminalColor  { return lipgloss.Color("#75FFBA") }
func (r *noTintRegistry) BrightPurple() lipgloss.TerminalColor { return lipgloss.Color("#FF75B7") }
func (r *noTintRegistry) BrightRed() lipgloss.TerminalColor    { return lipgloss.Color("#FF4444") }
func (r *noTintRegistry) BrightWhite() lipgloss.TerminalColor  { return lipgloss.Color("#FFFFFF") }
func (r *noTintRegistry) BrightYellow() lipgloss.TerminalColor { return lipgloss.Color("#FFD700") }
func (r *noTintRegistry) Black() lipgloss.TerminalColor        { return lipgloss.Color("#000000") }
func (r *noTintRegistry) Blue() lipgloss.TerminalColor         { return lipgloss.Color("#146fba") }
func (r *noTintRegistry) Cyan() lipgloss.TerminalColor         { return lipgloss.Color("#5DA9E9") }
func (r *noTintRegistry) Green() lipgloss.TerminalColor        { return lipgloss.Color("#a2ffd0") }
func (r *noTintRegistry) Purple() lipgloss.TerminalColor       { return lipgloss.Color("#874BFD") }
func (r *noTintRegistry) Red() lipgloss.TerminalColor          { return lipgloss.Color("#FF4400") }
func (r *noTintRegistry) White() lipgloss.TerminalColor        { return lipgloss.Color("#FFFFFF") }
func (r *noTintRegistry) Yellow() lipgloss.TerminalColor       { return lipgloss.Color("#ffdb1e") }
