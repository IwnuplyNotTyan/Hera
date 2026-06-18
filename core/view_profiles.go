package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewProfiles() string {
	title := m.Localizer.T("menu.profiles")
	gray := lipgloss.NewStyle().Foreground(m.Theme.SelectionBg())
	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	prefix := "   ● "
	for i := 0; i < 3; i++ {
		row := len(lines) + 2
		if m.Profiles[i] != nil {
			var namePart string
			if i == m.ProfileSlot {
				namePart = m.Styles.CursorStyle.Bold(true).Render(m.Profiles[i].Name)
			} else {
				namePart = m.Profiles[i].Name
			}

			if m.ProfileDeleteConfirm && i == m.ProfileSlot {
				var yesBtn, noBtn string
				yesFocused := m.ProfileConfirmChoice == 0 || m.HoveredConfirm == "confirm-yes"
				noFocused := m.ProfileConfirmChoice == 1 || m.HoveredConfirm == "confirm-no"
				if yesFocused {
					yesBtn = " " + lipgloss.NewStyle().Foreground(m.Theme.Green()).Background(m.Theme.SelectionBg()).Render("✓ ")
				} else {
					yesBtn = gray.Render("✓")
				}
				if noFocused {
					noBtn = " " + lipgloss.NewStyle().Foreground(m.Theme.Red()).Background(m.Theme.SelectionBg()).Render("✕ ")
				} else {
					noBtn = gray.Render("✕")
				}
				nameLine := prefix + namePart + "  " + yesBtn + " " + noBtn
				lines = append(lines, nameLine)
				scoreLine := lipgloss.NewStyle().Foreground(m.Theme.Purple()).Render("      " + m.Localizer.T("menu.score") + ": " + fmt.Sprint(m.Profiles[i].Score))
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
				ynStart := 6 + lipgloss.Width(m.Profiles[i].Name) + 4
				m.trackElement(Element{
					Type:   ElementProfileConfirm,
					X:      ynStart,
					Y:      row,
					Width:  lipgloss.Width("[Y]"),
					Height: 1,
					ID:     "confirm-yes",
					Index:  0,
				})
				m.trackElement(Element{
					Type:   ElementProfileConfirm,
					X:      ynStart + lipgloss.Width("[Y]") + 1,
					Y:      row,
					Width:  lipgloss.Width("[N]"),
					Height: 1,
					ID:     "confirm-no",
					Index:  0,
				})
			} else if m.SeedConfirmActive && i == m.ProfileSlot {
				var yesBtn, noBtn string
				if m.SeedConfirmChoice == 0 {
					yesBtn = lipgloss.NewStyle().Background(m.Theme.SelectionBg()).Render(" ? ")
				} else {
					yesBtn = gray.Render("?")
				}
				if m.SeedConfirmChoice == 1 {
					noBtn = lipgloss.NewStyle().Background(m.Theme.SelectionBg()).Render(" ● ")
				} else {
					noBtn = gray.Render("● ")
				}
				nameLine := prefix + namePart + "  " + yesBtn + " " + noBtn
				lines = append(lines, nameLine)
				scoreLine := lipgloss.NewStyle().Foreground(m.Theme.Purple()).Render("      " + m.Localizer.T("menu.score") + ": " + fmt.Sprint(m.Profiles[i].Score))
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
			} else {
				nameLine := prefix + namePart
				lines = append(lines, nameLine)
				scoreLine := lipgloss.NewStyle().Foreground(m.Theme.Cyan()).Render("      " + m.Localizer.T("menu.score") + ": " + fmt.Sprint(m.Profiles[i].Score))
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
			}
		} else {
			text := m.Localizer.T("menu.create")
			if i == m.ProfileSlot {
				lines = append(lines, prefix+m.Styles.CursorStyle.Bold(true).Render(text))
			} else {
				lines = append(lines, prefix+text)
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

	m.help.ShowAll = false
	hint := m.Styles.HelpStyle.Render(m.help.View(m.menuKeys))
	content = lipgloss.JoinVertical(lipgloss.Left, content, hint)

	return m.renderContent(content)
}

func (m *Model) viewProfileCreate() string {
	m.help.ShowAll = false
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

	hint := m.Styles.HelpStyle.Render(m.help.View(m.menuKeys))
	content = lipgloss.JoinVertical(lipgloss.Left, content, hint)

	return m.renderContent(content)
}
