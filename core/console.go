package generate

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConsoleCmd struct {
	Raw     string
	Command string
	Args    []string
}

func ParseConsoleCmd(input string) ConsoleCmd {
	input = strings.TrimSpace(input)
	if input == "" {
		raw := "help"
		return ConsoleCmd{Raw: raw, Command: "help"}
	}
	parts := strings.Fields(input)
	cmd := strings.ToLower(parts[0])
	args := parts[1:]
	return ConsoleCmd{Raw: input, Command: cmd, Args: args}
}

func ParseTarget(s string) (isPlayer bool, idx int, ok bool) {
	if len(s) < 2 {
		return false, 0, false
	}
	prefix := strings.ToLower(string(s[0]))
	n, err := strconv.Atoi(s[1:])
	if err != nil || n < 1 {
		return false, 0, false
	}
	switch prefix {
	case "p":
		return true, n - 1, true
	case "e":
		return false, n - 1, true
	default:
		return false, 0, false
	}
}

func (m *Model) InitConsole() {
	ti := textinput.New()
	ti.Placeholder = "type command... (help)"
	ti.CharLimit = 200
	ti.Width = 44
	ti.Prompt = "> "
	m.ConsoleInput = ti
	m.ConsoleOutput = nil
	m.ConsoleHistory = nil
	m.ConsoleHistoryIdx = -1
}

func (m *Model) addConsoleOutput(line string) {
	if line == "" {
		return
	}
	m.ConsoleOutput = append(m.ConsoleOutput, line)
	if len(m.ConsoleOutput) > 50 {
		m.ConsoleOutput = m.ConsoleOutput[len(m.ConsoleOutput)-50:]
	}
}

func (m *Model) UpdateConsole(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			input := m.ConsoleInput.Value()
			if input == "" {
				return m, nil
			}
			cmd := ParseConsoleCmd(input)
			result := m.ExecuteConsoleCmd(cmd)
			m.ConsoleHistory = append(m.ConsoleHistory, input)
			if len(m.ConsoleHistory) > 50 {
				m.ConsoleHistory = m.ConsoleHistory[1:]
			}
			m.ConsoleHistoryIdx = -1
			m.ConsoleInput.SetValue("")
			m.addConsoleOutput("> " + input)
			if result != "" {
				for _, line := range strings.Split(result, "\n") {
					m.addConsoleOutput(line)
				}
			}
			return m, m.ConsoleInput.Focus()
		case "esc", "q":
			m.ConsoleMode = false
			m.ShowEffectIdx = 0
			m.ConsoleInput.Blur()
			return m, nil
		case "up":
			if len(m.ConsoleHistory) == 0 {
				return m, nil
			}
			if m.ConsoleHistoryIdx < len(m.ConsoleHistory)-1 {
				m.ConsoleHistoryIdx++
			}
			idx := len(m.ConsoleHistory) - 1 - m.ConsoleHistoryIdx
			m.ConsoleInput.SetValue(m.ConsoleHistory[idx])
			m.ConsoleInput.SetCursor(len(m.ConsoleHistory[idx]))
			return m, nil
		case "down":
			if m.ConsoleHistoryIdx <= 0 {
				m.ConsoleHistoryIdx = -1
				m.ConsoleInput.SetValue("")
				return m, nil
			}
			m.ConsoleHistoryIdx--
			idx := len(m.ConsoleHistory) - 1 - m.ConsoleHistoryIdx
			m.ConsoleInput.SetValue(m.ConsoleHistory[idx])
			m.ConsoleInput.SetCursor(len(m.ConsoleHistory[idx]))
			return m, nil
		default:
			var cmd tea.Cmd
			m.ConsoleInput, cmd = m.ConsoleInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *Model) ConsoleView() string {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.Theme.Purple()).
		Padding(0, 2).
		Width(46)

	maxOut := 15
	if m.TerminalHeight > 0 {
		available := m.TerminalHeight - 25
		if available < 3 {
			available = 3
		}
		if available < maxOut {
			maxOut = available
		}
	}

	var lines []string
	start := 0
	if len(m.ConsoleOutput) > maxOut {
		start = len(m.ConsoleOutput) - maxOut
	}
	lines = append(lines, m.ConsoleOutput[start:]...)
	lines = append(lines, m.ConsoleInput.View())
	return borderStyle.Render(strings.Join(lines, "\n"))
}
