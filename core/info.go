package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"hera/i18n"
)

var effectSep = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render(" · ")

func renderEffects(effects []Effect, loc i18n.Localizer) string {
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
				Render(loc.T("effects.fire", e.Duration))
		case EffectWet:
			s = lipgloss.NewStyle().Foreground(lipgloss.Color("#146fba")).Bold(true).
				Render(loc.T("effects.wet", e.Duration))
		case EffectSmoke:
			s = lipgloss.NewStyle().Foreground(lipgloss.Color("#88AACC")).Bold(true).
				Render(loc.T("effects.smoke", e.Duration))
		default:
			s = icon + " " + fmt.Sprint(e.Duration)
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, effectSep)
}

func (m Model) cursorInfo() string {
	if len(m.Players) == 0 {
		return ""
	}
	loc := m.Localizer
	p := Point{m.CursorX, m.CursorY}
	current := m.Players[m.CurrentPlayer]
	wallBlocked := !m.UltMode && m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

	for i, pl := range m.Players {
		if pl.X == m.CursorX && pl.Y == m.CursorY {
			hp := strings.Repeat("♥ ", pl.HP) + strings.Repeat("♡ ", MaxHP-pl.HP)
			effectStr := renderEffects(pl.Effects, loc)

			if i == m.CurrentPlayer {
				result := pl.Style.Render(loc.T("cursor.player.you", hp))
				if effectStr != "" {
					result += "\n" + effectStr
				}
				return result
			}
			if wallBlocked {
				result := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(loc.T("cursor.player.wallBlocked", i+1, hp))
				if effectStr != "" {
					result += "\n" + effectStr
				}
				return result
			}
			result := pl.Style.Render(loc.T("cursor.player.other", i+1, hp))
			if effectStr != "" {
				result += "\n" + effectStr
			}
			return result
		}
	}

	for i, en := range m.Enemys {
		if en.X == m.CursorX && en.Y == m.CursorY {
			hp := strings.Repeat("♥ ", en.HP) + strings.Repeat("♡ ", MaxHP-en.HP)
			effectStr := renderEffects(en.Effects, loc)

			if wallBlocked {
				result := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(loc.T("cursor.enemy.wallBlocked", i+1, hp))
				if effectStr != "" {
					result += "\n" + effectStr
				}
				return result
			}
			result := en.Style.Render(loc.T("cursor.enemy.default", i+1, hp))
			if effectStr != "" {
				result += "\n" + effectStr
			}
			return result
		}
	}

	switch {
	case m.Walls[p]:
		return m.Styles.WallStyle.Render(loc.T("cursor.tile.wall"))
	case wallBlocked:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
			Render(loc.T("cursor.tile.wallInWay"))
	case m.SmokeTiles[p] > 0:
		return m.Styles.SteamStyle.Render(loc.T("cursor.tile.smoke", m.SmokeTiles[p]))
	case m.Water[p]:
		return m.Styles.WaterStyle.Render(loc.T("cursor.tile.water"))
	case m.FireTiles[p] > 0:
		return m.Styles.FireStyle.Render(loc.T("cursor.tile.fire", m.FireTiles[p]))
	case m.UltMode:
		if m.ultInAxisRange(m.CursorX, m.CursorY) {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4400")).Render(loc.T("cursor.range.ult"))
		}
		return m.Styles.CellStyle.Render(loc.T("cursor.range.outOfUltAxis"))
	case m.IsInRange(m.CursorX, m.CursorY):
		if m.ShootMode {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Render(loc.T("cursor.range.inShootRange"))
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render(loc.T("cursor.range.inMoveRange"))
	default:
		return m.Styles.CellStyle.Render(loc.T("cursor.range.empty"))
	}
}
