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
			} else if m.ShootMode || m.inRange(m.CursorX, newY) {
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
			} else if m.ShootMode || m.inRange(m.CursorX, newY) {
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
			} else if m.ShootMode || m.inRange(newX, m.CursorY) {
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
			} else if m.ShootMode || m.inRange(newX, m.CursorY) {
				m.CursorX = newX
			}

		case key.Matches(msg, m.keys.Ult):
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
				if hasEffect(m.Players[m.CurrentPlayer].Effects, EffectSteam) {
					m.Shot = true
					m.ShootMode = false
					break
				}
				if !m.Walls[p] && !m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY) {
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
				if !m.Walls[p] && !wallBlocked && !m.OccupiedByOther(m.CursorX, m.CursorY) {
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
					if m.SteamTiles[p] > 0 {
						m.Players[m.CurrentPlayer].Effects = resolveEffects(
							m.Players[m.CurrentPlayer].Effects,
							Effect{Type: EffectSteam, Duration: 2},
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
		return boxStyle.Render(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF4444")).
				Bold(true).
				Render("  ☠  Game Over  ☠  "),
		)
	}

	current := m.Players[m.CurrentPlayer]
	hp := strings.Repeat("♥ ", current.HP) + strings.Repeat("♡ ", MaxHP-current.HP)
	hpStyle := current.Style
	if current.HP == 1 {
		hpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Blink(true)
	}
	hpStr := hpStyle.Render(fmt.Sprintf("Player %d  %s", m.CurrentPlayer+1, hp))

	var modeStr string
	switch {
	case m.UltMode:
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4400")).
			Bold(true).
			Render("⽕ U ")
	case m.ShootMode:
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Bold(true).
			Render("♡ S ")
	default:
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Render("♧ M ")
	}

	ultCharges := m.Players[m.CurrentPlayer].UltCharges
	var ultStr string
	if ultCharges > 0 {
		ultStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4400")).
			Render(fmt.Sprintf(" ⽕×%d", ultCharges))
	} else {
		ultStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444")).
			Render(" ⽕×0")
	}

	var reachableZone map[Point]bool
	if !m.EnemyTurn && !m.UltMode && len(m.Players) > 0 {
		cur := m.Players[m.CurrentPlayer]
		r := m.currentRange()
		reachableZone = m.Reachable(cur.X, cur.Y, r)
	}

	// В UltMode показываем:
	//   ultAxisZone — вся ось (горизонталь или вертикаль) от игрока, зона прицела
	//   ultCrossZone — крест 5 клеток вокруг курсора, превью урона
	ultAxisZone := make(map[Point]bool)
	ultCrossZone := make(map[Point]bool)
	if m.UltMode && len(m.Players) > 0 {
		cur := m.Players[m.CurrentPlayer]
		cx, cy := m.CursorX, m.CursorY

		// ось прицела — по зафиксированной оси от игрока
		switch m.UltAxis {
		case "h":
			for x := 0; x < GridW; x++ {
				if x != cur.X {
					ultAxisZone[Point{x, cur.Y}] = true
				}
			}
		case "v":
			for y := 0; y < GridH; y++ {
				if y != cur.Y {
					ultAxisZone[Point{cur.X, y}] = true
				}
			}
		default:
			// ось не выбрана — показываем обе оси
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
		}

		// крест урона вокруг курсора (без стен)
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

			switch {
			case isCursor:
				if playerIdx >= 0 {
					cells = append(cells, cursorStyle.Render(m.Players[playerIdx].Style.Render(" ■ ")))
				} else if enemyIdx >= 0 {
					cells = append(cells, cursorStyle.Render(m.Enemys[enemyIdx].Style.Render(" ▲ ")))
				} else {
					cells = append(cells, cursorStyle.Render(" · "))
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
				cells = append(cells, st.Render(symbol))
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
				cells = append(cells, st.Render(symbol))
			case m.Walls[p]:
				cells = append(cells, wallStyle.Render(" ■ "))
			case m.SteamTiles[p] > 0:
				cells = append(cells, steamStyle.Render(" ~ "))
			case m.Water[p]:
				switch {
				case isUltCross:
					cells = append(cells, steamStyle.Background(lipgloss.Color("#001a2a")).Render(" ~ "))
				case isUltAxis:
					cells = append(cells, waterStyle.Background(lipgloss.Color("#0d0800")).Render(" ≈ "))
				default:
					cells = append(cells, waterStyle.Render(" ≈ "))
				}
			case m.FireTiles[p] > 0:
				cells = append(cells, fireStyle.Render(" ⁺ "))
			case isUltCross:
				cells = append(cells, ultZoneStyle.Render(" + "))
			case isUltAxis:
				cells = append(cells, ultAxisStyle.Render(" · "))
			case isReachable:
				if m.ShootMode {
					cells = append(cells, shootRangeStyle.Render(" · "))
				} else {
					cells = append(cells, rangeStyle.Render(" · "))
				}
			default:
				cells = append(cells, cellStyle.Render(" · "))
			}
		}
		rows = append(rows, strings.Join(cells, ""))
	}

	info := m.cursorInfo()
	line0 := m.turnOrder()

	line1 := lipgloss.JoinHorizontal(lipgloss.Top,
		modeStr,
		" ",
		hpStr,
		ultStr,
	)

	line2 := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(
			fmt.Sprintf("(%d, %d)  ", m.CursorX, m.CursorY),
		),
		info,
	)

	status := boxStyle.Render(fmt.Sprintf("%s\n%s\n%s", line1, line2, line0))
	grid := strings.Join(rows, "\n")
	box := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, grid))
	helpView := helpStyle.Render(m.help.View(m.keys))
	return lipgloss.JoinVertical(lipgloss.Left, box, status, helpView)
}
