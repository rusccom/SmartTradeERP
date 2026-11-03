package storefront

import "embed"

// themesFS embeds every theme folder (markup + manifest + assets). Adding a new
// theme is dropping a folder under themes/ — no Go change is required.
//
//go:embed themes/*/theme.json themes/*/templates themes/*/assets shared assets
var themesFS embed.FS
