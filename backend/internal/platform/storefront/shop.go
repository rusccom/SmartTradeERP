package storefront

import (
	"context"
	"net/http"

	"smarterp/backend/internal/platform/storefront/catalog"
	"smarterp/backend/internal/platform/storefront/themes/mamapack"
	"smarterp/backend/internal/shared/tenant"
)

type shop struct {
	service       *Service
	ctx           context.Context
	tenantID      string
	brand         brandRow
	preview, demo bool
}

func (s *Service) open(r *http.Request, preview bool) (*shop, error) {
	sh := &shop{service: s, ctx: r.Context(), tenantID: tenant.FromContext(r.Context()), preview: preview}
	if s.repo == nil {
		sh.brand = brandRow{TenantName: "MamaPack", ThemeID: "mamapack", Currency: currency{code: "PLN", symbol: "zł", places: 2}}
		sh.preview, sh.demo = true, true
		return sh, nil
	}
	brand, err := s.repo.LoadBrand(sh.ctx, sh.tenantID, preview)
	if err != nil {
		return nil, err
	}
	sh.brand = brand
	if preview && brand.ThemeID == "mamapack" {
		err = sh.selectSample()
	}
	return sh, err
}

func (sh *shop) catalog(f catalog.Filter) ([]catalog.Product, error) {
	if sh.demo {
		return mamapack.SampleProducts(f), nil
	}
	return sh.service.repo.catalog(sh.ctx, sh.tenantID, f)
}

func (sh *shop) layout(r *http.Request, seo SeoVM) LayoutVM {
	layout := sh.service.buildLayout(r, sh.brand, seo, sh.preview)
	layout.Demo = sh.demo
	layout.CartScope = sh.tenantID
	if sh.demo {
		layout.CartScope += "/sample"
	}
	return layout
}

func shopLanguage(brand brandRow) string {
	if brand.Currency.code == "PLN" {
		return "pl"
	}
	return "en"
}

func (sh *shop) images(p catalog.Product, transform string) []string {
	if sh.demo {
		return p.Images
	}
	return sh.service.imageList(p.Images, transform)
}

func (sh *shop) selectSample() error {
	const query = "SELECT NOT EXISTS (SELECT 1 FROM catalog.products WHERE tenant_id=$1 AND is_published)"
	err := sh.service.repo.store.QueryRow(sh.ctx, sh.tenantID, query, sh.tenantID).Scan(&sh.demo)
	if sh.demo {
		sh.brand.Currency = currency{code: "PLN", symbol: "zł", places: 2}
	}
	return err
}
