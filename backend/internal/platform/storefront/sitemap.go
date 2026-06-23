package storefront

import (
	"encoding/xml"
	"net/http"

	"smarterp/backend/internal/shared/tenant"
)

const sitemapLimit = 5000

type sitemapURL struct {
	Loc string `xml:"loc"`
}

type sitemapSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

func (s *Service) BuildSitemap(r *http.Request) ([]byte, error) {
	ctx := r.Context()
	slugs, err := s.repo.PublishedProductSlugs(ctx, tenant.FromContext(ctx), sitemapLimit)
	if err != nil {
		return nil, err
	}
	return renderSitemap(r, slugs)
}

func renderSitemap(r *http.Request, slugs []string) ([]byte, error) {
	set := sitemapSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	set.URLs = append(set.URLs, sitemapURL{Loc: absoluteURL(r, "/")})
	for _, slug := range slugs {
		set.URLs = append(set.URLs, sitemapURL{Loc: absoluteURL(r, productURL(slug))})
	}
	body, err := xml.Marshal(set)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), body...), nil
}
