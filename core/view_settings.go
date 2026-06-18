package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewSettings() string {
	title := m.Localizer.T("settings.title")
	lang := m.Localizer.GetLanguage()
	themeName := ""
	if m.Config != nil {
		themeName = m.Config.ThemeName
	}
	centerStr := "c"
	if m.Config != nil {
		centerStr = m.Config.CenterWindow
		if centerStr == "" {
			centerStr = "c"
		}
	}
	bgStr := "on"
	if m.Config != nil && !m.Config.Background {
		bgStr = "off"
	}
	debugStr := "off"
	if m.DebugMode {
		debugStr = "on"
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
		m.Localizer.T("settings.debug")+": "+debugStr,
		m.Localizer.T("settings.back"),
	)
	figures = append(figures, " ◆ ", " ◆ ", " » ", " ● ")

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

func (m *Model) viewCenterSelect() string {
	title := m.Localizer.T("settings.selectCenter")

	grid := [3][3]string{
		{"tl", "tc", "tr"},
		{"cl", "c", "cr"},
		{"bl", "bc", "br"},
	}

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	for row := 0; row < 3; row++ {
		var rowCells []string
		for col := 0; col < 3; col++ {
			code := grid[row][col]
			var label string
			if code == "c" {
				label = " C "
			} else {
				label = " " + strings.ToUpper(code) + " "
			}

			selected := row == m.CenterRow && col == m.CenterCol
			box := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(7).
				Align(lipgloss.Center)

			if selected {
				box = box.Background(m.Theme.SelectionBg()).Foreground(m.Theme.Fg())
			} else {
				box = box.Foreground(m.Theme.Fg())
			}

			rowCells = append(rowCells, box.Render(label))

			cellX := col * 9
			cellY := len(lines) + 2
			m.trackElement(Element{
				Type:   ElementCenterItem,
				X:      cellX,
				Y:      cellY,
				Width:  9,
				Height: 3,
				ID:     fmt.Sprintf("center-%d-%d", row, col),
				Index:  row*3 + col,
			})
		}
		rowStr := lipgloss.JoinHorizontal(lipgloss.Top, rowCells...)
		rowStr = lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Render(rowStr)
		lines = append(lines, rowStr)
		lines = append(lines, "")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	content = m.Styles.BoxStyle.Render(content)
	return m.renderContent(content)
}
