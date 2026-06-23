package storefront

import (
	_ "embed"
	"net/http"
)

// storefrontJS is the shared progressive-enhancement cart script, identical
// across themes. Themes provide the markup hooks (data-sf-*); this drives them.
//
//go:embed assets/storefront.js
var storefrontJS []byte

func (h *Handler) AppJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = w.Write(storefrontJS)
}
