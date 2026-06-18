package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
