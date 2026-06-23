package storefront

import (
	"encoding/json"
	"html/template"
	"net/http"
)

// buildLayout assembles the shared chrome: brand, resolved theme tokens, and
// SEO. An unknown theme id falls back to classic so a page never renders blank.
// Preview pages are marked noindex so a draft is never indexed.
func (s *Service) buildLayout(r *http.Request, brand brandRow, seo SeoVM, preview bool) LayoutVM {
	themeID := brand.ThemeID
	if !s.registry.Has(themeID) {
		themeID = "classic"
	}
	if preview {
		seo.Robots = "noindex,nofollow"
	}
	return LayoutVM{
		Brand: BrandVM{TenantName: brand.TenantName, LogoURL: s.image(brand.LogoKey, logoTransform), HomeURL: "/"},
		Theme: ThemeVM{ID: themeID, AssetVersion: s.registry.AssetVersion(themeID), TokenCSS: s.tokenCSS(themeID, brand.TokensJSON)},
		SEO:   seo,
		Nav:   defaultNav(),
	}
}

func (s *Service) tokenCSS(themeID string, tokensJSON []byte) template.CSS {
	return buildTokenCSS(s.registry.DefaultTokens(themeID), parseTokens(tokensJSON))
}

func (s *Service) image(key, transform string) string {
	return transformedImageURL(s.opts.MediaBaseURL, key, transform)
}

func defaultNav() []NavItemVM {
	return []NavItemVM{{Label: "Home", URL: "/"}, {Label: "Products", URL: "/products"}}
}

func parseTokens(raw []byte) map[string]string {
	out := map[string]string{}
	if len(raw) == 0 {
		return out
	}
	_ = json.Unmarshal(raw, &out)
	return out
}
