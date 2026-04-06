package generate

import (
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

type Styles struct {
	BoxStyle        lipgloss.Style
	HelpStyle       lipgloss.Style
	CellStyle       lipgloss.Style
	WallStyle       lipgloss.Style
	WaterStyle      lipgloss.Style
	WaterRangeStyle lipgloss.Style
	CursorStyle     lipgloss.Style
	RangeStyle      lipgloss.Style
	FireStyle       lipgloss.Style
	SteamStyle      lipgloss.Style
	UltZoneStyle    lipgloss.Style
	UltAxisStyle    lipgloss.Style
	ShootRangeStyle lipgloss.Style
	PlayerStyles    []lipgloss.Style
	EnemysStyles    []lipgloss.Style
}

func NewStyles(th *tint.Registry) Styles {
	return Styles{
		BoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(th.Purple()).
			Padding(1, 2).
			Width(46),
		HelpStyle: lipgloss.NewStyle().
			Padding(1, 2),
		CellStyle: lipgloss.NewStyle().
			Foreground(th.Fg()),
		WallStyle: lipgloss.NewStyle().
			Foreground(th.Purple()),
		WaterStyle: lipgloss.NewStyle().
			Foreground(th.Blue()),
		WaterRangeStyle: lipgloss.NewStyle().
			Foreground(th.Cyan()).
			Background(th.Bg()),
		CursorStyle: lipgloss.NewStyle().
			Foreground(th.BrightWhite()).
			Background(th.Cursor()),
		RangeStyle: lipgloss.NewStyle().
			Background(th.Bg()),
		FireStyle: lipgloss.NewStyle().
			Foreground(th.Red()).
			Bold(true),
		SteamStyle: lipgloss.NewStyle().
			Foreground(th.BrightCyan()).
			Bold(true),
		UltZoneStyle: lipgloss.NewStyle().
			Foreground(th.Red()).
			Background(lipgloss.Color("#1a0800")),
		UltAxisStyle: lipgloss.NewStyle().
			Foreground(th.Yellow()).
			Background(lipgloss.Color("#0d0500")),
		ShootRangeStyle: lipgloss.NewStyle().
			Foreground(th.BrightRed()).
			Background(lipgloss.Color("#1a0505")),
		PlayerStyles: []lipgloss.Style{
			lipgloss.NewStyle().Foreground(th.BrightPurple()).Bold(true),
			lipgloss.NewStyle().Foreground(th.BrightGreen()).Bold(true),
			lipgloss.NewStyle().Foreground(th.Yellow()).Bold(true),
			lipgloss.NewStyle().Foreground(th.BrightBlue()).Bold(true),
		},
		EnemysStyles: []lipgloss.Style{
			lipgloss.NewStyle().Foreground(th.Purple()).Bold(true),
			lipgloss.NewStyle().Foreground(th.Green()).Bold(true),
			lipgloss.NewStyle().Foreground(th.Yellow()).Bold(true),
			lipgloss.NewStyle().Foreground(th.Blue()).Bold(true),
		},
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
