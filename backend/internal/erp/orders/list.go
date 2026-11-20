package orders

import (
	"context"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
)

func list(ctx context.Context, tx db.DBTX, tenantID string, q httpx.ListQuery) ([]Summary, int, error) {
	const from = ` FROM storefront.orders o WHERE o.tenant_id=$1 AND ($2='' OR strpos(lower(o.number||' '||o.customer_name),lower($2))>0)`
	var total int
	if err := tx.QueryRow(ctx, `SELECT count(*)`+from, tenantID, q.Search).Scan(&total); err != nil {
		return nil, 0, err
	}
	query := `SELECT ` + summaryColumns + from + ` ORDER BY ` + orderBy(q) + `,o.id LIMIT $3 OFFSET $4`
	rows, err := tx.Query(ctx, query, tenantID, q.Search, q.PerPage, httpx.Offset(q.Page, q.PerPage))
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Summary, error) {
		var s Summary
		err := row.Scan(s.destinations()...)
		return s, err
	})
	return items, total, err
}

func orderBy(q httpx.ListQuery) string {
	columns := map[string]string{"number": "o.number", "customer_name": "o.customer_name", "total_with_delivery": "(o.total_amount+o.delivery_cost)", "status": "o.status"}
	col := columns[q.SortBy]
	if col == "" {
		col = "o.created_at"
	}
	if q.SortDir == "asc" {
		return col + " ASC NULLS LAST"
	}
	return col + " DESC NULLS LAST"
}
