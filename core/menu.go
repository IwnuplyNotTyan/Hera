package generate

import (
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bz "github.com/lrstanley/bubblezone"
)

func (m *Model) startGame() {
	players := []Player{}
	starts := []Point{
		{X: 1, Y: 1},
		{X: GridW - 2, Y: GridH - 2},
		{X: GridW - 2, Y: 1},
		{X: 1, Y: GridH - 2},
	}

	playerCount := rand.Intn(3) + 2
	enemyCount := rand.Intn(3) + 2
	if enemyCount > len(m.Styles.EnemysStyles) {
		enemyCount = len(m.Styles.EnemysStyles)
	}

	for i := 0; i < playerCount; i++ {
		players = append(players, Player{
			X:          starts[i].X,
			Y:          starts[i].Y,
			HP:         MaxHP,
			UltCharges: maxUltCharges,
			Style:      m.Styles.PlayerStyles[i%len(m.Styles.PlayerStyles)],
		})
	}

	blocked := make(map[Point]bool)
	for _, p := range players {
		blocked[Point{p.X, p.Y}] = true
	}

	walls := GenerateTiles(GridW/2, GridH/2, wallCount, blocked)
	for p := range walls {
		blocked[p] = true
	}

	water := GenerateTiles(GridW/2, GridH/2, waterCount, blocked)
	for p := range water {
		blocked[p] = true
	}

	enemyStarts := GenerateTiles(GridW/2, GridH/2, enemyCount, blocked)
	enemyPositions := make([]Point, 0, enemyCount)
	for p := range enemyStarts {
		enemyPositions = append(enemyPositions, p)
	}

	enemys := make([]Enemy, enemyCount)
	for i := range enemys {
		enemys[i] = Enemy{
			X:     enemyPositions[i].X,
			Y:     enemyPositions[i].Y,
			HP:    MaxHP,
			Style: m.Styles.EnemysStyles[i],
		}
	}

	m.Players = players
	m.Enemys = enemys
	m.Walls = walls
	m.Water = water
	m.FireTiles = make(map[Point]int)
	m.SmokeTiles = make(map[Point]int)
	m.CurrentPlayer = 0
	m.CurrentEnemy = 0
	m.CursorX = players[0].X
	m.CursorY = players[0].Y
	m.Moved = false
	m.Shot = false
	m.ShootMode = false
	m.UltMode = false
	m.UltAxis = ""
	m.EnemyTurn = false
	m.EnemyIdx = 0
	m.MenuSelected = 0
	m.Z = bz.New()
}

func (m *Model) navigateTheme(direction int) {
	currentIdx := 0
	for i, t := range m.AvailableThemes {
		if t == m.ThemeName {
			currentIdx = i
			break
		}
	}
	nextIdx := currentIdx + direction
	if nextIdx < 0 {
		nextIdx = len(m.AvailableThemes) - 1
	} else if nextIdx >= len(m.AvailableThemes) {
		nextIdx = 0
	}
	m.ThemeName = m.AvailableThemes[nextIdx]
	if m.Theme != nil {
		m.Theme.SetTintID(m.ThemeName)
	}
	m.Styles = NewStyles(m.Theme)
}

func (m Model) viewMenu() string {
	title := m.Localizer.T("menu.title")
	menuItems := []string{
		m.Localizer.T("menu.start"),
		m.Localizer.T("menu.settings"),
		m.Localizer.T("menu.exit"),
	}
	figures := []string{" ● ", " ■ ", " ◆ "}

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	for i, item := range menuItems {
		figure := figures[i]
		if i == m.MenuSelected {
			style := m.Styles.CursorStyle.Bold(true)
			lines = append(lines, "  "+figure+" "+style.Render(item))
		} else {
			lines = append(lines, "   "+figure+"  "+item)
		}
	}

	easterEgg := m.EasterEgg
	easterBox := m.Styles.BoxStyle.Render(easterEgg)
	lines = append(lines, "")

	menu := lipgloss.JoinVertical(lipgloss.Left, lines...)
	menu = m.Styles.BoxStyle.Render(menu)

	content := lipgloss.JoinVertical(lipgloss.Left, menu, easterBox)

	if m.CenterWindow && m.TerminalWidth > 0 && m.TerminalHeight > 0 {
		contentWidth := lipgloss.Width(content)
		contentHeight := lipgloss.Height(content)
		if contentWidth > m.TerminalWidth || contentHeight > m.TerminalHeight {
			content = m.Localizer.T("error.terminalTooSmall")
		} else {
			marginX := (m.TerminalWidth - contentWidth) / 2
			marginY := (m.TerminalHeight - contentHeight) / 2
			centerStyle := lipgloss.NewStyle().
				MarginLeft(marginX).
				MarginTop(marginY)
			content = centerStyle.Render(content)
		}
	}

	return content
}

