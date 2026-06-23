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
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	PreviewURL string            `json:"preview_url"`
	Tokens     map[string]string `json:"tokens"`
}

type themeEntry struct {
	manifest     Manifest
	pages        map[string]*template.Template
	css          []byte
	assetVersion string
}

// requiredPages every theme must define. validateTheme fails at boot otherwise,
// so a broken theme can never ship a blank shop.
var requiredPages = []string{"home", "list", "product", "cart", "404"}

// NewRegistry loads and validates every embedded theme folder.
func NewRegistry() (*Registry, error) {
	entries, err := fs.ReadDir(themesFS, "themes")
	if err != nil {
		return nil, err
	}
	reg := &Registry{themes: make(map[string]*themeEntry)}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		loaded, err := loadTheme(entry.Name())
		if err != nil {
			return nil, fmt.Errorf("storefront theme %q: %w", entry.Name(), err)
		}
		reg.themes[entry.Name()] = loaded
	}
	if _, ok := reg.themes["classic"]; !ok {
		return nil, errors.New("storefront: classic theme is required")
	}
	return reg, nil
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
	pages, err := loadPages(id)
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

func loadPages(id string) (map[string]*template.Template, error) {
	pages := make(map[string]*template.Template, len(requiredPages))
	for _, name := range requiredPages {
		tmpl, err := parsePage(id, name)
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
		path.Join(dir, "layout.html"),
		path.Join(dir, "partials", "*.html"),
		path.Join(dir, name+".html"),
	)
}

func hashVersion(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])[:8]
}
