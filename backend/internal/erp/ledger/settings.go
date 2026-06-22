package ledger

import (
	"context"

	"smarterp/backend/internal/shared/db"
)

func (s *Service) AllowNegativeStock(ctx context.Context, q db.DBTX, tenantID string) (bool, error) {
	query := `SELECT COALESCE((
        SELECT allow_negative_stock
        FROM platform.tenant_settings
        WHERE tenant_id=$1
    ), false)`
	allowed := false
	err := q.QueryRow(ctx, query, tenantID).Scan(&allowed)
	return allowed, err
}
