package storefront

import (
	"context"
	"net/http"

	"smarterp/backend/internal/platform/storefront/catalog"
	"smarterp/backend/internal/shared/tenant"
)

type shop struct {
	service  *Service
	ctx      context.Context
	tenantID string
	brand    brandRow
	preview  bool
}

func (s *Service) open(r *http.Request, preview bool) (*shop, error) {
	sh := &shop{service: s, ctx: r.Context(), tenantID: tenant.FromContext(r.Context()), preview: preview}
	brand, err := s.repo.LoadBrand(sh.ctx, sh.tenantID, preview)
	sh.brand = brand
	return sh, err
}

func (sh *shop) catalog(f catalog.Filter) ([]catalog.Product, error) {
	f.Preview = sh.preview
	return sh.service.repo.catalog(sh.ctx, sh.tenantID, f)
}

func (sh *shop) layout(r *http.Request, seo SeoVM) LayoutVM {
	return sh.service.buildLayout(r, sh.brand, seo, sh.preview)
}

func shopLanguage(brand brandRow) string {
	if brand.Currency.code == "PLN" {
		return "pl"
	}
	return "en"
}
