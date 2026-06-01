//go:build notint

package generate

func InitDefaultTheme() ThemeRegistry {
	return newNoTintRegistry()
}
