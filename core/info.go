package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) cursorInfo() string {
	p := Point{m.CursorX, m.CursorY}
	current := m.Players[m.CurrentPlayer]
	wallBlocked := m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

	for i, pl := range m.Players {
		if pl.X == m.CursorX && pl.Y == m.CursorY {
			hp := strings.Repeat("♥ ", pl.HP) + strings.Repeat("♡ ", MaxHP-pl.HP)
			if i == m.CurrentPlayer {
				return pl.Style.Render(fmt.Sprintf("● Player %s (you)", hp))
			}
			if wallBlocked {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(fmt.Sprintf("■ Player %d — wall in the way", i+1))
			}
			return pl.Style.Render(fmt.Sprintf("■ Player %s", hp))
		}
	}
	for i, en := range m.Enemys {
		if en.X == m.CursorX && en.Y == m.CursorY {
			hp := strings.Repeat("♥ ", en.HP) + strings.Repeat("♡ ", MaxHP-en.HP)
			if wallBlocked {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(fmt.Sprintf("▲ Enemy %d — wall in the way", i+1))
			}
			return en.Style.Render(fmt.Sprintf("▲ Enemy %d %s", i+1, hp))
		}
	}

	switch {
	case m.Walls[p]:
		return wallStyle.Render("■ Wall — impassable")
	case wallBlocked:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
			Render("⊘ Wall in the way")
	case m.Water[p]:
		return waterStyle.Render("≈ Water — passable")
	case m.IsInRange(m.CursorX, m.CursorY):
		if m.ShootMode {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Render("· In shoot range")
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render("· In move range")
	default:
		return cellStyle.Render("· Empty — out of range")
	}
}
