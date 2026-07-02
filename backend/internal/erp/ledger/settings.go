package ledger

import (
	"context"

	"smarterp/backend/internal/shared/db"
)

func (s *Service) AllowNegativeStock(ctx context.Context, q db.DBTX, tenantID string) (bool, error) {
	query := `SELECT COALESCE((
        SELECT allow_negative_stock
        FROM platform.tenant_settings
    ), false)`
	allowed := false
	err := q.QueryRow(ctx, query).Scan(&allowed)
	return allowed, err
}
