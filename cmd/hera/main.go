package main

import (
	"context"
	"fmt"
	"math/rand"

	generate "hera/core"
	"hera/i18n"

	"charm.land/fang/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	tint "github.com/lrstanley/bubbletint"
	"github.com/spf13/cobra"
)

func main() {
	var lang string
	var theme string
	var noCenter bool

	cmd := &cobra.Command{
		Use:   "hera",
		Short: "Hera - A tactical turn-based game",
		Long:  "A tactical turn-based game. Made with ♡",
		RunE: func(cmd *cobra.Command, _ []string) error {
			tint.NewDefaultRegistry()
			generate.RegisterDefaults()
			registry := tint.DefaultRegistry
			if theme != "" {
				if ok := registry.SetTintID(theme); !ok {
					return fmt.Errorf("unknown theme: %s", theme)
				}
			} else {
				registry.SetTintID("default")
			}
			loc, err := i18n.NewTranslator("./i18n/locales", lang)
			if err != nil {
				return err
			}
			centerWindow := !noCenter
			themeName := theme
			if themeName == "" {
				themeName = "default"
			}
			model := generate.NewModel(rand.Intn(3)+2, rand.Intn(3)+2, loc, registry, centerWindow, themeName)
			model.SetAvailableThemes()
			p := tea.NewProgram(
				model,
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)
			if _, err := p.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lang, "lang", "l", "en", "Language code (en, ru)")
	cmd.Flags().StringVarP(&theme, "theme", "t", "", "Theme name (e.g., dracula, tokyonight, gruvbox)")
	cmd.Flags().BoolVarP(&noCenter, "no-center", "c", false, "Disable centered window")

	if err := fang.Execute(context.Background(), cmd); err != nil {
		log.Error(err)
	}
}
