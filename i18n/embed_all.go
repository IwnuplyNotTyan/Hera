//go:build !eng

package i18n

import "embed"

var (
	//go:embed locales/*.json
	embeddedLocales embed.FS
)
