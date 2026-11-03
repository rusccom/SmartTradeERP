package storefront

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type brandRow struct {
	TenantName string
	ThemeID    string
	TokensJSON []byte
	Sections   []byte
	LogoKey    string
	Currency   currency
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
	err := r.store.QueryRow(ctx, tenantID, query, tenantID).Scan(
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

const slugsSQL = `SELECT slug FROM catalog.products
    WHERE is_published AND slug <> ''
    ORDER BY updated_at DESC LIMIT $1`

func (r *Repository) PublishedProductSlugs(ctx context.Context, tenantID string, limit int) ([]string, error) {
	rows, err := r.store.Query(ctx, tenantID, slugsSQL, limit)
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
