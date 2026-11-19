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

	"smarterp/backend/internal/platform/storefront/themes"
)

type Manifest struct {
	Collections map[string]string         `json:"collections"`
	ID          string                    `json:"-"`
	Pages       map[string]PageDefinition `json:"pages"`
}

type PageDefinition struct {
	Template string      `json:"template"`
	Title    string      `json:"title"`
	Plugin   string      `json:"plugin"`
	Catalog  PageCatalog `json:"catalog"`
}

type PageCatalog struct {
	ProductType     string `json:"product_type"`
	Collection      string `json:"collection"`
	CategoryParam   string `json:"category_param"`
	RequireVariants bool   `json:"require_variants"`
	PageSize        int    `json:"page_size"`
}

type themeEntry struct {
	manifest     Manifest
	pages        map[string]*template.Template
	css          []byte
	assetVersion string
}

func NewRegistry() (*Registry, error) {
	idx, err := themes.Load()
	if err != nil {
		return nil, err
	}
	cart, err := template.ParseFS(themesFS, "cart/cart.html", "shared/seo-head.html")
	if err != nil {
		return nil, err
	}
	r := &Registry{cart: cart, defaultID: idx.Default, themes: map[string]*themeEntry{}}
	for _, info := range idx.Items {
		t, err := loadTheme(info.ID)
		if err != nil {
			return nil, fmt.Errorf("storefront theme %q: %w", info.ID, err)
		}
		r.themes[info.ID] = t
	}
	return r, nil
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
	m := Manifest{ID: id}
	raw, err := themesFS.ReadFile(path.Join("themes", id, "routes.json"))
	if errors.Is(err, fs.ErrNotExist) {
		return m, nil
	}
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(raw, &m)
	return m, err
}

func loadPages(m Manifest) (map[string]*template.Template, error) {
	names := []string{"home", "list", "product", "404"}
	for _, p := range m.Pages {
		names = append(names, p.Template)
	}
	pages := make(map[string]*template.Template, len(names))
	for _, name := range names {
		t, err := parsePage(m.ID, name)
		if err != nil {
			return nil, err
		}
		pages[name] = t
	}
	return pages, nil
}

func parsePage(id, name string) (*template.Template, error) {
	dir := path.Join("themes", id, "templates")
	partials, err := fs.Glob(themesFS, path.Join(dir, "partials", "*.html"))
	if err != nil {
		return nil, err
	}
	files := []string{"shared/seo-head.html", "variants/variant-picker.html", path.Join(dir, "layout.html")}
	files = append(files, partials...)
	return template.New(name).Funcs(pageHelpers).ParseFS(themesFS, append(files, path.Join(dir, name+".html"))...)
}

func hashVersion(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])[:8]
}
