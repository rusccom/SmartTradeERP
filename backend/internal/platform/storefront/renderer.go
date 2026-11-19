package storefront

import (
	"bytes"
	"html/template"
	"net/http"
)

// Registry holds every embedded theme parsed and validated at boot.
type Registry struct {
	defaultID string
	cart      *template.Template
	themes    map[string]*themeEntry
}

func (rg *Registry) Has(themeID string) bool {
	_, ok := rg.themes[themeID]
	return ok
}

func (rg *Registry) AssetVersion(themeID string) string {
	return rg.theme(themeID).assetVersion
}

func (rg *Registry) CSS(themeID string) ([]byte, string, bool) {
	entry, ok := rg.themes[themeID]
	if !ok {
		return nil, "", false
	}
	return entry.css, entry.assetVersion, true
}

// Render executes a theme page through its layout and writes the buffered HTML,
// so a render error never produces a half-written response.
func (rg *Registry) Render(w http.ResponseWriter, status int, themeID, page string, data any) {
	tmpl := rg.theme(themeID).pages[page]
	if page == "cart" {
		tmpl = rg.cart
	}
	if tmpl == nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout", data); err != nil {
		http.Error(w, "render failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}

// theme returns the requested theme or falls back to the installed default so an unknown id
// never produces a blank page.
func (rg *Registry) theme(themeID string) *themeEntry {
	if entry, ok := rg.themes[themeID]; ok {
		return entry
	}
	return rg.themes[rg.defaultID]
}
