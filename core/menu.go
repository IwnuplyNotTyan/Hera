package generate

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) startGame() {
	players := []Player{}
	starts := []Point{
		{X: 1, Y: 1},
		{X: GridW - 2, Y: GridH - 2},
		{X: GridW - 2, Y: 1},
		{X: 1, Y: GridH - 2},
	}

	playerCount := m.startPlayers
	if playerCount <= 0 {
		playerCount = rand.Intn(3) + 2
	}
	enemyCount := m.startEnemies
	if enemyCount <= 0 {
		enemyCount = rand.Intn(3) + 2
	}
	if enemyCount > len(m.Styles.EnemysStyles) {
		enemyCount = len(m.Styles.EnemysStyles)
	}

	for i := 0; i < playerCount; i++ {
		effs := make([]Effect, len(m.startPlayerEffects))
		copy(effs, m.startPlayerEffects)
		players = append(players, Player{
			X:          starts[i].X,
			Y:          starts[i].Y,
			HP:         MaxHP,
			UltCharges: maxUltCharges,
			Effects:    effs,
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
		effs := make([]Effect, len(m.startEnemyEffects))
		copy(effs, m.startEnemyEffects)
		enemys[i] = Enemy{
			X:       enemyPositions[i].X,
			Y:       enemyPositions[i].Y,
			HP:      MaxHP,
			Effects: effs,
			Style:   m.Styles.EnemysStyles[i],
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
}

func (m *Model) navigateTheme(direction int) {
	currentIdx := 0
	for i, t := range m.AvailableThemes {
		if m.Config != nil && t == m.Config.ThemeName {
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
	if m.Config != nil {
		oldTheme := m.Config.ThemeName
		m.Config.ThemeName = m.AvailableThemes[nextIdx]
		if m.Theme != nil {
			m.Theme.SetTintID(m.Config.ThemeName)
		}
		m.Styles = NewStyles(m.Theme)
		if err := m.SaveConfig(); err != nil {
			m.Config.ThemeName = oldTheme
			if m.Theme != nil {
				m.Theme.SetTintID(oldTheme)
			}
			m.Styles = NewStyles(m.Theme)
		}
	}
}

func (m *Model) viewMenu() string {
	title := m.Localizer.T("menu.title")
	if m.BannerText != "" {
		title = m.Styles.CursorStyle.Render(m.BannerText)
	}
	menuItems := []string{
		m.Localizer.T("menu.start"),
		m.Localizer.T("menu.settings"),
		m.Localizer.T("menu.exit"),
	}
	figures := []string{" ● ", " ■ ", " ◆ "}

	var lines []string
	renderedLines := 0
	lines = append(lines, title)
	renderedLines += strings.Count(title, "\n") + 1
	lines = append(lines, "")
	renderedLines++

	for i, item := range menuItems {
		figure := figures[i]
		row := renderedLines + 2
		itemWidth := lipgloss.Width(figure) + 1 + lipgloss.Width(item)
		if i == m.MenuSelected {
			style := m.Styles.CursorStyle.Bold(true)
			lines = append(lines, "  "+figure+" "+style.Render(item))
		} else {
			lines = append(lines, "   "+figure+"  "+item)
		}
		m.trackElement(Element{
			Type:   ElementMenuItem,
			X:      4,
			Y:      row,
			Width:  itemWidth,
			Height: 1,
			ID:     fmt.Sprintf("menu-%d", i),
			Index:  i,
		})
		renderedLines++
	}

	easterEgg := m.EasterEgg
	easterBox := m.Styles.BoxStyle.Render(easterEgg)
	lines = append(lines, "")

	menu := lipgloss.JoinVertical(lipgloss.Left, lines...)
	menu = m.Styles.BoxStyle.Render(menu)

	content := lipgloss.JoinVertical(lipgloss.Left, menu, easterBox)

	return m.renderContent(content)
}

func (m *Model) viewSettings() string {
	title := m.Localizer.T("settings.title")
	lang := m.Localizer.GetLanguage()
	themeName := ""
	if m.Config != nil {
		themeName = m.Config.ThemeName
	}
	centerStr := "on"
	if m.Config != nil && !m.Config.CenterWindow {
		centerStr = "off"
	}
	bgStr := "on"
	if m.Config != nil && !m.Config.Background {
		bgStr = "off"
	}

	hasLang := len(m.Localizer.AvailableLanguages()) > 1
	hasThemes := len(m.AvailableThemes) > 1
	menuItems := []string{}
	figures := []string{}
	if hasLang {
		menuItems = append(menuItems, m.Localizer.T("settings.language")+": "+lang)
		figures = append(figures, " ● ")
	}
	if hasThemes {
		menuItems = append(menuItems, m.Localizer.T("settings.theme")+": "+themeName)
		figures = append(figures, " ■ ")
	}
	menuItems = append(menuItems,
		m.Localizer.T("settings.center")+": "+centerStr,
		m.Localizer.T("settings.background")+": "+bgStr,
		m.Localizer.T("settings.back"),
	)
	figures = append(figures, " ◆ ", " ◆ ", " ● ")

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	for i, item := range menuItems {
		figure := figures[i]
		row := len(lines) + 2
		itemWidth := lipgloss.Width(figure) + 1 + lipgloss.Width(item)
		if i == m.MenuSelected {
			style := m.Styles.CursorStyle.Bold(true)
			lines = append(lines, "  "+figure+" "+style.Render(item))
		} else {
			lines = append(lines, "   "+figure+"  "+item)
		}
		m.trackElement(Element{
			Type:   ElementSettingsItem,
			X:      4,
			Y:      row,
			Width:  itemWidth,
			Height: 1,
			ID:     fmt.Sprintf("settings-%d", i),
			Index:  i,
		})
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	content = m.Styles.BoxStyle.Render(content)

	return m.renderContent(content)
}

func (m *Model) viewThemeSelect() string {
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
			if m.Config != nil && t == m.Config.ThemeName {
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
		row := len(lines) + 2
		if (m.Config != nil && theme == m.Config.ThemeName) || (m.Config == nil && theme == "default") {
			style := m.Styles.CursorStyle.Bold(true)
			lines = append(lines, "  ● "+style.Render(theme))
		} else {
			lines = append(lines, "   ●  "+theme)
		}
		m.trackElement(Element{
			Type:   ElementThemeItem,
			X:      4,
			Y:      row,
			Width:  lipgloss.Width("   ●  ") + lipgloss.Width(theme),
			Height: 1,
			ID:     theme,
			Index:  i,
		})
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

	return m.renderContent(content)
}

func (m *Model) viewProfiles() string {
	title := m.Localizer.T("menu.title")
	gray := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	for i := 0; i < 3; i++ {
		row := len(lines) + 2
		if m.Profiles[i] != nil {
			del := gray.Render("✕")
			var nameLine string
			if i == m.ProfileSlot {
				nameLine = "  ●  " + m.Styles.CursorStyle.Bold(true).Render(m.Profiles[i].Name) + "  " + del
			} else {
				nameLine = "   ● " + m.Profiles[i].Name + "  " + del
			}
			lines = append(lines, nameLine)
			scoreLine := gray.Render("      " + m.Localizer.T("menu.score") + ": " + fmt.Sprint(m.Profiles[i].Score))
			lines = append(lines, scoreLine)

			m.trackElement(Element{
				Type:   ElementMenuItem,
				X:      4,
				Y:      row,
				Width:  lipgloss.Width(m.Profiles[i].Name) + 6,
				Height: 1,
				ID:     fmt.Sprintf("profile-%d", i),
				Index:  i,
			})
			m.trackElement(Element{
				Type:   ElementProfileDelete,
				X:      lipgloss.Width(nameLine) - lipgloss.Width(del) + 4,
				Y:      row,
				Width:  lipgloss.Width(del),
				Height: 1,
				ID:     fmt.Sprintf("delete-%d", i),
				Index:  i,
			})
		} else {
			text := m.Localizer.T("menu.create")
			if i == m.ProfileSlot {
				lines = append(lines, "  ●  "+m.Styles.CursorStyle.Bold(true).Render(text))
			} else {
				lines = append(lines, "   ● "+text)
			}
			lines = append(lines, "")

			m.trackElement(Element{
				Type:   ElementMenuItem,
				X:      4,
				Y:      row,
				Width:  lipgloss.Width(text) + 6,
				Height: 1,
				ID:     fmt.Sprintf("profile-%d", i),
				Index:  i,
			})
		}
	}

	lines = append(lines, "")
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	content = m.Styles.BoxStyle.Render(content)
	return m.renderContent(content)
}

func (m *Model) viewProfileCreate() string {
	title := m.Localizer.T("menu.createProfile")

	letterRenders := make([][]string, 3)
	for i, r := range m.ProfileLetters {
		if m.Font != nil {
			rendered := m.Font.Render(string(r))
			letterRenders[i] = strings.Split(rendered, "\n")
		} else {
			letterRenders[i] = []string{string(r), "", "", ""}
		}
	}

	var letterLines []string
	for row := 0; row < 4; row++ {
		var combined string
		for i := 0; i < 3; i++ {
			part := ""
			if row < len(letterRenders[i]) {
				part = letterRenders[i][row]
			}
			if i == m.ProfileCursor && strings.TrimSpace(part) != "" {
				combined += m.Styles.CursorStyle.Render(part) + "  "
			} else {
				combined += part + "  "
			}
		}
		letterLines = append(letterLines, combined)
	}

	var boxLines []string
	boxLines = append(boxLines, title)
	boxLines = append(boxLines, "")
	boxLines = append(boxLines, letterLines...)
	boxLines = append(boxLines, "")

	minW := 42
	maxW := 0
	for _, l := range boxLines {
		if w := lipgloss.Width(l); w > maxW {
			maxW = w
		}
	}
	if maxW < minW {
		maxW = minW
	}
	for i, l := range boxLines {
		if w := lipgloss.Width(l); w < maxW {
			pad := (maxW - w) / 2
			boxLines[i] = strings.Repeat(" ", pad) + l
		}
	}

	boxContent := lipgloss.JoinVertical(lipgloss.Left, boxLines...)
	content := m.Styles.BoxStyle.Render(boxContent)

	hint := m.Localizer.T("menu.profileHint")
	hint = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(hint)
	content = lipgloss.JoinVertical(lipgloss.Left, content, hint)

	return m.renderContent(content)
}

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
				m.Screen = ScreenProfiles
			case "esc", "q":
				m.Screen = ScreenProfiles
			}
			return m, nil
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
					maxItem := 2
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
					maxItem := 2
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
				m.ProfileSlot--
				if m.ProfileSlot < 0 {
					m.ProfileSlot = 2
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
				if err := m.Localizer.SetLanguage(languages[currentIdx]); err != nil {
					return m, nil
				}
			}
		case "right", "l", "L":
			if m.Screen == ScreenProfiles {
				m.ProfileSlot++
				if m.ProfileSlot > 2 {
					m.ProfileSlot = 0
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
				if err := m.Localizer.SetLanguage(languages[currentIdx]); err != nil {
					return m, nil
				}
			}
		case "enter", "x", "X":
			return m.doMenuConfirm()
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
				m.Screen = ScreenMenu
				m.MenuSelected = 0
			case ScreenProfileCreate:
				m.Screen = ScreenProfiles
			}
		}
	}
	return m, nil
}
