//go:build !notint

package generate

import tint "github.com/lrstanley/bubbletint"

func InitDefaultTheme() ThemeRegistry {
	tint.NewDefaultRegistry()
	RegisterDefaults()
	return tint.DefaultRegistry
}
