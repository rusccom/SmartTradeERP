package storefront

import "embed"

//go:embed themes/*/*.json themes/*/templates themes/*/assets cart variants shared
var themesFS embed.FS
