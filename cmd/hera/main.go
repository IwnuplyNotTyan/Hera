package main

import (
	"context"
	"math/rand"

	generate "hera/core"
	"hera/i18n"

	"charm.land/fang/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func main() {
	var lang string

	cmd := &cobra.Command{
		Use:   "hera",
		Short: "Hera - A tactical turn-based game",
		Long:  "A tactical turn-based game. Made with ♡",
		RunE: func(cmd *cobra.Command, _ []string) error {
			loc, err := i18n.NewTranslator("./i18n/locales", lang)
			if err != nil {
				return err
			}
			p := tea.NewProgram(generate.NewModel(rand.Intn(3)+2, rand.Intn(3)+2, loc), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lang, "lang", "l", "en", "Language code (en, ru)")

	if err := fang.Execute(context.Background(), cmd); err != nil {
		log.Error(err)
	}
}
