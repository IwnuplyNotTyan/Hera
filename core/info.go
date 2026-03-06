package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) cursorInfo() string {
	if len(m.Players) == 0 {
		return ""
	}
	p := Point{m.CursorX, m.CursorY}
	current := m.Players[m.CurrentPlayer]
	wallBlocked := m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

	for i, pl := range m.Players {
		if pl.X == m.CursorX && pl.Y == m.CursorY {
			hp := strings.Repeat("♥ ", pl.HP) + strings.Repeat("♡ ", MaxHP-pl.HP)

			effectStr := ""
			for _, e := range pl.Effects {
				effectStr += fmt.Sprintf(" %s %d ", effectIcon(e.Type), e.Duration)
			}

			if i == m.CurrentPlayer {
				result := pl.Style.Render(fmt.Sprintf("● Player %s (you)", hp))
				if effectStr != "" {
					result += "\n" + lipgloss.NewStyle().
						Foreground(lipgloss.Color("#146fba")).
						Render(effectStr)
				}
				return result
			}
			if wallBlocked {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(fmt.Sprintf("■ Player %d — wall in the way", i+1))
			}
			result := pl.Style.Render(fmt.Sprintf("■ Player %d %s", i+1, hp))
			if effectStr != "" {
				result += "\n" + lipgloss.NewStyle().
					Foreground(lipgloss.Color("#146fba")).
					Render(effectStr)
			}
			return result
		}
	}
	for i, en := range m.Enemys {
		if en.X == m.CursorX && en.Y == m.CursorY {
			hp := strings.Repeat("♥ ", en.HP) + strings.Repeat("♡ ", MaxHP-en.HP)

			effectStr := ""
			for _, e := range en.Effects {
				effectStr += fmt.Sprintf(" %s %d ", effectIcon(e.Type), e.Duration)
			}

			if wallBlocked {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(fmt.Sprintf("▲ Enemy %d — wall in the way", i+1))
			}
			result := en.Style.Render(fmt.Sprintf("▲ Enemy %d %s", i+1, hp))
			if effectStr != "" {
				result += "\n" + lipgloss.NewStyle().
					Foreground(lipgloss.Color("#146fba")).
					Render(effectStr)
			}
			return result
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
