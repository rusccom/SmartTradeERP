package ledger

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/shared/db"
)

func (s *Service) HasVariantMovements(ctx context.Context, tenantID, variantID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM ledger.inventory_movements
        WHERE variant_id=$1
	)`
	return scanExists(s.store.QueryRow(ctx, tenantID, query, variantID))
}

func (s *Service) HasWarehouseMovements(ctx context.Context, tenantID, warehouseID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM ledger.inventory_movements
        WHERE warehouse_id=$1
	)`
	return scanExists(s.store.QueryRow(ctx, tenantID, query, warehouseID))
}

func (s *Service) HasProductMovements(ctx context.Context, tenantID, productID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM ledger.inventory_movements m
        JOIN catalog.product_variants v ON v.id=m.variant_id
        WHERE v.product_id=$1
	)`
	return scanExists(s.store.QueryRow(ctx, tenantID, query, productID))
}

func scanExists(row pgx.Row) (bool, error) {
	exists := false
	err := row.Scan(&exists)
	return exists, err
}

const globalStockQuery = `SELECT COALESCE(running_qty,0), COALESCE(running_avg_cost,0)
        FROM ledger.cost_movement_results r
        JOIN ledger.inventory_movements m ON m.id=r.movement_id
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        WHERE r.variant_id=$1 AND b.status='active'
        ORDER BY r.sequence_num DESC LIMIT 1`

func (s *Service) GlobalStock(ctx context.Context, tenantID, variantID string) (decimal.Decimal, decimal.Decimal, error) {
	return scanStockPair(s.store.QueryRow(ctx, tenantID, globalStockQuery, variantID))
}

func (s *Service) GlobalStockTx(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	variantID string,
) (decimal.Decimal, decimal.Decimal, error) {
	return s.globalStock(ctx, tx, tenantID, variantID)
}

func (s *Service) globalStock(
	ctx context.Context,
	q db.DBTX,
	tenantID string,
	variantID string,
) (decimal.Decimal, decimal.Decimal, error) {
	return scanStockPair(q.QueryRow(ctx, globalStockQuery, variantID))
}

func scanStockPair(row pgx.Row) (decimal.Decimal, decimal.Decimal, error) {
	qty := decimal.Zero
	avg := decimal.Zero
	if err := row.Scan(&qty, &avg); err != nil {
		if err == pgx.ErrNoRows {
			return decimal.Zero, decimal.Zero, nil
		}
		return decimal.Zero, decimal.Zero, err
	}
	return qty, avg, nil
}

func (s *Service) ProfitByPeriod(ctx context.Context, filter ProfitFilter) (decimal.Decimal, error) {
	query, args := buildProfitPeriodQuery(filter)
	row := s.store.QueryRow(ctx, filter.TenantID, query, args...)
	profit := decimal.Zero
	err := row.Scan(&profit)
	return profit, err
}

func buildProfitPeriodQuery(filter ProfitFilter) (string, []any) {
	parts := []string{"b.status='active'", "m.movement_date BETWEEN $1 AND $2"}
	args := []any{filter.FromDate, filter.ToDate}
	parts, args = appendProfitFilter(parts, args, "m.warehouse_id", filter.WarehouseID)
	parts, args = appendProfitFilter(parts, args, "m.variant_id", filter.VariantID)
	query := `SELECT COALESCE(SUM(COALESCE(r.gross_profit,0)), 0)
        FROM ledger.cost_movement_results r
        JOIN ledger.inventory_movements m ON m.id=r.movement_id
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        WHERE ` + strings.Join(parts, " AND ")
	return query, args
}

func appendProfitFilter(parts []string, args []any, field string, value string) ([]string, []any) {
	if value == "" {
		return parts, args
	}
	args = append(args, value)
	parts = append(parts, field+"=$"+strconv.Itoa(len(args)))
	return parts, args
}

func (s *Service) ProfitByDocumentItem(ctx context.Context, tenantID, documentItemID string) (decimal.Decimal, error) {
	query := `SELECT COALESCE(gross_profit,0)
        FROM ledger.document_item_financials
        WHERE document_item_id=$1`
	row := s.store.QueryRow(ctx, tenantID, query, documentItemID)
	return scanDecimal(row)
}

func (s *Service) ProfitByDocument(ctx context.Context, tenantID, documentID string) (decimal.Decimal, error) {
	query := `SELECT COALESCE(SUM(gross_profit),0)
        FROM ledger.document_item_financials
        WHERE document_id=$1`
	row := s.store.QueryRow(ctx, tenantID, query, documentID)
	return scanDecimal(row)
}

func scanDecimal(row pgx.Row) (decimal.Decimal, error) {
	value := decimal.Zero
	if err := row.Scan(&value); err != nil {
		if err == pgx.ErrNoRows {
			return decimal.Zero, nil
		}
		return decimal.Zero, err
	}
	return value, nil
}

func (s *Service) Movements(ctx context.Context, tenantID, variantID string) ([]Movement, error) {
	query := `SELECT r.sequence_num, m.movement_date::text, m.direction, m.reason,
        r.qty_delta, r.running_qty, r.running_avg_cost, COALESCE(r.cogs_amount,0),
        COALESCE(r.gross_profit,0)
        FROM ledger.inventory_movements m
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        JOIN ledger.cost_movement_results r ON r.movement_id=m.id
        WHERE m.variant_id=$1 AND b.status='active'
        ORDER BY r.sequence_num`
	rows, err := s.store.Query(ctx, tenantID, query, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMovements(rows)
}

func scanMovements(rows pgx.Rows) ([]Movement, error) {
	items := make([]Movement, 0)
	for rows.Next() {
		item := Movement{}
		err := rows.Scan(&item.SequenceNum, &item.Date, &item.Type, &item.Reason, &item.Qty,
			&item.RunningQty, &item.RunningAvg, &item.COGS, &item.Profit)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
