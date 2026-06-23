package storefront

import (
	"bytes"
	"net/http"
	"sort"
)

// Registry holds every embedded theme parsed and validated at boot.
type Registry struct {
	themes map[string]*themeEntry
}

// ThemeInfo is the admin-facing description of an installed theme: its id, its
// display name, a preview URL, and its default tokens (the customizer form).
type ThemeInfo struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	PreviewURL string            `json:"preview_url"`
	Tokens     map[string]string `json:"tokens"`
}

// Themes lists installed themes sorted by id, for the admin picker.
func (rg *Registry) Themes() []ThemeInfo {
	infos := make([]ThemeInfo, 0, len(rg.themes))
	for id, entry := range rg.themes {
		infos = append(infos, ThemeInfo{ID: id, Name: entry.manifest.Name, PreviewURL: entry.manifest.PreviewURL, Tokens: entry.manifest.Tokens})
	}
	sort.Slice(infos, func(i, j int) bool { return infos[i].ID < infos[j].ID })
	return infos
}

func (rg *Registry) Has(themeID string) bool {
	_, ok := rg.themes[themeID]
	return ok
}

func (rg *Registry) AssetVersion(themeID string) string {
	return rg.theme(themeID).assetVersion
}

func (rg *Registry) DefaultTokens(themeID string) map[string]string {
	return rg.theme(themeID).manifest.Tokens
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

// theme returns the requested theme or falls back to classic so an unknown id
// never produces a blank page.
func (rg *Registry) theme(themeID string) *themeEntry {
	if entry, ok := rg.themes[themeID]; ok {
		return entry
	}
	return rg.themes["classic"]
}
