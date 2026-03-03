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
			if m.inRange(m.CursorX, newY) {
				m.CursorY = newY
			}
		case key.Matches(msg, m.keys.Down):
			newY := utils.Clamp(m.CursorY+1, 0, GridH-1)
			if m.inRange(m.CursorX, newY) {
				m.CursorY = newY
			}
		case key.Matches(msg, m.keys.Left):
			newX := utils.Clamp(m.CursorX-1, 0, GridW-1)
			if m.inRange(newX, m.CursorY) {
				m.CursorX = newX
			}
		case key.Matches(msg, m.keys.Right):
			newX := utils.Clamp(m.CursorX+1, 0, GridW-1)
			if m.inRange(newX, m.CursorY) {
				m.CursorX = newX
			}
		case key.Matches(msg, m.keys.Shoot):
			if !m.Shot {
				m.ShootMode = !m.ShootMode
				current := m.Players[m.CurrentPlayer]
				m.CursorX = current.X
				m.CursorY = current.Y
			}
		case key.Matches(msg, m.keys.Confirm):
			p := Point{m.CursorX, m.CursorY}
			current := m.Players[m.CurrentPlayer]
			wallBlocked := m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

			if m.ShootMode && !m.Shot {
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
			} else if !m.ShootMode && !m.Moved {
				if !m.Walls[p] && !wallBlocked && !m.OccupiedByOther(m.CursorX, m.CursorY) {
					m.Players[m.CurrentPlayer].X = m.CursorX
					m.Players[m.CurrentPlayer].Y = m.CursorY
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
	current := m.Players[m.CurrentPlayer]
	hp := strings.Repeat("♥ ", current.HP) + strings.Repeat("♡ ", MaxHP-current.HP)
	hpStyle := current.Style
	if current.HP == 1 {
		hpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Blink(true)
	}
	hpStr := hpStyle.Render(fmt.Sprintf("Player %d  %s", m.CurrentPlayer+1, hp))

	var modeStr string
	if m.ShootMode {
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Bold(true).
			Render("♡ S ")
	} else {
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Render("♧ M ")
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
			switch {
			case col == m.CursorX && row == m.CursorY:
				if playerIdx >= 0 {
					cells = append(cells, cursorStyle.Render(
						m.Players[playerIdx].Style.Render(" ■ "),
					))
				} else if enemyIdx >= 0 {
					cells = append(cells, cursorStyle.Render(
						m.Enemys[enemyIdx].Style.Render(" ▲ "),
					))
				} else {
					cells = append(cells, cursorStyle.Render(" · "))
				}
			case playerIdx >= 0:
				symbol := " ■ "
				if playerIdx == m.CurrentPlayer {
					symbol = " ● "
				}
				cells = append(cells, m.Players[playerIdx].Style.Render(symbol))
			case enemyIdx >= 0:
				symbol := " ▲ "
				if enemyIdx == m.CurrentEnemy {
					symbol = " ♦ "
				}
				cells = append(cells, m.Enemys[enemyIdx].Style.Render(symbol))
			case m.Walls[p]:
				cells = append(cells, wallStyle.Render(" ■ "))
			case m.Water[p]:
				cells = append(cells, waterStyle.Render(" ≈ "))
			case m.IsInRange(col, row):
				cells = append(cells, rangeStyle.Render(" · "))
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
