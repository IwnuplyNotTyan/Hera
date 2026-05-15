package generate

import (
	"fmt"
	"strings"

	"hera/utils"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress {
			elem := m.hitTest(msg.X, msg.Y)
			if elem != nil {
				switch elem.Type {
				case ElementGridCell:
					if m.Screen == ScreenGame && !m.EnemyTurn {
						var col, row int
						if _, err := fmt.Sscanf(elem.ID, "cell-%d-%d", &col, &row); err == nil {
							m.CursorX = col
							m.CursorY = row
						}
					}
				case ElementMenuItem:
					if m.Screen == ScreenMenu {
						m.MenuSelected = elem.Index
					}
				case ElementSettingsItem:
					if m.Screen == ScreenSettings {
						m.MenuSelected = elem.Index
					}
				case ElementThemeItem:
					if m.Screen == ScreenThemeSelect {
						m.ThemeName = elem.ID
						if m.Theme != nil {
							m.Theme.SetTintID(elem.ID)
						}
						m.Styles = NewStyles(m.Theme)
					}
				}
			}
		}
		if msg.Button == tea.MouseButtonRight && msg.Action == tea.MouseActionPress {
			if m.Screen == ScreenGame && !m.EnemyTurn {
				m = m.doConfirm()
				if m.Moved && m.Shot {
					m = m.nextTurn()
					if m.CurrentPlayer == 0 {
						m.EnemyTurn = true
						return m, enemyTurnCmd(0)
					}
				}
				return m, nil
			}
			if m.Screen == ScreenMenu || m.Screen == ScreenSettings || m.Screen == ScreenThemeSelect {
				return m.doMenuConfirm()
			}
		}

	case tea.WindowSizeMsg:
		m.TerminalWidth = msg.Width
		m.TerminalHeight = msg.Height
		return m, nil

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
	}

	if m.Screen == ScreenMenu || m.Screen == ScreenSettings || m.Screen == ScreenThemeSelect {
		return m.updateMenu(msg)
	}

	if m.EnemyTurn {
		return m, nil
	}

	if !m.Moved && !m.Shot {
	} else if m.Moved {
		m.ShootMode = true
	} else {
		m.ShootMode = false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
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
			m = m.doConfirm()
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

func (m *Model) renderContent(s string) string {
	if m.TerminalWidth <= 0 || m.TerminalHeight <= 0 {
		return s
	}

	contentWidth := lipgloss.Width(s)
	contentHeight := lipgloss.Height(s)
	if contentWidth > m.TerminalWidth || contentHeight > m.TerminalHeight {
		return m.Localizer.T("error.terminalTooSmall")
	}

	if !m.EnableBackground && !m.CenterWindow {
		m.gridOffsetX = 0
		m.gridOffsetY = 0
		return s
	}

	hPos := lipgloss.Left
	vPos := lipgloss.Top
	if m.CenterWindow {
		hPos = lipgloss.Center
		vPos = lipgloss.Center
		m.gridOffsetX = (m.TerminalWidth - contentWidth) / 2
		m.gridOffsetY = (m.TerminalHeight - contentHeight) / 2
	} else {
		m.gridOffsetX = 0
		m.gridOffsetY = 0
	}

	opts := []lipgloss.WhitespaceOption{}
	if m.EnableBackground {
		opts = append(opts, lipgloss.WithWhitespaceBackground(m.Theme.Bg()))
	}

	return lipgloss.Place(m.TerminalWidth, m.TerminalHeight, hPos, vPos, s, opts...)
}

func (m *Model) View() string {
	m.resetLayout()

	if m.Screen == ScreenMenu {
		return m.viewMenu()
	}
	if m.Screen == ScreenSettings {
		return m.viewSettings()
	}
	if m.Screen == ScreenThemeSelect {
		return m.viewThemeSelect()
	}

	if len(m.Players) == 0 {
		gameOver := m.Styles.BoxStyle.Render(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF4444")).
				Bold(true).
				Render(m.Localizer.T("game.gameOver")),
		)
		return m.renderContent(gameOver)
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

	return m.renderContent(content)
}
