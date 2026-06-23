package storefront

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type brandRow struct {
	TenantName string
	ThemeID    string
	TokensJSON []byte
	Sections   []byte
	LogoKey    string
	Currency   currency
}

type productCardRow struct {
	Name     string
	Slug     string
	ImageKey string
	MinPrice decimal.Decimal
	HasPrice bool
}

type productDetailRow struct {
	ID             string
	Name           string
	Slug           string
	SEOTitle       string
	SEODescription string
}

type variantRow struct {
	ID       string
	Name     string
	SKU      string
	Price    decimal.Decimal
	HasPrice bool
}

const brandPublishedSQL = `SELECT t.name,
        COALESCE(ts.storefront_theme_id, 'classic'),
        COALESCE(ts.published_tokens, '{}'::jsonb),
        COALESCE(ts.sections, '[]'::jsonb),
        COALESCE(m.object_key, ''),
        COALESCE(c.code, ''), COALESCE(c.symbol, ''), COALESCE(c.decimal_places, 2)
    FROM platform.tenants t
    LEFT JOIN platform.tenant_settings ts ON ts.tenant_id = t.id
    LEFT JOIN platform.currencies c ON c.id = ts.base_currency_id
    LEFT JOIN platform.media_objects m ON m.id = ts.logo_media_id AND m.tenant_id = t.id AND m.status = 'ready'
    WHERE t.id = $1`

const brandPreviewSQL = `SELECT t.name,
        COALESCE(NULLIF(ts.draft_theme_id, ''), ts.storefront_theme_id, 'classic'),
        COALESCE(ts.draft_tokens, '{}'::jsonb),
        COALESCE(ts.draft_sections, '[]'::jsonb),
        COALESCE(m.object_key, ''),
        COALESCE(c.code, ''), COALESCE(c.symbol, ''), COALESCE(c.decimal_places, 2)
    FROM platform.tenants t
    LEFT JOIN platform.tenant_settings ts ON ts.tenant_id = t.id
    LEFT JOIN platform.currencies c ON c.id = ts.base_currency_id
    LEFT JOIN platform.media_objects m ON m.id = ts.logo_media_id AND m.tenant_id = t.id AND m.status = 'ready'
    WHERE t.id = $1`

// LoadBrand reads the tenant brand, theme, tokens and sections. In preview mode
// it reads the unpublished draft columns so an owner can review changes.
func (r *Repository) LoadBrand(ctx context.Context, tenantID string, preview bool) (brandRow, error) {
	query := brandPublishedSQL
	if preview {
		query = brandPreviewSQL
	}
	out := brandRow{}
	var code, symbol string
	var places int
	err := r.store.Pool.QueryRow(ctx, query, tenantID).Scan(
		&out.TenantName, &out.ThemeID, &out.TokensJSON, &out.Sections, &out.LogoKey, &code, &symbol, &places)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return brandRow{}, ErrTenantNotFound
		}
		return brandRow{}, err
	}
	out.Currency = currency{code: code, symbol: symbol, places: int32(places)}
	return out, nil
}

const listProductsSQL = `SELECT p.name, p.slug,
        COALESCE(pi.object_key, ''),
        COALESCE(pr.min_price, 0),
        (pr.min_price IS NOT NULL)
    FROM catalog.products p
    LEFT JOIN LATERAL (
        SELECT m.object_key FROM platform.media_objects m
        WHERE m.tenant_id = p.tenant_id AND m.owner_type = 'product' AND m.owner_id = p.id AND m.status = 'ready'
        ORDER BY m.is_primary DESC, m.sort_order, m.created_at LIMIT 1
    ) pi ON true
    LEFT JOIN LATERAL (
        SELECT MIN(v.price) AS min_price FROM catalog.product_variants v
        WHERE v.tenant_id = p.tenant_id AND v.product_id = p.id AND v.is_published AND v.price IS NOT NULL
    ) pr ON true
    WHERE p.tenant_id = $1 AND p.is_published
    ORDER BY p.created_at DESC
    LIMIT $2 OFFSET $3`

func (r *Repository) ListPublishedProducts(ctx context.Context, tenantID string, limit, offset int) ([]productCardRow, error) {
	rows, err := r.store.Pool.Query(ctx, listProductsSQL, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCards(rows)
}

func scanCards(rows pgx.Rows) ([]productCardRow, error) {
	cards := make([]productCardRow, 0)
	for rows.Next() {
		card := productCardRow{}
		if err := rows.Scan(&card.Name, &card.Slug, &card.ImageKey, &card.MinPrice, &card.HasPrice); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func (r *Repository) CountPublishedProducts(ctx context.Context, tenantID string) (int, error) {
	const query = `SELECT COUNT(*) FROM catalog.products WHERE tenant_id = $1 AND is_published`
	total := 0
	err := r.store.Pool.QueryRow(ctx, query, tenantID).Scan(&total)
	return total, err
}

const productBySlugSQL = `SELECT p.id::text, p.name, p.slug, COALESCE(p.seo_title, ''), COALESCE(p.seo_description, '')
    FROM catalog.products p
    WHERE p.tenant_id = $1 AND p.slug = $2 AND p.is_published
    LIMIT 1`

func (r *Repository) ProductBySlug(ctx context.Context, tenantID, slug string) (productDetailRow, error) {
	out := productDetailRow{}
	err := r.store.Pool.QueryRow(ctx, productBySlugSQL, tenantID, slug).Scan(
		&out.ID, &out.Name, &out.Slug, &out.SEOTitle, &out.SEODescription)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDetailRow{}, ErrProductNotFound
		}
		return productDetailRow{}, err
	}
	return out, nil
}

const variantsSQL = `SELECT v.id::text, COALESCE(v.name, 'Default'), COALESCE(v.sku_code, ''), COALESCE(v.price, 0), (v.price IS NOT NULL)
    FROM catalog.product_variants v
    WHERE v.tenant_id = $1 AND v.product_id = $2 AND v.is_published
    ORDER BY v.created_at`

func (r *Repository) VariantsForProduct(ctx context.Context, tenantID, productID string) ([]variantRow, error) {
	rows, err := r.store.Pool.Query(ctx, variantsSQL, tenantID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVariants(rows)
}

func scanVariants(rows pgx.Rows) ([]variantRow, error) {
	variants := make([]variantRow, 0)
	for rows.Next() {
		v := variantRow{}
		if err := rows.Scan(&v.ID, &v.Name, &v.SKU, &v.Price, &v.HasPrice); err != nil {
			return nil, err
		}
		variants = append(variants, v)
	}
	return variants, rows.Err()
}

const productImagesSQL = `SELECT m.object_key FROM platform.media_objects m
    WHERE m.tenant_id = $1 AND m.owner_type = 'product' AND m.owner_id = $2 AND m.status = 'ready'
    ORDER BY m.is_primary DESC, m.sort_order, m.created_at`

func (r *Repository) ProductImageKeys(ctx context.Context, tenantID, productID string) ([]string, error) {
	rows, err := r.store.Pool.Query(ctx, productImagesSQL, tenantID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStrings(rows)
}

const slugsSQL = `SELECT slug FROM catalog.products
    WHERE tenant_id = $1 AND is_published AND slug <> ''
    ORDER BY updated_at DESC LIMIT $2`

func (r *Repository) PublishedProductSlugs(ctx context.Context, tenantID string, limit int) ([]string, error) {
	rows, err := r.store.Pool.Query(ctx, slugsSQL, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStrings(rows)
}

func scanStrings(rows pgx.Rows) ([]string, error) {
	values := make([]string, 0)
	for rows.Next() {
		value := ""
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, rows.Err()
}
