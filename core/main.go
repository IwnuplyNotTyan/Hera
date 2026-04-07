package generate

import (
	"fmt"
	"strings"

	"hera/utils"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.Moved && !m.Shot {
	} else if m.Moved {
		m.ShootMode = true
	} else {
		m.ShootMode = false
	}

	switch msg := msg.(type) {
	case enemyTurnMsg:
		if len(m.Players) == 0 {
			return m, tea.Quit
		}
		if msg.enemyIdx >= len(m.Enemys) {
			m.EnemyTurn = false
			m.EnemyIdx = 0
			if len(m.Players) > 0 {
				cur := m.Players[m.CurrentPlayer]
				m.CursorX = cur.X
				m.CursorY = cur.Y
				m.UltMode = false
				m.UltAxis = ""
			}
			return m, nil
		}
		m.EnemyIdx = msg.enemyIdx
		m = m.doEnemyTurn(msg.enemyIdx)
		if len(m.Players) == 0 {
			return m, tea.Quit
		}
		return m, enemyTurnCmd(msg.enemyIdx + 1)

	case tea.MouseMsg:
		if m.EnemyTurn {
			return m, nil
		}
		if msg.Button != tea.MouseButtonLeft || msg.Action != tea.MouseActionPress {
			return m, nil
		}
		for col := 0; col < GridW; col++ {
			for row := 0; row < GridH; row++ {
				if m.Z.Get(fmt.Sprintf("cell-%d-%d", col, row)).InBounds(msg) {
					m.CursorX = col
					m.CursorY = row
					break
				}
			}
		}

	case tea.KeyMsg:
		if m.EnemyTurn {
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Up):
			newY := utils.Clamp(m.CursorY-1, 0, GridH-1)
			if m.UltMode {
				cur := m.Players[m.CurrentPlayer]
				if m.CursorX == cur.X && m.CursorY == cur.Y {
					m.UltAxis = ""
				}
				if m.UltAxis == "" || m.UltAxis == "v" {
					m.UltAxis = "v"
					m.CursorY = newY
				}
			} else {
				m.CursorY = newY
			}
		case key.Matches(msg, m.keys.Down):
			newY := utils.Clamp(m.CursorY+1, 0, GridH-1)
			if m.UltMode {
				cur := m.Players[m.CurrentPlayer]
				if m.CursorX == cur.X && m.CursorY == cur.Y {
					m.UltAxis = ""
				}
				if m.UltAxis == "" || m.UltAxis == "v" {
					m.UltAxis = "v"
					m.CursorY = newY
				}
			} else {
				m.CursorY = newY
			}
		case key.Matches(msg, m.keys.Left):
			newX := utils.Clamp(m.CursorX-1, 0, GridW-1)
			if m.UltMode {
				cur := m.Players[m.CurrentPlayer]
				if m.CursorX == cur.X && m.CursorY == cur.Y {
					m.UltAxis = ""
				}
				if m.UltAxis == "" || m.UltAxis == "h" {
					m.UltAxis = "h"
					m.CursorX = newX
				}
			} else {
				m.CursorX = newX
			}
		case key.Matches(msg, m.keys.Right):
			newX := utils.Clamp(m.CursorX+1, 0, GridW-1)
			if m.UltMode {
				cur := m.Players[m.CurrentPlayer]
				if m.CursorX == cur.X && m.CursorY == cur.Y {
					m.UltAxis = ""
				}
				if m.UltAxis == "" || m.UltAxis == "h" {
					m.UltAxis = "h"
					m.CursorX = newX
				}
			} else {
				m.CursorX = newX
			}

		case key.Matches(msg, m.keys.Ult):
			cur := m.Players[m.CurrentPlayer]
			m.CursorX = cur.X
			m.CursorY = cur.Y
			if !m.Shot && m.Players[m.CurrentPlayer].UltCharges > 0 {
				m.UltMode = !m.UltMode
				m.UltAxis = ""
				m.ShootMode = false
				if m.UltMode {
					cur := m.Players[m.CurrentPlayer]
					m.CursorX = cur.X
					m.CursorY = cur.Y
				}
			}

		case key.Matches(msg, m.keys.Shoot):
			if !m.Shot {
				m.ShootMode = !m.ShootMode
				m.UltMode = false
				cur := m.Players[m.CurrentPlayer]
				m.CursorX = cur.X
				m.CursorY = cur.Y
			}

		case key.Matches(msg, m.keys.Confirm):
			p := Point{m.CursorX, m.CursorY}
			current := m.Players[m.CurrentPlayer]
			wallBlocked := m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

			if m.UltMode && !m.Shot {
				m = m.doUlt()
				cur := m.Players[m.CurrentPlayer]
				m.CursorX = cur.X
				m.CursorY = cur.Y
			} else if m.ShootMode && !m.Shot {
				if hasEffect(m.Players[m.CurrentPlayer].Effects, EffectSmoke) {
					m.Shot = true
					m.ShootMode = false
					break
				}
				if m.IsInRange(m.CursorX, m.CursorY) && !m.Walls[p] && !m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY) {
					for i, pl := range m.Players {
						if i != m.CurrentPlayer && pl.X == m.CursorX && pl.Y == m.CursorY {
							m.Players[i].HP--
							if m.Players[i].HP <= 0 {
								m.Players = append(m.Players[:i], m.Players[i+1:]...)
								if m.CurrentPlayer >= len(m.Players) {
									m.CurrentPlayer = 0
								}
							}
							break
						}
					}
					for i, en := range m.Enemys {
						if en.X == m.CursorX && en.Y == m.CursorY {
							m.Enemys[i].HP--
							if m.Enemys[i].HP <= 0 {
								m.Enemys = append(m.Enemys[:i], m.Enemys[i+1:]...)
							}
							break
						}
					}
					m.Shot = true
					m.ShootMode = false
					cur := m.Players[m.CurrentPlayer]
					m.CursorX = cur.X
					m.CursorY = cur.Y
				}
			} else if !m.ShootMode && !m.UltMode && !m.Moved {
				if m.IsInRange(m.CursorX, m.CursorY) && !m.Walls[p] && !wallBlocked && !m.OccupiedByOther(m.CursorX, m.CursorY) {
					m.Players[m.CurrentPlayer].X = m.CursorX
					m.Players[m.CurrentPlayer].Y = m.CursorY

					if m.Water[p] {
						m.Players[m.CurrentPlayer].Effects = resolveEffects(
							m.Players[m.CurrentPlayer].Effects,
							Effect{Type: EffectWet, Duration: 2},
						)
					}
					if m.FireTiles[p] > 0 {
						if !hasEffect(m.Players[m.CurrentPlayer].Effects, EffectWet) {
							m.Players[m.CurrentPlayer].Effects = resolveEffects(
								m.Players[m.CurrentPlayer].Effects,
								Effect{Type: EffectFire, Duration: 2},
							)
						}
					}
					if m.SmokeTiles[p] > 0 {
						m.Players[m.CurrentPlayer].Effects = resolveEffects(
							m.Players[m.CurrentPlayer].Effects,
							Effect{Type: EffectSmoke, Duration: 2},
						)
					}

					m.Moved = true
					m.CursorX = m.Players[m.CurrentPlayer].X
					m.CursorY = m.Players[m.CurrentPlayer].Y
				}
			}

			if m.Moved && m.Shot {
				m = m.nextTurn()
				if m.CurrentPlayer == 0 {
					m.EnemyTurn = true
					return m, enemyTurnCmd(0)
				}
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if len(m.Players) == 0 {
		gameOver := m.Styles.BoxStyle.Render(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF4444")).
				Bold(true).
				Render(m.Localizer.T("game.gameOver")),
		)
		return gameOver
	}

	current := m.Players[m.CurrentPlayer]
	hp := strings.Repeat("♥ ", current.HP) + strings.Repeat("♡ ", MaxHP-current.HP)
	hpStyle := current.Style
	if current.HP == 1 {
		hpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Blink(true)
	}
	hpStr := hpStyle.Render(m.Localizer.T("status.player", m.CurrentPlayer+1, hp))

	var modeStr string
	switch {
	case m.UltMode:
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4400")).
			Bold(true).
			Render(m.Localizer.T("status.ult"))
	case m.ShootMode:
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Bold(true).
			Render(m.Localizer.T("status.shoot"))
	default:
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Render(m.Localizer.T("status.move"))
	}

	ultCharges := m.Players[m.CurrentPlayer].UltCharges
	var ultStr string
	if ultCharges > 0 {
		ultStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4400")).
			Render(m.Localizer.T("status.ultCharges", ultCharges))
	} else {
		ultStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444")).
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
					st = st.Background(lipgloss.Color("#171717"))
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
					st = st.Background(lipgloss.Color("#171717"))
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
				cellContent = m.Styles.FireStyle.Render(" ⁺ ")
			case isUltCross:
				cellContent = m.Styles.UltZoneStyle.Render(" + ")
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
			cells = append(cells, m.Z.Mark(fmt.Sprintf("cell-%d-%d", col, row), cellContent))
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
		lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(
			m.Localizer.T("cursor.coordinates", map[string]interface{}{"x": m.CursorX, "y": m.CursorY}),
		),
		info,
	)

	status := m.Styles.BoxStyle.Render(line1 + "\n" + line2 + "\n" + line0)
	grid := strings.Join(rows, "\n")
	box := m.Styles.BoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, grid))
	helpView := m.Styles.HelpStyle.Render(m.help.View(m.keys))
	content := lipgloss.JoinVertical(lipgloss.Left, box, status, helpView)
	content = m.Z.Scan(content)
	return content
}
