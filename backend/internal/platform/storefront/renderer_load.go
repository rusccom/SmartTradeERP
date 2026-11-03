package storefront

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"path"
)

type Manifest struct {
	Collections map[string]string `json:"collections"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	PreviewURL  string            `json:"preview_url"`
	Pages       map[string]string `json:"pages"`
	Tokens      map[string]string `json:"tokens"`
}

type themeEntry struct {
	manifest     Manifest
	pages        map[string]*template.Template
	css          []byte
	assetVersion string
}

// requiredPages every theme must define. validateTheme fails at boot otherwise,
// so a broken theme can never ship a blank shop.
var requiredPages = []string{"home", "list", "product", "404"}

// NewRegistry loads and validates every embedded theme folder.
func NewRegistry() (*Registry, error) {
	entries, err := fs.ReadDir(themesFS, "themes")
	if err != nil {
		return nil, err
	}
	cart, err := template.ParseFS(themesFS, "shared/cart.html", "shared/seo-head.html")
	if err != nil {
		return nil, err
	}
	reg := &Registry{cart: cart, themes: make(map[string]*themeEntry)}
	if err := reg.loadThemes(entries); err != nil {
		return nil, err
	}
	if !reg.Has("classic") {
		return nil, errors.New("storefront: classic theme is required")
	}
	return reg, nil
}

func (reg *Registry) loadThemes(entries []fs.DirEntry) error {
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		loaded, err := loadTheme(entry.Name())
		if err != nil {
			return fmt.Errorf("storefront theme %q: %w", entry.Name(), err)
		}
		reg.themes[entry.Name()] = loaded
	}
	return nil
}

func loadTheme(id string) (*themeEntry, error) {
	manifest, err := loadManifest(id)
	if err != nil {
		return nil, err
	}
	css, err := fs.ReadFile(themesFS, path.Join("themes", id, "assets", "theme.css"))
	if err != nil {
		return nil, err
	}
	pages, err := loadPages(manifest)
	if err != nil {
		return nil, err
	}
	return &themeEntry{manifest: manifest, pages: pages, css: css, assetVersion: hashVersion(css)}, nil
}

func loadManifest(id string) (Manifest, error) {
	raw, err := fs.ReadFile(themesFS, path.Join("themes", id, "theme.json"))
	if err != nil {
		return Manifest{}, err
	}
	manifest := Manifest{}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func loadPages(manifest Manifest) (map[string]*template.Template, error) {
	pages := make(map[string]*template.Template, len(requiredPages))
	names := append([]string(nil), requiredPages...)
	for _, name := range manifest.Pages {
		names = append(names, name)
	}
	for _, name := range names {
		tmpl, err := parsePage(manifest.ID, name)
		if err != nil {
			return nil, err
		}
		pages[name] = tmpl
	}
	return pages, nil
}

func parsePage(id, name string) (*template.Template, error) {
	dir := path.Join("themes", id, "templates")
	return template.New(name).ParseFS(themesFS,
		"shared/seo-head.html", "shared/variant-picker.html", path.Join(dir, "layout.html"),
		path.Join(dir, "partials", "*.html"),
		path.Join(dir, name+".html"),
	)
}

func hashVersion(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])[:8]
}
