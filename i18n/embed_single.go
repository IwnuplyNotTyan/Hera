//go:build eng

package i18n

import "embed"

var (
	//go:embed locales/en.json
	embeddedLocales embed.FS
)
