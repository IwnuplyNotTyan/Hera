package generate

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) endGame(scoreBonus int) {
	score := m.CurrentScore + scoreBonus
	if score < 0 {
		score = 0
	}
	if !m.SeedLocked && !m.DebugMode {
		if m.ActiveSlot >= 0 && m.ActiveSlot < len(m.Profiles) && m.Profiles[m.ActiveSlot] != nil {
			m.Profiles[m.ActiveSlot].Score += score
		}
		if err := m.SaveProfiles(); err != nil {
			// Save failed — non-fatal, continue
		}
	}
	m.ActiveSlot = -1
	if len(m.Players) == 0 {
		m.Screen = ScreenGameOver
	} else {
		m.Screen = ScreenMenu
	}
	m.MenuSelected = 0
}

func (m *Model) doMenuConfirm() (tea.Model, tea.Cmd) {
	if m.Screen == ScreenMenu {
		switch m.MenuSelected {
		case 0:
			m.Screen = ScreenProfiles
			m.ProfileSlot = 0
		case 1:
			m.Screen = ScreenSettings
			m.MenuSelected = 0
		case 2:
			return m, tea.Quit
		}
	} else if m.Screen == ScreenSettings {
		hasLang := len(m.Localizer.AvailableLanguages()) > 1
		hasThemes := len(m.AvailableThemes) > 1
		n := 0
		if hasLang {
			if m.MenuSelected == n {
				languages := m.Localizer.AvailableLanguages()
				currentIdx := 0
				for i, l := range languages {
					if l == m.Localizer.GetLanguage() {
						currentIdx = i
						break
					}
				}
				nextIdx := (currentIdx + 1) % len(languages)
				if err := m.SetLanguage(languages[nextIdx]); err != nil {
					return m, nil
				}
				return m, nil
			}
			n++
		}
		if hasThemes {
			if m.MenuSelected == n {
				m.Screen = ScreenThemeSelect
				m.MenuSelected = 0
				return m, nil
			}
			n++
		}
		if m.MenuSelected == n {
			m.Screen = ScreenCenterSelect
			m.CenterRow, m.CenterCol = centerToGrid(m.Config.CenterWindow)
			return m, nil
		}
		n++
		if m.MenuSelected == n {
			if m.Config != nil {
				m.Config.Background = !m.Config.Background
				if err := m.SaveConfig(); err != nil {
					m.Config.Background = !m.Config.Background
				}
			}
			return m, nil
		}
		n++
		if m.MenuSelected == n {
			m.DebugMode = !m.DebugMode
			return m, nil
		}
		n++
		if m.MenuSelected == n {
			m.Screen = ScreenMenu
			m.MenuSelected = 0
			return m, nil
		}
		m.Screen = ScreenMenu
		m.MenuSelected = 0
	} else if m.Screen == ScreenThemeSelect {
		m.Screen = ScreenSettings
		m.MenuSelected = 0
	} else if m.Screen == ScreenProfiles {
		slot := m.ProfileSlot
		if m.Profiles[slot] != nil {
			m.ActiveSlot = slot
			m.SeedConfirmActive = true
			m.SeedConfirmChoice = 0
		} else {
			m.ProfileLetters = [3]rune{'A', 'A', 'A'}
			m.ProfileCursor = 0
			m.Screen = ScreenProfileCreate
		}
	}
	return m, nil
}
