package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewGameOver() string {
	var titleLines []string
	if m.Font != nil {
		titleLines = strings.Split(m.Font.Render("GAME OVER"), "\n")
		if len(titleLines) > 0 && strings.TrimSpace(titleLines[len(titleLines)-1]) == "" {
			titleLines = titleLines[:len(titleLines)-1]
		}
	} else {
		titleLines = []string{m.Styles.CursorStyle.Bold(true).Render(m.Localizer.T("game.gameOver"))}
	}

	for i, l := range titleLines {
		titleLines[i] = m.Styles.CursorStyle.Render(l)
	}

	score := m.Localizer.T("menu.score") + ": " + fmt.Sprint(m.CurrentScore)
	if m.Seed != 0 {
		score += "  |  " + m.Localizer.T("game.seedLabel") + ": " + fmt.Sprint(m.Seed)
	}
	prompt := m.Localizer.T("game.anyKey")

	var lines []string
	lines = append(lines, titleLines...)
	lines = append(lines, "", score, "", prompt)

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	content = m.boxStyle().Align(lipgloss.Center).Render(content)
	return m.renderContent(content)
}

func (m *Model) viewSeedPrompt() string {
	title := m.Localizer.T("game.seedLabel")
	warning := m.Localizer.T("game.seedWarning")

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")
	lines = append(lines, warning)
	lines = append(lines, "")
	lines = append(lines, m.Localizer.T("game.seedEnter"))
	lines = append(lines, "")
	lines = append(lines, m.SeedPromptInput.View())

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	content = m.Styles.BoxStyle.Render(content)
	return m.renderContent(content)
}
