package generate

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2).
			Width(46)
	helpStyle = lipgloss.NewStyle().
			Padding(1, 2)
	cellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))
	wallStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#874BFD"))
	waterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#146fba"))
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#222222"))
	rangeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#171717"))
	fireStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4400")).
			Bold(true)
	steamStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#88AACC")).
			Bold(true)
	ultZoneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4400")).
			Background(lipgloss.Color("#1a0800"))
)

var playerStyles = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#75FFBA")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#75B7FF")).Bold(true),
}

var enemysStyles = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcbe4")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#a2ffd0")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#ffdb1e")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#84befe")).Bold(true),
}
