package currencies

import (
	"context"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
)

type Repository struct {
	store *db.Store
}

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) List(ctx context.Context, tenantID string, page, perPage int) ([]Currency, int, error) {
	total, err := r.countTenant(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	items, err := r.loadTenant(ctx, tenantID, page, perPage)
	return items, total, err
}

func (r *Repository) Options(ctx context.Context, page, perPage int) ([]CurrencyOption, int, error) {
	total, err := r.countOptions(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := r.loadOptions(ctx, page, perPage)
	return items, total, err
}

func (r *Repository) countTenant(ctx context.Context, tenantID string) (int, error) {
	row := r.store.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM platform.tenant_currencies WHERE tenant_id=$1`, tenantID)
	count := 0
	err := row.Scan(&count)
	return count, err
}

func (r *Repository) countOptions(ctx context.Context) (int, error) {
	row := r.store.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM platform.currencies WHERE is_active=true`)
	count := 0
	err := row.Scan(&count)
	return count, err
}

func (r *Repository) loadTenant(ctx context.Context, tenantID string, page, perPage int) ([]Currency, error) {
	query := `SELECT tc.id::text, c.id::text, c.code, c.name, c.symbol,
		COALESCE(tc.display_symbol,''), c.decimal_places, tc.is_base,
		tc.is_enabled, tc.created_at::text, tc.updated_at::text
		FROM platform.tenant_currencies tc
		JOIN platform.currencies c ON c.id=tc.currency_id
		WHERE tc.tenant_id=$1
		ORDER BY tc.is_base DESC, c.code ASC
		LIMIT $2 OFFSET $3`
	rows, err := r.store.Pool.Query(ctx, query, tenantID, perPage, httpx.Offset(page, perPage))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCurrencies(rows)
}

func scanCurrencies(rows pgx.Rows) ([]Currency, error) {
	items := make([]Currency, 0)
	for rows.Next() {
		item := Currency{}
		err := rows.Scan(&item.ID, &item.CurrencyID, &item.Code, &item.Name,
			&item.Symbol, &item.DisplaySymbol, &item.DecimalPlaces, &item.IsBase,
			&item.IsEnabled, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) loadOptions(ctx context.Context, page, perPage int) ([]CurrencyOption, error) {
	query := `SELECT id::text, code, name, symbol, decimal_places
		FROM platform.currencies
		WHERE is_active=true
		ORDER BY code ASC
		LIMIT $1 OFFSET $2`
	rows, err := r.store.Pool.Query(ctx, query, perPage, httpx.Offset(page, perPage))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func scanOptions(rows pgx.Rows) ([]CurrencyOption, error) {
	items := make([]CurrencyOption, 0)
	for rows.Next() {
		item := CurrencyOption{}
		err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Symbol, &item.DecimalPlaces)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) Count(ctx context.Context, tx db.DBTX, tenantID string) (int, error) {
	row := tx.QueryRow(ctx, `SELECT COUNT(*) FROM platform.tenant_currencies WHERE tenant_id=$1`, tenantID)
	count := 0
	err := row.Scan(&count)
	return count, err
}

func (r *Repository) ClearBase(ctx context.Context, tx db.DBTX, tenantID string) error {
	_, err := tx.Exec(ctx, `UPDATE platform.tenant_currencies SET is_base=false WHERE tenant_id=$1`, tenantID)
	return err
}

func (r *Repository) Create(ctx context.Context, tx db.DBTX, tenantID, id string, req CreateRequest) error {
	query := `INSERT INTO platform.tenant_currencies
		(id, tenant_id, currency_id, is_base, is_enabled, display_symbol)
		VALUES ($1,$2,$3,$4,true,$5)`
	_, err := tx.Exec(ctx, query, id, tenantID, req.CurrencyID, req.IsBase, req.DisplaySymbol)
	return err
}

func (r *Repository) SaveBaseSetting(ctx context.Context, tx db.DBTX, tenantID, currencyID string) error {
	query := `INSERT INTO platform.tenant_settings (tenant_id, allow_negative_stock, base_currency_id)
		VALUES ($1,false,$2)
		ON CONFLICT (tenant_id) DO UPDATE SET base_currency_id=EXCLUDED.base_currency_id, updated_at=now()`
	_, err := tx.Exec(ctx, query, tenantID, currencyID)
	return err
}
