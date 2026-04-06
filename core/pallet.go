package generate

import (
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

func initStyles(th *tint.Registry) {
	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.Purple()).
		Padding(1, 2).
		Width(46)
	helpStyle = lipgloss.NewStyle().
		Padding(1, 2)
	cellStyle = lipgloss.NewStyle().
		Foreground(th.Fg())
	wallStyle = lipgloss.NewStyle().
		Foreground(th.Purple())
	waterStyle = lipgloss.NewStyle().
		Foreground(th.Blue())
	waterRangeStyle = lipgloss.NewStyle().
		Foreground(th.Cyan()).
		Background(th.Bg())
	cursorStyle = lipgloss.NewStyle().
		Foreground(th.BrightWhite()).
		Background(th.Cursor())
	rangeStyle = lipgloss.NewStyle().
		Background(th.Bg())
	fireStyle = lipgloss.NewStyle().
		Foreground(th.Red()).
		Bold(true)
	steamStyle = lipgloss.NewStyle().
		Foreground(th.BrightCyan()).
		Bold(true)
	ultZoneStyle = lipgloss.NewStyle().
		Foreground(th.Red()).
		Background(lipgloss.Color("#1a0800"))
	ultAxisStyle = lipgloss.NewStyle().
		Foreground(th.Yellow()).
		Background(lipgloss.Color("#0d0500"))
	shootRangeStyle = lipgloss.NewStyle().
		Foreground(th.BrightRed()).
		Background(lipgloss.Color("#1a0505"))

	playerStyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(th.BrightPurple()).Bold(true),
		lipgloss.NewStyle().Foreground(th.BrightGreen()).Bold(true),
		lipgloss.NewStyle().Foreground(th.Yellow()).Bold(true),
		lipgloss.NewStyle().Foreground(th.BrightBlue()).Bold(true),
	}

	enemysStyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(th.Purple()).Bold(true),
		lipgloss.NewStyle().Foreground(th.Green()).Bold(true),
		lipgloss.NewStyle().Foreground(th.Yellow()).Bold(true),
		lipgloss.NewStyle().Foreground(th.Blue()).Bold(true),
	}
}

func RegisterDefaults() {
	tint.Register(&DefaultTheme{})
}

type DefaultTheme struct{}

func (t *DefaultTheme) DisplayName() string                  { return "Default" }
func (t *DefaultTheme) ID() string                           { return "default" }
func (t *DefaultTheme) About() string                        { return "Original Hera colors" }
func (t *DefaultTheme) Fg() lipgloss.TerminalColor           { return lipgloss.Color("#AAAAAA") }
func (t *DefaultTheme) Bg() lipgloss.TerminalColor           { return lipgloss.Color("#171717") }
func (t *DefaultTheme) SelectionBg() lipgloss.TerminalColor  { return lipgloss.Color("#444444") }
func (t *DefaultTheme) Cursor() lipgloss.TerminalColor       { return lipgloss.Color("#222222") }
func (t *DefaultTheme) BrightBlack() lipgloss.TerminalColor  { return lipgloss.Color("#555555") }
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

var (
	boxStyle        lipgloss.Style
	helpStyle       lipgloss.Style
	cellStyle       lipgloss.Style
	wallStyle       lipgloss.Style
	waterStyle      lipgloss.Style
	waterRangeStyle lipgloss.Style
	cursorStyle     lipgloss.Style
	rangeStyle      lipgloss.Style
	fireStyle       lipgloss.Style
	steamStyle      lipgloss.Style
	ultZoneStyle    lipgloss.Style
	ultAxisStyle    lipgloss.Style
	shootRangeStyle lipgloss.Style
)

var (
	playerStyles []lipgloss.Style
	enemysStyles []lipgloss.Style
)
