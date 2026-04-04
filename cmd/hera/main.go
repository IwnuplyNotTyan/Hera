package main

import (
	"math/rand"

	generate "hera/core"
	"hera/i18n"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	loc, err := i18n.NewTranslator("./locales", "en")
	if err != nil {
		log.Error(err)
		return
	}
	p := tea.NewProgram(generate.NewModel(rand.Intn(3)+2, rand.Intn(3)+2, loc), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Error(err)
	}
}
