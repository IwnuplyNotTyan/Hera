package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"hera/i18n"
)

var effectSep = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render(" · ")

func renderEffects(effects []Effect, loc i18n.Localizer, highlightIdx int, th ThemeRegistry) string {
	if th == nil {
		return ""
	}
	var parts []string
	for i, e := range effects {
		icon := EffectIcon(e.Type)
		var s string
		st := lipgloss.NewStyle().Bold(true)
		switch e.Type {
		case EffectFire:
			st = st.Foreground(th.Red())
			s = st.Render(loc.T("effects.fire", e.Duration))
		case EffectWet:
			st = st.Foreground(th.Blue())
			s = st.Render(loc.T("effects.wet", e.Duration))
		case EffectSmoke:
			st = st.Foreground(th.BrightCyan())
			s = st.Render(loc.T("effects.smoke", e.Duration))
		default:
			s = icon + " " + fmt.Sprint(e.Duration)
		}
		if i == highlightIdx {
			s = lipgloss.NewStyle().Background(lipgloss.Color("#3a3a3a")).Render(s)
		}
		parts = append(parts, s)
	}
	if len(parts) == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render("-")
	}
	return strings.Join(parts, effectSep)
}

func (m *Model) effectsAtCursor() []Effect {
	p := Point{m.CursorX, m.CursorY}
	for _, pl := range m.Players {
		if pl.X == p.X && pl.Y == p.Y && len(pl.Effects) > 0 {
			return pl.Effects
		}
	}
	for _, en := range m.Enemys {
		if en.X == p.X && en.Y == p.Y && len(en.Effects) > 0 {
			return en.Effects
		}
	}
	return nil
}

func (m *Model) effectDescLine(e Effect) string {
	var icon string
	var desc string
	switch e.Type {
	case EffectFire:
		icon = "⚹"
		desc = m.Localizer.T("effects.desc.fire")
	case EffectWet:
		icon = "≈"
		desc = m.Localizer.T("effects.desc.wet")
	case EffectSmoke:
		icon = "~"
		desc = m.Localizer.T("effects.desc.smoke")
	}
	return m.Styles.BoxStyle.Padding(0, 1).Width(40).Render(icon + " " + desc)
}

func (m *Model) cursorInfo() string {
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

			if i == m.CurrentPlayer {
				return pl.Style.Render(loc.T("cursor.player.you", hp))
			}
			if wallBlocked {
				return lipgloss.NewStyle().Foreground(m.Theme.BrightRed()).
					Render(loc.T("cursor.player.wallBlocked", i+1, hp))
			}
			return pl.Style.Render(loc.T("cursor.player.other", i+1, hp))
		}
	}

	for i, en := range m.Enemys {
		if en.X == m.CursorX && en.Y == m.CursorY {
			hp := strings.Repeat("♥ ", en.HP) + strings.Repeat("♡ ", MaxHP-en.HP)

			if wallBlocked {
				return lipgloss.NewStyle().Foreground(m.Theme.BrightRed()).
					Render(loc.T("cursor.enemy.wallBlocked", i+1, hp))
			}
			return en.Style.Render(loc.T("cursor.enemy.default", i+1, hp))
		}
	}

	switch {
	case m.Walls[p]:
		return m.Styles.WallStyle.Render(loc.T("cursor.tile.wall"))
	case wallBlocked:
		return m.Styles.BlockedWallStyle.Render(loc.T("cursor.tile.wallInWay"))
	case m.SmokeTiles[p] > 0:
		return m.Styles.SteamStyle.Render(loc.T("cursor.tile.smoke", m.SmokeTiles[p]))
	case m.Water[p]:
		return m.Styles.WaterStyle.Render(loc.T("cursor.tile.water"))
	case m.FireTiles[p] > 0:
		return m.Styles.FireStyle.Render(loc.T("cursor.tile.fire", m.FireTiles[p]))
	case m.UltMode:
		if m.ultInAxisRange(m.CursorX, m.CursorY) {
			return m.Styles.UltRangeStyle.Render(loc.T("cursor.range.ult"))
		}
		return m.Styles.CellStyle.Render(loc.T("cursor.range.outOfUltAxis"))
	case m.IsInRange(m.CursorX, m.CursorY):
		if m.ShootMode {
			return m.Styles.ShootRangeStyle.Render(loc.T("cursor.range.inShootRange"))
		}
		return m.Styles.MoveRangeStyle.Render(loc.T("cursor.range.inMoveRange"))
	default:
		return m.Styles.CellStyle.Render(loc.T("cursor.range.empty"))
	}
}

func (m *Model) effectsLine() string {
	if m.ShowEffectIdx == 0 {
		return ""
	}

	effects := m.effectsAtCursor()
	skipIdx := len(effects) + 1

	loc := m.Localizer
	highlightIdx := -1
	if m.ShowEffectIdx > 0 && m.ShowEffectIdx <= len(effects) {
		highlightIdx = m.ShowEffectIdx - 1
	}

	var result string
	if len(effects) > 0 {
		result = renderEffects(effects, loc, highlightIdx, m.Theme)
	}

	var skipIcon string
	if m.ShowEffectIdx == skipIdx {
		skipIcon = lipgloss.NewStyle().Background(lipgloss.Color("#3a3a3a")).Render("⏭")
	} else {
		skipIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render("⏭")
	}
	if result == "" {
		result = skipIcon
	} else {
		result += effectSep + skipIcon
	}

	if m.DebugMode {
		consoleIcon := lipgloss.NewStyle().
			Foreground(m.Theme.BrightBlack()).
			Render("»")
		result += effectSep + consoleIcon
	}

	if m.ShowEffectIdx > 0 && m.ShowEffectIdx <= len(effects) {
		e := effects[m.ShowEffectIdx-1]
		result += "\n" + m.effectDescLine(e)
	} else if m.ShowEffectIdx == skipIdx {
		result += "\n" + m.Styles.BoxStyle.Padding(0, 1).Width(40).
			Render(m.Localizer.T("help.skipTurn"))
	}
	return result
}
