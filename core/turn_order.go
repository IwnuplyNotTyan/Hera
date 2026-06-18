package generate

import (
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) turnOrder() string {
	var parts []string

	for i, pl := range m.Players {
		symbol := " ■ "
		style := pl.Style
		if i == m.CurrentPlayer && !m.EnemyTurn {
			style = style.Underline(true).Bold(true)
			symbol = " ● "
		}
		parts = append(parts, style.Render(symbol))
	}

	parts = append(parts, lipgloss.NewStyle().
		Foreground(m.Theme.SelectionBg()).Render(" · "))

	for i, en := range m.Enemys {
		symbol := " ▲ "
		style := en.Style
		if m.EnemyTurn && i == m.EnemyIdx {
			style = style.Underline(true).Bold(true)
			symbol = " ♦ "
		}
		parts = append(parts, style.Render(symbol))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}
