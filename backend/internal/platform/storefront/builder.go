package storefront

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/platform/storefront/catalog"
	"smarterp/backend/internal/platform/storefront/themes/mamapack"
)

const builderPageSize = 36

func (sh *shop) builder(r *http.Request) (builderVM, error) {
	page, selected := pageParam(r), r.URL.Query().Get("kategoria")
	filter := catalog.Filter{Kind: "builder", Collection: selected, Limit: builderPageSize + 1, Offset: (page - 1) * builderPageSize}
	products, err := sh.catalog(filter)
	if err != nil {
		return builderVM{}, err
	}
	vm := builderVM{Layout: sh.pageLayout(r, "Złóż własną wyprawkę do szpitala"), Selected: selected}
	if len(products) > builderPageSize {
		vm.Next = builderPageURL(r.URL.Path, page+1, selected)
		products = products[:builderPageSize]
	}
	if page > 1 {
		vm.Prev = builderPageURL(r.URL.Path, page-1, selected)
	}
	vm.Groups = mamapack.BuilderGroups(sh.cards(products), selected)
	vm.Categories, err = sh.builderCategories()
	return vm, err
}

func builderPageURL(path string, page int, category string) string {
	query := url.Values{}
	if page > 1 {
		query.Set("page", strconv.Itoa(page))
	}
	if category != "" {
		query.Set("kategoria", category)
	}
	return (&url.URL{Path: path, RawQuery: query.Encode()}).String()
}

func (sh *shop) builderCategories() ([]catalog.Category, error) {
	const query = `SELECT c.slug,c.name FROM catalog.categories c
		WHERE c.tenant_id=$1 AND (c.is_published OR $2) AND c.slug<>'' AND EXISTS (
			SELECT 1 FROM catalog.product_categories pc JOIN catalog.products p ON p.id=pc.product_id AND p.tenant_id=pc.tenant_id
			JOIN catalog.product_variants v ON v.product_id=p.id AND v.tenant_id=p.tenant_id AND (v.is_published OR $2)
			WHERE pc.tenant_id=c.tenant_id AND pc.category_id=c.id AND (p.is_published OR $2) AND NOT p.is_composite)
		ORDER BY c.sort_order,c.name`
	rows, err := sh.service.repo.store.Query(sh.ctx, sh.tenantID, query, sh.tenantID, sh.preview)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (catalog.Category, error) {
		var c catalog.Category
		err := row.Scan(&c.Slug, &c.Name)
		return c, err
	})
}
