package storefront

import (
	"context"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"smarterp/backend/internal/shared/tenant"
)

const (
	homeProductLimit = 12
	listPageSize     = 24
)

func (s *Service) BuildHome(r *http.Request, preview bool) (HomeVM, error) {
	ctx := r.Context()
	tenantID := tenant.FromContext(ctx)
	brand, err := s.repo.LoadBrand(ctx, tenantID, preview)
	if err != nil {
		return HomeVM{}, err
	}
	rows, err := s.repo.ListPublishedProducts(ctx, tenantID, homeProductLimit, 0)
	if err != nil {
		return HomeVM{}, err
	}
	layout := s.buildLayout(r, brand, seoHome(r, brand), preview)
	return HomeVM{Layout: layout, Products: s.toCards(brand, rows), Sections: homeSections(brand.Sections)}, nil
}

func (s *Service) BuildList(r *http.Request, preview bool) (ListVM, error) {
	ctx := r.Context()
	tenantID := tenant.FromContext(ctx)
	brand, err := s.repo.LoadBrand(ctx, tenantID, preview)
	if err != nil {
		return ListVM{}, err
	}
	page := pageParam(r)
	rows, err := s.repo.ListPublishedProducts(ctx, tenantID, listPageSize, (page-1)*listPageSize)
	if err != nil {
		return ListVM{}, err
	}
	total, err := s.repo.CountPublishedProducts(ctx, tenantID)
	if err != nil {
		return ListVM{}, err
	}
	return s.assembleList(r, brand, rows, page, total, preview), nil
}

func (s *Service) assembleList(r *http.Request, brand brandRow, rows []productCardRow, page, total int, preview bool) ListVM {
	offset := (page - 1) * listPageSize
	return ListVM{
		Layout:   s.buildLayout(r, brand, seoList(r, brand), preview),
		Products: s.toCards(brand, rows),
		Page:     page,
		HasPrev:  page > 1,
		HasNext:  offset+len(rows) < total,
		PrevURL:  listPageURL(page - 1),
		NextURL:  listPageURL(page + 1),
	}
}

func (s *Service) BuildProduct(r *http.Request, slug string, preview bool) (ProductDetailVM, error) {
	ctx := r.Context()
	tenantID := tenant.FromContext(ctx)
	brand, err := s.repo.LoadBrand(ctx, tenantID, preview)
	if err != nil {
		return ProductDetailVM{}, err
	}
	detail, err := s.repo.ProductBySlug(ctx, tenantID, slug)
	if err != nil {
		return ProductDetailVM{}, err
	}
	return s.assembleProduct(ctx, r, brand, detail, preview)
}

func (s *Service) assembleProduct(ctx context.Context, r *http.Request, brand brandRow, detail productDetailRow, preview bool) (ProductDetailVM, error) {
	tenantID := tenant.FromContext(ctx)
	variants, err := s.repo.VariantsForProduct(ctx, tenantID, detail.ID)
	if err != nil {
		return ProductDetailVM{}, err
	}
	keys, err := s.repo.ProductImageKeys(ctx, tenantID, detail.ID)
	if err != nil {
		return ProductDetailVM{}, err
	}
	images := s.imageList(keys, mainTransform)
	price, hasPrice := minVariantPrice(variants, brand.Currency)
	seo := seoProduct(r, detail, firstString(images), price, hasPrice)
	return ProductDetailVM{
		Layout: s.buildLayout(r, brand, seo, preview), Name: detail.Name, Description: detail.SEODescription,
		Images: images, Price: price, HasPrice: hasPrice, Variants: toVariantVMs(variants, brand.Currency),
	}, nil
}

func (s *Service) BuildCart(r *http.Request, preview bool) (CartVM, error) {
	ctx := r.Context()
	brand, err := s.repo.LoadBrand(ctx, tenant.FromContext(ctx), preview)
	if err != nil {
		return CartVM{}, err
	}
	return CartVM{Layout: s.buildLayout(r, brand, seoCart(r, brand), preview)}, nil
}

func (s *Service) BuildNotFound(r *http.Request, preview bool) (NotFoundVM, error) {
	ctx := r.Context()
	brand, err := s.repo.LoadBrand(ctx, tenant.FromContext(ctx), preview)
	if err != nil {
		return NotFoundVM{}, err
	}
	return NotFoundVM{Layout: s.buildLayout(r, brand, seoNotFound(r, brand), preview)}, nil
}

func (s *Service) toCards(brand brandRow, rows []productCardRow) []ProductCardVM {
	cards := make([]ProductCardVM, 0, len(rows))
	for _, row := range rows {
		cards = append(cards, s.toCard(brand, row))
	}
	return cards
}

func (s *Service) toCard(brand brandRow, row productCardRow) ProductCardVM {
	return ProductCardVM{
		Name:     row.Name,
		URL:      productURL(row.Slug),
		ImageURL: s.image(row.ImageKey, cardTransform),
		Price:    formatMoney(row.MinPrice, brand.Currency),
		HasPrice: row.HasPrice,
	}
}

func (s *Service) imageList(keys []string, transform string) []string {
	urls := make([]string, 0, len(keys))
	for _, key := range keys {
		if url := s.image(key, transform); url != "" {
			urls = append(urls, url)
		}
	}
	return urls
}

func minVariantPrice(variants []variantRow, cur currency) (MoneyVM, bool) {
	min := decimal.Zero
	found := false
	for _, v := range variants {
		if v.HasPrice && (!found || v.Price.LessThan(min)) {
			min = v.Price
			found = true
		}
	}
	if !found {
		return MoneyVM{}, false
	}
	return formatMoney(min, cur), true
}

func toVariantVMs(variants []variantRow, cur currency) []ProductVariantVM {
	vms := make([]ProductVariantVM, 0, len(variants))
	for _, v := range variants {
		vms = append(vms, ProductVariantVM{VariantID: v.ID, Name: v.Name, SKU: v.SKU, Price: formatMoney(v.Price, cur), HasPrice: v.HasPrice})
	}
	return vms
}

func firstString(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func productURL(slug string) string {
	if slug == "" {
		return "/products"
	}
	return "/products/" + slug
}

func listPageURL(page int) string {
	if page <= 1 {
		return "/products"
	}
	return "/products?page=" + strconv.Itoa(page)
}

func pageParam(r *http.Request) int {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}
