package generate

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func (m *Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TerminalWidth = msg.Width
		m.TerminalHeight = msg.Height
		return m, nil

	case tea.KeyMsg:
		keyStr := msg.String()

		if m.Screen == ScreenThemeSelect && m.ThemeSearch {
			switch keyStr {
			case "esc":
				m.ThemeSearch = false
				m.ThemeQuery = ""
			case "backspace":
				if len(m.ThemeQuery) > 0 {
					m.ThemeQuery = m.ThemeQuery[:len(m.ThemeQuery)-1]
				} else {
					m.ThemeSearch = false
				}
			case "enter":
				m.ThemeSearch = false
				m.ThemeQuery = ""
			case "/":
				m.ThemeSearch = false
				m.ThemeQuery = ""
			default:
				if len(keyStr) == 1 {
					r := []rune(keyStr)[0]
					if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ' ' || r == '-' || r == '_' {
						m.ThemeQuery += keyStr
					}
				}
			}
			return m, nil
		}

		if m.Screen == ScreenProfileCreate {
			switch keyStr {
			case "j", "J":
				m.ProfileLetters[m.ProfileCursor]++
				if m.ProfileLetters[m.ProfileCursor] > 'Z' {
					m.ProfileLetters[m.ProfileCursor] = 'A'
				}
			case "k", "K":
				m.ProfileLetters[m.ProfileCursor]--
				if m.ProfileLetters[m.ProfileCursor] < 'A' {
					m.ProfileLetters[m.ProfileCursor] = 'Z'
				}
			case "h", "H":
				m.ProfileCursor--
				if m.ProfileCursor < 0 {
					m.ProfileCursor = 2
				}
			case "l", "L":
				m.ProfileCursor++
				if m.ProfileCursor > 2 {
					m.ProfileCursor = 0
				}
			case "enter", "x", "X":
				name := string(m.ProfileLetters[:])
				m.Profiles[m.ProfileSlot] = &Profile{Name: name, Score: 0}
				if err := m.SaveProfiles(); err != nil {
					log.Error("failed to save profiles", "err", err)
				}
				m.Screen = ScreenProfiles
			case "esc", "q":
				m.Screen = ScreenProfiles
			}
			return m, nil
		}

		if m.Screen == ScreenCenterSelect {
			switch keyStr {
			case "up", "k", "K":
				m.CenterRow--
				if m.CenterRow < 0 {
					m.CenterRow = 2
				}
			case "down", "j", "J":
				m.CenterRow++
				if m.CenterRow > 2 {
					m.CenterRow = 0
				}
			case "left", "h", "H":
				m.CenterCol--
				if m.CenterCol < 0 {
					m.CenterCol = 2
				}
			case "right", "l", "L":
				m.CenterCol++
				if m.CenterCol > 2 {
					m.CenterCol = 0
				}
			case "enter", "x", "X":
				if m.Config != nil {
					old := m.Config.CenterWindow
					m.Config.CenterWindow = gridToCenter(m.CenterRow, m.CenterCol)
					if err := m.SaveConfig(); err != nil {
						m.Config.CenterWindow = old
					}
				}
				m.Screen = ScreenSettings
				m.MenuSelected = 0
			case "esc", "q":
				m.Screen = ScreenSettings
				m.MenuSelected = 0
			}
			return m, nil
		}

		if m.Screen == ScreenSeedPrompt {
			var cmd tea.Cmd
			m.SeedPromptInput, cmd = m.SeedPromptInput.Update(msg)

			switch keyStr {
			case "enter":
				val := m.SeedPromptInput.Value()
				if val != "" {
					m.Seed = ParseSeed(val)
					m.SeedLocked = true
				}
				m.Screen = ScreenGame
				m.CurrentScore = 0
				m.SeedConfirmChoice = 0
				m.SeedPromptInput.SetValue("")
				m.startGame()
				return m, cmd
			case "esc", "q":
				m.Screen = ScreenProfiles
				m.SeedConfirmChoice = 0
				m.SeedPromptInput.SetValue("")
				return m, cmd
			}
			return m, cmd
		}

		switch keyStr {
		case "up", "k", "K":
			switch m.Screen {
			case ScreenThemeSelect:
				m.navigateTheme(-1)
			case ScreenProfiles:
				m.ProfileSlot--
				if m.ProfileSlot < 0 {
					m.ProfileSlot = 2
				}
			default:
				m.MenuSelected--
				if m.Screen == ScreenMenu && m.MenuSelected < 0 {
					m.MenuSelected = 2
				} else if m.Screen == ScreenSettings && m.MenuSelected < 0 {
					maxItem := 3
					if len(m.Localizer.AvailableLanguages()) > 1 {
						maxItem++
					}
					if len(m.AvailableThemes) > 1 {
						maxItem++
					}
					m.MenuSelected = maxItem
				}
			}
		case "down", "j", "J":
			switch m.Screen {
			case ScreenThemeSelect:
				m.navigateTheme(1)
			case ScreenProfiles:
				m.ProfileSlot++
				if m.ProfileSlot > 2 {
					m.ProfileSlot = 0
				}
			default:
				m.MenuSelected++
				if m.Screen == ScreenMenu && m.MenuSelected > 2 {
					m.MenuSelected = 0
				} else if m.Screen == ScreenSettings {
					maxItem := 3
					if len(m.Localizer.AvailableLanguages()) > 1 {
						maxItem++
					}
					if len(m.AvailableThemes) > 1 {
						maxItem++
					}
					if m.MenuSelected > maxItem {
						m.MenuSelected = 0
					}
				}
			}
		case "left", "h", "H":
			if m.Screen == ScreenProfiles {
				if m.SeedConfirmActive {
					m.SeedConfirmChoice = 0
				} else if m.ProfileDeleteConfirm {
					m.ProfileConfirmChoice = 0
				} else {
					m.ProfileSlot--
					if m.ProfileSlot < 0 {
						m.ProfileSlot = 2
					}
				}
			} else if m.Screen == ScreenThemeSelect {
				m.navigateTheme(-1)
			} else if m.Screen == ScreenSettings && m.MenuSelected == 0 && len(m.Localizer.AvailableLanguages()) > 1 {
				languages := m.Localizer.AvailableLanguages()
				currentIdx := len(languages) - 1
				for i, l := range languages {
					if l == m.Localizer.GetLanguage() {
						currentIdx = i
						break
					}
				}
				currentIdx--
				if currentIdx < 0 {
					currentIdx = len(languages) - 1
				}
				if err := m.SetLanguage(languages[currentIdx]); err != nil {
					return m, nil
				}
			}
		case "right", "l", "L":
			if m.Screen == ScreenProfiles {
				if m.SeedConfirmActive {
					m.SeedConfirmChoice = 1
				} else if m.ProfileDeleteConfirm {
					m.ProfileConfirmChoice = 1
				} else {
					m.ProfileSlot++
					if m.ProfileSlot > 2 {
						m.ProfileSlot = 0
					}
				}
			} else if m.Screen == ScreenThemeSelect && !m.ThemeSearch {
				m.navigateTheme(1)
			} else if m.Screen == ScreenSettings && m.MenuSelected == 0 && len(m.Localizer.AvailableLanguages()) > 1 {
				languages := m.Localizer.AvailableLanguages()
				currentIdx := 0
				for i, l := range languages {
					if l == m.Localizer.GetLanguage() {
						currentIdx = i
						break
					}
				}
				currentIdx++
				if currentIdx >= len(languages) {
					currentIdx = 0
				}
				if err := m.SetLanguage(languages[currentIdx]); err != nil {
					return m, nil
				}
			}
		case "enter", "x", "X":
			if m.Screen == ScreenProfiles && m.ProfileDeleteConfirm {
				if m.ProfileConfirmChoice == 0 {
					m.Profiles[m.ProfileSlot] = nil
					if err := m.SaveProfiles(); err != nil {
						log.Error("failed to save profiles", "err", err)
					}
				}
				m.ProfileDeleteConfirm = false
				return m, nil
			}
			if m.Screen == ScreenProfiles && m.SeedConfirmActive {
				m.SeedConfirmActive = false
				if m.SeedConfirmChoice == 0 {
					m.Screen = ScreenGame
					m.CurrentScore = 0
					m.SeedConfirmChoice = 0
					m.SeedPromptInput.SetValue("")
					m.startGame()
				} else {
					m.Screen = ScreenSeedPrompt
					ti := textinput.New()
					ti.Placeholder = "123456789"
					ti.CharLimit = 64
					ti.Width = 30
					ti.Prompt = ""
					m.SeedPromptInput = ti
					cmd := m.SeedPromptInput.Focus()
					return m, cmd
				}
				return m, nil
			}
			return m.doMenuConfirm()
		case "backspace", "z", "Z":
			if m.Screen == ScreenProfiles {
				if m.SeedConfirmActive {
					m.SeedConfirmActive = false
				} else if m.ProfileDeleteConfirm {
					m.ProfileDeleteConfirm = false
				} else if m.Profiles[m.ProfileSlot] != nil {
					m.ProfileDeleteConfirm = true
					m.ProfileConfirmChoice = 0
				}
			}
		case "esc", "q":
			switch m.Screen {
			case ScreenSettings:
				m.Screen = ScreenMenu
				m.MenuSelected = 0
			case ScreenThemeSelect:
				if m.ThemeSearch {
					m.LastSearchQuery = m.ThemeQuery
					m.ThemeSearch = false
				} else {
					m.Screen = ScreenSettings
					m.MenuSelected = 0
				}
			case ScreenMenu:
				return m, tea.Quit
			case ScreenProfiles:
				if m.SeedConfirmActive {
					m.SeedConfirmActive = false
				} else {
					m.Screen = ScreenMenu
					m.MenuSelected = 0
				}
			case ScreenProfileCreate:
				m.Screen = ScreenProfiles
			case ScreenSeedPrompt:
				m.Screen = ScreenProfiles
				m.SeedConfirmChoice = 0
				m.SeedPromptInput.SetValue("")
			}
		}
	}
	return m, nil
}
