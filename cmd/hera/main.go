package main

import (
	"math/rand"

	"hera/generate"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	p := tea.NewProgram(generate.NewModel(rand.Intn(3) + 2), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Error(err)
	}
}
