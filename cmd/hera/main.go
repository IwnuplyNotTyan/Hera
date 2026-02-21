package main

import (
	"hera/generate"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	p := tea.NewProgram(generate.NewModel())
	if _, err := p.Run(); err != nil {
		log.Error(err)
	}
}
