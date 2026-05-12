package generate

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	BoxStyle         lipgloss.Style
	HelpStyle        lipgloss.Style
	CellStyle        lipgloss.Style
	WallStyle        lipgloss.Style
	WaterStyle       lipgloss.Style
	WaterRangeStyle  lipgloss.Style
	CursorStyle      lipgloss.Style
	RangeStyle       lipgloss.Style
	FireStyle        lipgloss.Style
	SteamStyle       lipgloss.Style
	UltZoneStyle     lipgloss.Style
	UltAxisStyle     lipgloss.Style
	UltRangeStyle    lipgloss.Style
	ShootRangeStyle  lipgloss.Style
	MoveRangeStyle   lipgloss.Style
	BlockedWallStyle lipgloss.Style
	PlayerStyles     []lipgloss.Style
	EnemysStyles     []lipgloss.Style
}

func NewStyles(th ThemeRegistry) Styles {
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
			Background(th.BrightBlack()),
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
			Background(th.Bg()),
		UltAxisStyle: lipgloss.NewStyle().
			Foreground(th.Yellow()).
			Background(th.Bg()),
		ShootRangeStyle: lipgloss.NewStyle().
			Foreground(th.BrightRed()).
			Background(th.Bg()),
		UltRangeStyle: lipgloss.NewStyle().
			Foreground(th.Red()),
		MoveRangeStyle: lipgloss.NewStyle().
			Foreground(th.Fg()),
		BlockedWallStyle: lipgloss.NewStyle().
			Foreground(th.BrightRed()),
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
