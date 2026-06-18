package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"hera/utils"
)

func (m *Model) View() string {
	m.resetLayout()

	if m.Screen == ScreenProfiles {
		return m.viewProfiles()
	}
	if m.Screen == ScreenProfileCreate {
		return m.viewProfileCreate()
	}
	if m.Screen == ScreenMenu {
		return m.viewMenu()
	}
	if m.Screen == ScreenSettings {
		return m.viewSettings()
	}
	if m.Screen == ScreenThemeSelect {
		return m.viewThemeSelect()
	}
	if m.Screen == ScreenCenterSelect {
		return m.viewCenterSelect()
	}
	if m.Screen == ScreenGameOver {
		return m.viewGameOver()
	}
	if m.Screen == ScreenSeedPrompt {
		return m.viewSeedPrompt()
	}

	current := m.Players[m.CurrentPlayer]
	hp := strings.Repeat("♥ ", current.HP) + strings.Repeat("♡ ", MaxHP-current.HP)
	hpStyle := current.Style
	if current.HP == 1 {
		hpStyle = lipgloss.NewStyle().Foreground(m.Theme.Red()).Bold(true).Blink(true)
	}
	hpStr := hpStyle.Render(m.Localizer.T("status.player", m.CurrentPlayer+1, hp))

	var modeStr string
	switch {
	case m.UltMode:
		modeStr = lipgloss.NewStyle().
			Foreground(m.Theme.Red()).
			Bold(true).
			Render(m.Localizer.T("status.ult"))
	case m.ShootMode:
		modeStr = lipgloss.NewStyle().
			Foreground(m.Theme.BrightRed()).
			Bold(true).
			Render(m.Localizer.T("status.shoot"))
	default:
		modeStr = lipgloss.NewStyle().
			Foreground(m.Theme.Fg()).
			Render(m.Localizer.T("status.move"))
	}

	ultCharges := m.Players[m.CurrentPlayer].UltCharges
	var ultStr string
	if ultCharges > 0 {
		ultStr = lipgloss.NewStyle().
			Foreground(m.Theme.Red()).
			Render(m.Localizer.T("status.ultCharges", ultCharges))
	} else {
		ultStr = lipgloss.NewStyle().
			Foreground(m.Theme.SelectionBg()).
			Render(m.Localizer.T("status.ultChargesZero"))
	}

	var reachableZone map[Point]bool
	if !m.EnemyTurn && !m.UltMode && len(m.Players) > 0 {
		cur := m.Players[m.CurrentPlayer]
		r := m.currentRange()
		reachableZone = m.Reachable(cur.X, cur.Y, r)
	}

	ultAxisZone := make(map[Point]bool)
	ultCrossZone := make(map[Point]bool)
	if m.UltMode && len(m.Players) > 0 {
		cur := m.Players[m.CurrentPlayer]
		cx, cy := m.CursorX, m.CursorY
		for x := 0; x < GridW; x++ {
			if x != cur.X {
				ultAxisZone[Point{x, cur.Y}] = true
			}
		}
		for y := 0; y < GridH; y++ {
			if y != cur.Y {
				ultAxisZone[Point{cur.X, y}] = true
			}
		}
		for _, dp := range []Point{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			np := Point{cx + dp.X, cy + dp.Y}
			if np.X >= 0 && np.X < GridW && np.Y >= 0 && np.Y < GridH && !m.Walls[np] {
				ultCrossZone[np] = true
			}
		}
	}

	var rows []string
	for row := 0; row < GridH; row++ {
		var cells []string
		for col := 0; col < GridW; col++ {
			p := Point{col, row}
			playerIdx := -1
			enemyIdx := -1
			for i, pl := range m.Players {
				if pl.X == col && pl.Y == row {
					playerIdx = i
					break
				}
			}
			for i, en := range m.Enemys {
				if en.X == col && en.Y == row {
					enemyIdx = i
					break
				}
			}

			isCursor := col == m.CursorX && row == m.CursorY
			isUltCross := ultCrossZone[p]
			isUltAxis := ultAxisZone[p]
			isReachable := reachableZone[p]

			cellContent := ""
			switch {
			case isCursor:
				if playerIdx >= 0 {
					cellContent = m.Styles.CursorStyle.Render(m.Players[playerIdx].Style.Render(" ■ "))
				} else if enemyIdx >= 0 {
					cellContent = m.Styles.CursorStyle.Render(m.Enemys[enemyIdx].Style.Render(" ▲ "))
				} else {
					cellContent = m.Styles.CursorStyle.Render(" · ")
				}
			case playerIdx >= 0:
				symbol := " ■ "
				if playerIdx == m.CurrentPlayer {
					symbol = " ● "
				}
				st := m.Players[playerIdx].Style
				switch {
				case isUltCross:
					st = st.Background(lipgloss.Color("#2a0800"))
				case isUltAxis:
					st = st.Background(lipgloss.Color("#1a0a00"))
				case isReachable && m.ShootMode:
					st = st.Background(lipgloss.Color("#1a0505"))
				case isReachable:
					st = st.Background(m.Theme.Bg())
				}
				cellContent = st.Render(symbol)
			case enemyIdx >= 0:
				symbol := " ▲ "
				if enemyIdx == m.CurrentEnemy {
					symbol = " ♦ "
				}
				st := m.Enemys[enemyIdx].Style
				switch {
				case isUltCross:
					st = st.Background(lipgloss.Color("#2a0800"))
				case isUltAxis:
					st = st.Background(lipgloss.Color("#1a0a00"))
				case isReachable && m.ShootMode:
					st = st.Background(lipgloss.Color("#1a0505"))
				case isReachable:
					st = st.Background(m.Theme.Bg())
				}
				cellContent = st.Render(symbol)
			case m.Walls[p]:
				cellContent = m.Styles.WallStyle.Render(" ■ ")
			case m.SmokeTiles[p] > 0:
				cellContent = m.Styles.SteamStyle.Render(" ~ ")
			case m.Water[p]:
				switch {
				case isUltCross:
					cellContent = m.Styles.SteamStyle.Background(lipgloss.Color("#001a2a")).Render(" ~ ")
				case isUltAxis:
					cellContent = m.Styles.WaterStyle.Background(lipgloss.Color("#0d0800")).Render(" ≈ ")
				case m.IsInRange(col, row):
					cellContent = m.Styles.WaterRangeStyle.Render(" ≈ ")
				default:
					cellContent = m.Styles.WaterStyle.Render(" ≈ ")
				}
			case m.FireTiles[p] > 0:
				cellContent = m.Styles.FireStyle.Render(" ⚹ ")
			case isUltCross:
				cellContent = m.Styles.UltZoneStyle.Render(" ⚹ ")
			case isUltAxis:
				cellContent = m.Styles.UltAxisStyle.Render(" · ")
			case m.IsInRange(col, row):
				if m.ShootMode {
					cellContent = m.Styles.ShootRangeStyle.Render(" · ")
				} else if m.UltMode {
					cellContent = m.Styles.CellStyle.Render(" · ")
				} else {
					cellContent = m.Styles.RangeStyle.Render(" · ")
				}
			default:
				cellContent = m.Styles.CellStyle.Render(" · ")
			}
			cells = append(cells, cellContent)
			m.trackElement(Element{
				Type:   ElementGridCell,
				X:      col*cellWidth() + 2,
				Y:      row + 2,
				Width:  cellWidth(),
				Height: cellHeight(),
				ID:     fmt.Sprintf("cell-%d-%d", col, row),
				Index:  -1,
			})
		}
		rows = append(rows, strings.Join(cells, ""))
	}

	info := m.cursorInfo()
	info = utils.PadString(info, 40)
	line0 := m.turnOrder()

	line1 := lipgloss.JoinHorizontal(lipgloss.Top,
		modeStr,
		" ",
		hpStr,
		ultStr,
	)

	line2 := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Foreground(m.Theme.Cyan()).Render(
			m.Localizer.T("cursor.coordinates", map[string]interface{}{"x": m.CursorX, "y": m.CursorY}),
		),
		info,
	)

	line3 := m.effectsLine()

	var status string
	if m.ShowEffectIdx > 0 {
		status = m.boxStyle().Render(line3)
	} else {
		status = m.boxStyle().Render(line1 + "\n" + line2 + "\n" + line3 + "\n" + line0)
	}
	grid := strings.Join(rows, "\n")
	box := m.boxStyle().Render(lipgloss.JoinVertical(lipgloss.Left, grid))
	helpView := m.Styles.HelpStyle.Render(m.help.View(m.keys))

	var content string
	if m.ConsoleMode {
		consoleView := m.ConsoleView()
		content = lipgloss.JoinVertical(lipgloss.Left, box, consoleView, helpView)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left, box, status, helpView)
	}

	return m.renderContent(content)
}
