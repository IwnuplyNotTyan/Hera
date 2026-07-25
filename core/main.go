package generate

import (
	"fmt"

	"hera/utils"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
							m.BoxTrigger = TriggerNone
							effects := m.effectsAtCursor()
							if len(effects) > 0 {
								m.ShowEffectIdx++
								if m.ShowEffectIdx > len(effects)+1 {
									m.ShowEffectIdx = 0
								}
							} else {
								m.ShowEffectIdx = 0
							}
						}
					}
				case ElementMenuItem:
					if m.Screen == ScreenMenu {
						m.MenuSelected = elem.Index
					}
					if m.Screen == ScreenProfiles {
						m.ProfileSlot = elem.Index
					}
				case ElementSettingsItem:
					if m.Screen == ScreenSettings {
						m.MenuSelected = elem.Index
					}
				case ElementThemeItem:
					if m.Screen == ScreenThemeSelect && m.Config != nil {
						old := m.Config.ThemeName
						m.Config.ThemeName = elem.ID
						if m.Theme != nil {
							m.Theme.SetTintID(elem.ID)
						}
						m.Styles = NewStyles(m.Theme)
						if err := m.SaveConfig(); err != nil {
							m.Config.ThemeName = old
							if m.Theme != nil {
								m.Theme.SetTintID(old)
							}
							m.Styles = NewStyles(m.Theme)
						}
					}
				case ElementCenterItem:
					if m.Screen == ScreenCenterSelect {
						var row, col int
						if _, err := fmt.Sscanf(elem.ID, "center-%d-%d", &row, &col); err == nil {
							if m.Config != nil {
								m.CenterRow = row
								m.CenterCol = col
								old := m.Config.CenterWindow
								m.Config.CenterWindow = gridToCenter(row, col)
								if err := m.SaveConfig(); err != nil {
									m.Config.CenterWindow = old
								}
							}
							m.Screen = ScreenSettings
							m.MenuSelected = 0
						}
					}
				case ElementProfileConfirm:
					if m.Screen == ScreenProfiles {
						if elem.ID == "confirm-yes" {
							m.Profiles[m.ProfileSlot] = nil
							_ = m.SaveProfiles()
						}
						m.ProfileDeleteConfirm = false
					}
				}
			}
		}
		if msg.Button == tea.MouseButtonRight && msg.Action == tea.MouseActionPress {
			if m.Screen == ScreenGame && !m.EnemyTurn {
				m = m.doConfirm()
				if len(m.Players) == 0 {
					m.endGame(3)
					return m, nil
				}
				if len(m.Enemys) == 0 {
					m.endGame(15)
					return m, nil
				}
				if m.Moved && m.Shot {
					m = m.nextTurn()
					if m.CurrentPlayer == 0 {
						m.EnemyTurn = true
						return m, tea.Batch(m.triggerTickCmd(), enemyTurnCmd(0))
					}
				}
				return m, m.triggerTickCmd()
			}
			if m.Screen == ScreenMenu || m.Screen == ScreenSettings || m.Screen == ScreenThemeSelect || m.Screen == ScreenProfiles {
				return m.doMenuConfirm()
			}
		}

		if msg.Action == tea.MouseActionMotion {
			m.HoveredConfirm = ""
			elem := m.hitTest(msg.X, msg.Y)
			if elem != nil && elem.Type == ElementProfileConfirm {
				m.HoveredConfirm = elem.ID
			}
		}

	case tea.WindowSizeMsg:
		m.TerminalWidth = msg.Width
		m.TerminalHeight = msg.Height
		return m, nil

	case enemyTurnMsg:
		if len(m.Players) == 0 {
			m.endGame(3)
			return m, nil
		}
		if msg.enemyIdx >= len(m.Enemys) {
			m.EnemyTurn = false
			m.EnemyIdx = 0
			if len(m.Players) > 0 {
				cur := m.Players[m.CurrentPlayer]
				m.CursorX = cur.X
				m.CursorY = cur.Y
				m.UltMode = false
				m.PushStrikeMode = false
				m.RamMode = false
				m.UltAxis = ""
			}
			return m, m.triggerTickCmd()
		}
		m.EnemyIdx = msg.enemyIdx
		m = m.doEnemyTurn(msg.enemyIdx)
		if len(m.Players) == 0 {
			m.endGame(3)
			return m, nil
		}
		return m, tea.Batch(m.triggerTickCmd(), enemyTurnCmd(msg.enemyIdx+1))

	case triggerTickMsg:
		m.TriggerTimer--
		if m.TriggerTimer <= 0 {
			m.BoxTrigger = TriggerNone
		} else {
			return m, m.triggerTickCmd()
		}
	}

	if m.Screen == ScreenGameOver {
		switch msg.(type) {
		case tea.KeyMsg, tea.MouseMsg:
			m.Screen = ScreenMenu
			m.MenuSelected = 0
			return m, nil
		}
	}

	if m.Screen == ScreenMenu || m.Screen == ScreenSettings || m.Screen == ScreenThemeSelect || m.Screen == ScreenCenterSelect || m.Screen == ScreenProfiles || m.Screen == ScreenProfileCreate || m.Screen == ScreenSeedPrompt {
		return m.updateMenu(msg)
	}

	if m.ConsoleMode {
		return m.UpdateConsole(msg)
	}

	if m.EnemyTurn {
		return m, nil
	}

	if !m.Moved && !m.Shot {
	} else if m.Moved {
		switch m.Players[m.CurrentPlayer].AttackType {
		case AttackPushStrike:
			m.PushStrikeMode = true
		case AttackRam:
			m.RamMode = true
		default:
			m.ShootMode = true
		}
	} else {
		m.ShootMode = false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.ShowEffectIdx = 0
			m.BoxTrigger = TriggerNone
			newY := utils.Clamp(m.CursorY-1, 0, GridH-1)
			if m.UltMode || m.PushStrikeMode || m.RamMode {
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
			m.ShowEffectIdx = 0
			m.BoxTrigger = TriggerNone
			newY := utils.Clamp(m.CursorY+1, 0, GridH-1)
			if m.UltMode || m.PushStrikeMode || m.RamMode {
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
			m.ShowEffectIdx = 0
			m.BoxTrigger = TriggerNone
			newX := utils.Clamp(m.CursorX-1, 0, GridW-1)
			if m.UltMode || m.PushStrikeMode || m.RamMode {
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
			m.ShowEffectIdx = 0
			m.BoxTrigger = TriggerNone
			newX := utils.Clamp(m.CursorX+1, 0, GridW-1)
			if m.UltMode || m.PushStrikeMode || m.RamMode {
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
			m.ShowEffectIdx = 0
			cur := m.Players[m.CurrentPlayer]
			m.CursorX = cur.X
			m.CursorY = cur.Y
			if !m.Shot && m.Players[m.CurrentPlayer].UltCharges > 0 {
				m.UltMode = !m.UltMode
				m.UltAxis = ""
				m.ShootMode = false
				m.PushStrikeMode = false
				m.RamMode = false
				if m.UltMode {
					cur := m.Players[m.CurrentPlayer]
					m.CursorX = cur.X
					m.CursorY = cur.Y
				}
			}

		case key.Matches(msg, m.keys.Shoot):
			m.ShowEffectIdx = 0
			if !m.Shot {
				cur := m.Players[m.CurrentPlayer]
				m.CursorX = cur.X
				m.CursorY = cur.Y
				m.UltMode = false
				m.UltAxis = ""
				switch cur.AttackType {
				case AttackPushStrike:
					m.PushStrikeMode = !m.PushStrikeMode
					m.ShootMode = false
					m.RamMode = false
				case AttackRam:
					m.RamMode = !m.RamMode
					m.ShootMode = false
					m.PushStrikeMode = false
				default:
					m.ShootMode = !m.ShootMode
					m.PushStrikeMode = false
					m.RamMode = false
				}
			}

		case key.Matches(msg, m.keys.Confirm):
			if m.ShowEffectIdx > 0 {
				skipIdx := len(m.effectsAtCursor()) + 1
				if m.ShowEffectIdx == skipIdx {
					m.ShowEffectIdx = 0
					m.Moved = true
					m.Shot = true
					m = m.nextTurn()
					cur := m.Players[m.CurrentPlayer]
					m.CursorX = cur.X
					m.CursorY = cur.Y
					if m.CurrentPlayer == 0 {
						m.EnemyTurn = true
						return m, tea.Batch(m.triggerTickCmd(), enemyTurnCmd(0))
					}
					return m, m.triggerTickCmd()
				}
				m.ShowEffectIdx = 0
				break
			}
			m = m.doConfirm()
			if len(m.Players) == 0 {
				m.endGame(3)
				return m, nil
			}
			if len(m.Enemys) == 0 {
				m.endGame(15)
				return m, nil
			}
			if m.Moved && m.Shot {
				m = m.nextTurn()
				if m.CurrentPlayer == 0 {
					m.EnemyTurn = true
					return m, tea.Batch(m.triggerTickCmd(), enemyTurnCmd(0))
				}
			}

		case key.Matches(msg, m.keys.EffectInfo):
			if m.Screen == ScreenGame {
				if m.ConsoleMode {
					m.ConsoleMode = false
					m.ShowEffectIdx = 0
					m.ConsoleInput.Blur()
					break
				}
				effects := m.effectsAtCursor()
				maxIdx := 0
				if len(effects) > 0 {
					maxIdx = len(effects) + 1
				} else {
					maxIdx = 1
				}
				if m.ShowEffectIdx > 0 {
					m.ShowEffectIdx++
					if m.ShowEffectIdx > maxIdx {
						m.ShowEffectIdx = 0
						if m.DebugMode {
							m.ConsoleMode = true
							m.help.ShowAll = false
							fi := m.ConsoleInput.Focus()
							return m, fi
						}
					}
				} else {
					m.ShowEffectIdx = 1
				}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			if m.Screen == ScreenGame && m.ActiveSlot >= 0 {
				m.endGame(0)
				return m, nil
			}
			return m, tea.Quit
		}
	}
	return m, m.triggerTickCmd()
}
