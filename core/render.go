package generate

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) renderContent(s string) string {
	if m.TerminalWidth <= 0 || m.TerminalHeight <= 0 {
		m.gridOffsetX = 0
		m.gridOffsetY = 0
		return s
	}

	contentWidth := lipgloss.Width(s)
	contentHeight := lipgloss.Height(s)
	if contentWidth > m.TerminalWidth || contentHeight > m.TerminalHeight {
		m.gridOffsetX = 0
		m.gridOffsetY = 0
		return m.Localizer.T("error.terminalTooSmall")
	}

	if m.Config != nil && m.Config.CenterWindow != "" {
		var marginX, marginY int

		switch m.Config.CenterWindow {
		case "tl":
			marginX, marginY = 0, 0
		case "tc":
			marginX, marginY = (m.TerminalWidth-contentWidth)/2, 0
		case "tr":
			marginX, marginY = m.TerminalWidth-contentWidth, 0
		case "cl":
			marginX, marginY = 0, (m.TerminalHeight-contentHeight)/2
		case "c":
			marginX, marginY = (m.TerminalWidth-contentWidth)/2, (m.TerminalHeight-contentHeight)/2
		case "cr":
			marginX, marginY = m.TerminalWidth-contentWidth, (m.TerminalHeight-contentHeight)/2
		case "bl":
			marginX, marginY = 0, m.TerminalHeight-contentHeight
		case "bc":
			marginX, marginY = (m.TerminalWidth-contentWidth)/2, m.TerminalHeight-contentHeight
		case "br":
			marginX, marginY = m.TerminalWidth-contentWidth, m.TerminalHeight-contentHeight
		default:
			marginX, marginY = 0, 0
		}

		s = lipgloss.NewStyle().
			MarginLeft(marginX).
			MarginTop(marginY).
			Render(s)
		m.gridOffsetX = marginX
		m.gridOffsetY = marginY
	} else {
		m.gridOffsetX = 0
		m.gridOffsetY = 0
	}

	if m.Config == nil || !m.Config.Background {
		if m.TerminalWidth > contentWidth {
			lines := strings.Split(s, "\n")
			for i, line := range lines {
				if w := lipgloss.Width(line); w < m.TerminalWidth {
					lines[i] = line + strings.Repeat(" ", m.TerminalWidth-w)
				}
			}
			s = strings.Join(lines, "\n")
		}
		return s
	}

	lines := strings.Split(s, "\n")
	bgStyle := lipgloss.NewStyle().Background(m.Theme.Bg())
	for i, line := range lines {
		w := lipgloss.Width(line)
		if w < m.TerminalWidth {
			lines[i] = bgStyle.Render(line + strings.Repeat(" ", m.TerminalWidth-w))
		} else {
			lines[i] = bgStyle.Render(line)
		}
	}
	for len(lines) < m.TerminalHeight {
		lines = append(lines, bgStyle.Render(strings.Repeat(" ", m.TerminalWidth)))
	}

	s = strings.Join(lines, "\n")
	if m.Config != nil && m.Config.Background {
		prefix := strings.SplitN(bgStyle.Render("|"), "|", 2)[0]
		s = strings.ReplaceAll(s, "\x1b[0m", "\x1b[0m"+prefix) + "\x1b[0m"
	}

	return s
}
