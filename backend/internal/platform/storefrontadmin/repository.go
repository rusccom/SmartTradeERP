package storefrontadmin

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
)

type settingsRow struct {
	ThemeID         string
	DraftThemeID    string
	PublishedTokens []byte
	DraftTokens     []byte
	Sections        []byte
	DraftSections   []byte
}

type Repository struct {
	store *db.Store
}

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

const loadSQL = `SELECT COALESCE(storefront_theme_id, 'classic'),
        COALESCE(draft_theme_id, ''),
        COALESCE(published_tokens, '{}'::jsonb),
        COALESCE(draft_tokens, '{}'::jsonb),
        COALESCE(sections, '[]'::jsonb),
        COALESCE(draft_sections, '[]'::jsonb)
    FROM platform.tenant_settings`

func (r *Repository) Load(ctx context.Context, tenantID string) (settingsRow, error) {
	out := settingsRow{}
	err := r.store.QueryRow(ctx, tenantID, loadSQL).Scan(
		&out.ThemeID, &out.DraftThemeID, &out.PublishedTokens, &out.DraftTokens, &out.Sections, &out.DraftSections)
	if errors.Is(err, pgx.ErrNoRows) {
		return defaultRow(), nil
	}
	return out, err
}

func defaultRow() settingsRow {
	return settingsRow{
		ThemeID: "classic", PublishedTokens: []byte("{}"), DraftTokens: []byte("{}"),
		Sections: []byte("[]"), DraftSections: []byte("[]"),
	}
}

const saveDraftSQL = `INSERT INTO platform.tenant_settings
        (tenant_id, allow_negative_stock, draft_theme_id, draft_tokens, draft_sections)
    VALUES ($1, false, $2, $3::jsonb, $4::jsonb)
    ON CONFLICT (tenant_id) DO UPDATE SET
        draft_theme_id = EXCLUDED.draft_theme_id,
        draft_tokens = EXCLUDED.draft_tokens,
        draft_sections = EXCLUDED.draft_sections,
        updated_at = now()`

func (r *Repository) SaveDraft(ctx context.Context, tx db.DBTX, tenantID, themeID string, tokens, sections []byte) error {
	_, err := tx.Exec(ctx, saveDraftSQL, tenantID, themeID, string(tokens), string(sections))
	return err
}

const publishSQL = `UPDATE platform.tenant_settings
    SET storefront_theme_id = COALESCE(NULLIF(draft_theme_id, ''), storefront_theme_id),
        published_tokens = draft_tokens,
        sections = draft_sections,
        updated_at = now()`

func (r *Repository) Publish(ctx context.Context, tx db.DBTX, tenantID string) error {
	_, err := tx.Exec(ctx, publishSQL)
	return err
}

const primaryHostSQL = `SELECT host FROM platform.storefront_domains
    WHERE tenant_id = $1 AND status = 'active'
    ORDER BY (kind = 'custom') DESC, created_at
    LIMIT 1`

func (r *Repository) PrimaryHost(ctx context.Context, tenantID string) (string, error) {
	host := ""
	err := r.store.QueryRow(ctx, tenantID, primaryHostSQL, tenantID).Scan(&host)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return host, err
}

// TenantName reads the display name used to derive the shop slug.
// platform.tenants is pre-tenant infrastructure, outside RLS: raw pool access.
func (r *Repository) TenantName(ctx context.Context, tenantID string) (string, error) {
	name := ""
	err := r.store.Pool.QueryRow(ctx,
		`SELECT name FROM platform.tenants WHERE id=$1`, tenantID).Scan(&name)
	return name, err
}

// InsertDomain registers a shop host for the tenant. The host column is
// globally UNIQUE, so a taken subdomain surfaces as a 23505 for the caller to
// retry. platform.storefront_domains is outside RLS (it is the host->tenant
// map itself): raw pool access.
func (r *Repository) InsertDomain(ctx context.Context, id, tenantID, host string) error {
	_, err := r.store.Pool.Exec(ctx,
		`INSERT INTO platform.storefront_domains (id, tenant_id, host, kind, status)
		 VALUES ($1, $2, $3, 'subdomain', 'active')`, id, tenantID, host)
	return err
}