func (m Model) viewSettings() string {
	title := m.Localizer.T("settings.title")
	lang := m.Localizer.GetLanguage()
	themeName := m.ThemeName
	centerStr := "on"
	if !m.CenterWindow {
		centerStr = "off"
	}

	menuItems := []string{
		m.Localizer.T("settings.language") + ": " + lang,
		m.Localizer.T("settings.theme") + ": " + themeName,
		m.Localizer.T("settings.center") + ": " + centerStr,
		m.Localizer.T("settings.back"),
	}
	figures := []string{" ● ", " ■ ", " ◆ ", " ● "}

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	for i, item := range menuItems {
		figure := figures[i]
		if i == m.MenuSelected {
			style := m.Styles.CursorStyle.Bold(true)
			lines = append(lines, "  "+figure+" "+style.Render(item))
		} else {
			lines = append(lines, "   "+figure+"  "+item)
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	content = m.Styles.BoxStyle.Render(content)

	if m.CenterWindow && m.TerminalWidth > 0 && m.TerminalHeight > 0 {
		contentWidth := lipgloss.Width(content)
		contentHeight := lipgloss.Height(content)
		if contentWidth > m.TerminalWidth || contentHeight > m.TerminalHeight {
			content = m.Localizer.T("error.terminalTooSmall")
		} else {
			marginX := (m.TerminalWidth - contentWidth) / 2
			marginY := (m.TerminalHeight - contentHeight) / 2
			centerStyle := lipgloss.NewStyle().
				MarginLeft(marginX).
				MarginTop(marginY)
			content = centerStyle.Render(content)
		}
	}

	return content
}

func (m Model) viewThemeSelect() string {
	title := m.Localizer.T("settings.selectTheme")

	themes := m.AvailableThemes
	searchQuery := m.ThemeQuery
	if m.ThemeSearch && m.ThemeQuery != "" {
		var filtered []string
		query := strings.ToLower(m.ThemeQuery)
		for _, t := range m.AvailableThemes {
			if strings.Contains(strings.ToLower(t), query) {
				filtered = append(filtered, t)
			}
		}
		themes = filtered
	} else if m.LastSearchQuery != "" {
		var filtered []string
		query := strings.ToLower(m.LastSearchQuery)
		for _, t := range m.AvailableThemes {
			if strings.Contains(strings.ToLower(t), query) {
				filtered = append(filtered, t)
			}
		}
		themes = filtered
		searchQuery = m.LastSearchQuery
	}

	var lines []string

	if len(themes) == 0 && (m.ThemeSearch || m.LastSearchQuery != "") && searchQuery != "" {
		noResultsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
		lines = append(lines, noResultsStyle.Render("  "+m.Localizer.T("settings.noResults")))
	}

	currentIdx := 0
	if len(themes) > 0 {
		for i, t := range themes {
			if t == m.ThemeName {
				currentIdx = i
				break
			}
		}
	}

	pageSize := 5
	endIdx := currentIdx + pageSize
	if endIdx > len(themes) {
		endIdx = len(themes)
	}
	startIdx := endIdx - pageSize
	if startIdx < 0 {
		startIdx = 0
	}

	lines = append(lines, title)
	lines = append(lines, "")

	for i := startIdx; i < endIdx; i++ {
		theme := themes[i]
		if theme == m.ThemeName {
			style := m.Styles.CursorStyle.Bold(true)
			lines = append(lines, "  ● "+style.Render(theme))
		} else {
			lines = append(lines, "   ●  "+theme)
		}
	}

	themeContent := lipgloss.JoinVertical(lipgloss.Left, lines...)

	var searchContent string
	boxStyle := m.Styles.BoxStyle
	boxStyle = boxStyle.Border(lipgloss.RoundedBorder())
	if m.ThemeSearch {
		searchLine := "/ " + m.ThemeQuery + "_"
		searchStyle := m.Styles.CursorStyle
		searchStyle = searchStyle.Bold(true)
		searchContent = searchStyle.Render(searchLine)
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
		searchContent += "\n" + hintStyle.Render("  esc to close")
		searchContent = boxStyle.Render(searchContent)
	} else if m.LastSearchQuery != "" {
		searchLine := "/ " + m.LastSearchQuery + ""
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
		searchContent = hintStyle.Render(searchLine)
		searchContent = boxStyle.Render(searchContent)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		boxStyle.Render(themeContent),
	)
	if searchContent != "" {
		content = lipgloss.JoinVertical(lipgloss.Left,
			content,
			searchContent,
		)
	}

	if m.CenterWindow && m.TerminalWidth > 0 && m.TerminalHeight > 0 {
		contentWidth := lipgloss.Width(content)
		contentHeight := lipgloss.Height(content)
		if contentWidth > m.TerminalWidth || contentHeight > m.TerminalHeight {
			content = m.Localizer.T("error.terminalTooSmall")
		} else {
			marginX := (m.TerminalWidth - contentWidth) / 2
			marginY := (m.TerminalHeight - contentHeight) / 2
			centerStyle := lipgloss.NewStyle().
				MarginLeft(marginX).
				MarginTop(marginY)
			content = centerStyle.Render(content)
		}
	}

	return content
}

func (m Model) updateMenu(msg tea.Msg) (Model, tea.Cmd) {
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

		switch keyStr {
		case "up", "k", "K":
			if m.Screen == ScreenThemeSelect {
				m.navigateTheme(-1)
			} else {
				m.MenuSelected--
				if m.Screen == ScreenMenu && m.MenuSelected < 0 {
					m.MenuSelected = 2
				} else if m.Screen == ScreenSettings && m.MenuSelected < 0 {
					m.MenuSelected = 3
				}
			}
		case "down", "j", "J":
			if m.Screen == ScreenThemeSelect {
				m.navigateTheme(1)
			} else {
				m.MenuSelected++
				if m.Screen == ScreenMenu && m.MenuSelected > 2 {
					m.MenuSelected = 0
				} else if m.Screen == ScreenSettings && m.MenuSelected > 3 {
					m.MenuSelected = 0
				}
			}
		case "left", "h", "H":
			if m.Screen == ScreenThemeSelect {
				m.navigateTheme(-1)
			} else if m.Screen == ScreenSettings && m.MenuSelected == 0 {
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
				if err := m.Localizer.SetLanguage(languages[currentIdx]); err != nil {
					return m, nil
				}
			}
		case "right", "l", "L":
			if m.Screen == ScreenThemeSelect && !m.ThemeSearch {
				m.navigateTheme(1)
			} else if m.Screen == ScreenSettings && m.MenuSelected == 0 {
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
				if err := m.Localizer.SetLanguage(languages[currentIdx]); err != nil {
					return m, nil
				}
			}
		case "enter", "x", "X":
			if m.Screen == ScreenMenu {
				switch m.MenuSelected {
				case 0:
					m.Screen = ScreenGame
					m.startGame()
				case 1:
					m.Screen = ScreenSettings
					m.MenuSelected = 0
				case 2:
					return m, tea.Quit
				}
			} else if m.Screen == ScreenSettings {
				switch m.MenuSelected {
				case 0:
					languages := m.Localizer.AvailableLanguages()
					currentIdx := 0
					for i, l := range languages {
						if l == m.Localizer.GetLanguage() {
							currentIdx = i
							break
						}
					}
					nextIdx := (currentIdx + 1) % len(languages)
					if err := m.Localizer.SetLanguage(languages[nextIdx]); err != nil {
						return m, nil
					}
				case 1:
					m.Screen = ScreenThemeSelect
					m.MenuSelected = 0
				case 2:
					m.CenterWindow = !m.CenterWindow
				case 3:
					m.Screen = ScreenMenu
					m.MenuSelected = 0
				}
			} else if m.Screen == ScreenThemeSelect {
				m.Screen = ScreenSettings
				m.MenuSelected = 0
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
			}
		}
	}
	return m, nil
}
