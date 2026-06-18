package main

import (
	"context"
	"fmt"
	"strings"

	"hera/assets"
	generate "hera/core"
	"hera/i18n"

	"charm.land/fang/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var version = "dev"
var commit = ""

func parseEffects(s string) []generate.Effect {
	if s == "" {
		return nil
	}
	var out []generate.Effect
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		var t generate.EffectType
		switch p {
		case "fire":
			t = generate.EffectFire
		case "wet":
			t = generate.EffectWet
		case "smoke":
			t = generate.EffectSmoke
		default:
			continue
		}
		out = append(out, generate.Effect{Type: t, Duration: 3})
	}
	return out
}

func main() {
	var lang string
	var theme string
	var centerMode string
	var background bool
	var flagPlayers, flagEnemies int
	var flagPlayerEffects, flagEnemyEffects string
	var asciiText string
	var seedStr string
	var debugMode bool

	cmd := &cobra.Command{
		Use:   "hera",
		Short: "Hera - A tactical turn-based game",
		Long:  "A tactical turn-based game. Made with ♡",
		RunE: func(cmd *cobra.Command, _ []string) error {
			registry := generate.InitDefaultTheme()
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
			centerWindow := centerMode
			enableBackground := background
			themeName := theme
			if themeName == "" {
				themeName = "default"
			}
			pc := flagPlayers
			if pc < 0 {
				pc = 0
			}
			ec := flagEnemies
			if ec < 0 {
				ec = 0
			}
			model := generate.NewModel(pc, ec, parseEffects(flagPlayerEffects), parseEffects(flagEnemyEffects), loc, registry, centerWindow, enableBackground, themeName)

			if asciiText != "" {
				font, ferr := generate.ParseFLF(assets.FontData)
				if ferr == nil {
					model.BannerText = font.Render(asciiText)
				}
			}

			if err := model.InitConfig(); err != nil {
				log.Warn("config", "err", err)
			}
			if model.Config != nil {
				if theme != "" {
					model.Config.ThemeName = theme
				}
				if cmd.Flags().Changed("background") {
					model.Config.Background = background
				}
				if cmd.Flags().Changed("center") {
					model.Config.CenterWindow = centerMode
				}
				if model.Config.Language != "" && !cmd.Flags().Changed("lang") {
					if err := model.SetLanguage(model.Config.Language); err != nil {
						log.Warn("language", "err", err)
					}
				}
				registry.SetTintID(model.Config.ThemeName)
				model.Theme = registry
				model.Styles = generate.NewStyles(model.Theme)
			}
			model.SetAvailableThemes()

			if err := model.LoadProfiles(); err != nil {
				log.Warn("profiles", "err", err)
			}

			if model.Font == nil {
				font, ferr := generate.ParseFLF(assets.FontData)
				if ferr == nil {
					model.Font = font
				}
			}

			if seedStr != "" {
				model.Seed = generate.ParseSeed(seedStr)
				model.SeedLocked = true
			}

			model.DebugMode = debugMode

			p := tea.NewProgram(
				&model,
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
	cmd.Flags().StringVarP(&centerMode, "center", "C", "", "Window position: tl, tc, tr, cl, c, cr, bl, bc, br (default: c)")
	cmd.Flags().BoolVarP(&background, "background", "b", false, "Background fill")
	cmd.Flags().IntVarP(&flagPlayers, "players", "p", 0, "Number of players (0=random, 2-4)")
	cmd.Flags().IntVarP(&flagEnemies, "enemies", "e", 0, "Number of enemies (0=random)")
	cmd.Flags().StringVarP(&flagPlayerEffects, "player-effects", "P", "", "Starting effects for players (comma-separated: fire,wet,smoke)")
	cmd.Flags().StringVarP(&flagEnemyEffects, "enemy-effects", "E", "", "Starting effects for enemies (comma-separated: fire,wet,smoke)")
	cmd.Flags().StringVarP(&asciiText, "ascii", "a", "", "Render text as ASCII banner using assets/Tmplr.flf")
	cmd.Flags().StringVarP(&seedStr, "seed", "s", "", "Seed for deterministic world generation")
	cmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug mode (console, no score saving)")

	opts := []fang.Option{fang.WithVersion(version)}
	if commit != "" {
		opts = append(opts, fang.WithCommit(commit))
	}
	if err := fang.Execute(context.Background(), cmd, opts...); err != nil {
		log.Error(err)
	}
}
