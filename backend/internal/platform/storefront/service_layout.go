package storefront

import (
	"net/http"
)

// buildLayout supplies the active theme, store data and SEO for this request.
func (s *Service) buildLayout(r *http.Request, brand brandRow, seo SeoVM, preview bool) LayoutVM {
	themeID := s.registry.theme(brand.ThemeID).manifest.ID
	if preview {
		seo.Robots = "noindex,nofollow"
	}
	return LayoutVM{
		Brand:    BrandVM{TenantName: brand.TenantName, LogoURL: s.image(brand.LogoKey, logoTransform), HomeURL: "/"},
		Theme:    ThemeVM{ID: themeID, AssetVersion: s.registry.AssetVersion(themeID)},
		SEO:      seo,
		Plugins:  brand.Plugins,
		Nav:      defaultNav(),
		Preview:  preview,
		Language: shopLanguage(brand),
	}
}

func (s *Service) image(key, transform string) string {
	return transformedImageURL(s.opts.MediaBaseURL, key, transform)
}

func defaultNav() []NavItemVM {
	return []NavItemVM{{Label: "Home", URL: "/"}, {Label: "Products", URL: "/products"}}
}
