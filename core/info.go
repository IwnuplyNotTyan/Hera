package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var effectSep = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render(" · ")

func renderEffects(effects []Effect) string {
	if len(effects) == 0 {
		return ""
	}
	var parts []string
	for _, e := range effects {
		icon := effectIcon(e.Type)
		var s string
		switch e.Type {
		case EffectFire:
			s = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4400")).Bold(true).
				Render(fmt.Sprintf("%s %d", icon, e.Duration))
		case EffectWet:
			s = lipgloss.NewStyle().Foreground(lipgloss.Color("#146fba")).Bold(true).
				Render(fmt.Sprintf("%s %d", icon, e.Duration))
		case EffectSmoke:
			s = lipgloss.NewStyle().Foreground(lipgloss.Color("#88AACC")).Bold(true).
				Render(fmt.Sprintf("%s %d", icon, e.Duration))
		default:
			s = fmt.Sprintf("%s %d", icon, e.Duration)
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, effectSep)
}

func (m Model) cursorInfo() string {
	if len(m.Players) == 0 {
		return ""
	}
	p := Point{m.CursorX, m.CursorY}
	current := m.Players[m.CurrentPlayer]
	// wallBlocked — только для shoot/move режимов, не для ult
	wallBlocked := !m.UltMode && m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

	for i, pl := range m.Players {
		if pl.X == m.CursorX && pl.Y == m.CursorY {
			hp := strings.Repeat("♥ ", pl.HP) + strings.Repeat("♡ ", MaxHP-pl.HP)
			effectStr := renderEffects(pl.Effects)

			if i == m.CurrentPlayer {
				result := pl.Style.Render(fmt.Sprintf("● Player %s (you)", hp))
				if effectStr != "" {
					result += "\n" + effectStr
				}
				return result
			}
			if wallBlocked {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(fmt.Sprintf("■ Player %d — wall in the way", i+1))
			}
			result := pl.Style.Render(fmt.Sprintf("■ Player %d %s", i+1, hp))
			if effectStr != "" {
				result += "\n" + effectStr
			}
			return result
		}
	}

	for i, en := range m.Enemys {
		if en.X == m.CursorX && en.Y == m.CursorY {
			hp := strings.Repeat("♥ ", en.HP) + strings.Repeat("♡ ", MaxHP-en.HP)
			effectStr := renderEffects(en.Effects)

			if wallBlocked {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(fmt.Sprintf("▲ Enemy %d — wall in the way", i+1))
			}
			result := en.Style.Render(fmt.Sprintf("▲ Enemy %d %s", i+1, hp))
			if effectStr != "" {
				result += "\n" + effectStr
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
	case m.SmokeTiles[p] > 0:
		return steamStyle.Render(fmt.Sprintf("~ Smoke — %d turns left", m.SmokeTiles[p]))
	case m.Water[p]:
		return waterStyle.Render("≈ Water — passable")
	case m.FireTiles[p] > 0:
		return fireStyle.Render(fmt.Sprintf("⽕ Fire — %d turns left", m.FireTiles[p]))
	case m.UltMode:
		if m.ultInAxisRange(m.CursorX, m.CursorY) {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4400")).Render("⽕ Ult target")
		}
		return cellStyle.Render("· Out of ult axis")
	case m.IsInRange(m.CursorX, m.CursorY):
		if m.ShootMode {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Render("· In shoot range")
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render("· In move range")
	default:
		return cellStyle.Render("· Empty — out of range")
	}
}
