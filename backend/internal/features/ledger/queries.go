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
        WHERE tenant_id=$1 AND variant_id=$2
	)`
	return scanExists(s.store.Pool.QueryRow(ctx, query, tenantID, variantID))
}

func (s *Service) HasWarehouseMovements(ctx context.Context, tenantID, warehouseID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM ledger.inventory_movements
        WHERE tenant_id=$1 AND warehouse_id=$2
	)`
	return scanExists(s.store.Pool.QueryRow(ctx, query, tenantID, warehouseID))
}

func (s *Service) HasProductMovements(ctx context.Context, tenantID, productID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM ledger.inventory_movements m
        JOIN catalog.product_variants v ON v.id=m.variant_id
        WHERE m.tenant_id=$1 AND v.tenant_id=$1 AND v.product_id=$2
	)`
	return scanExists(s.store.Pool.QueryRow(ctx, query, tenantID, productID))
}

func scanExists(row pgx.Row) (bool, error) {
	exists := false
	err := row.Scan(&exists)
	return exists, err
}

func (s *Service) GlobalStock(ctx context.Context, tenantID, variantID string) (decimal.Decimal, decimal.Decimal, error) {
	return s.globalStock(ctx, s.store.Pool, tenantID, variantID)
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
	query := `SELECT COALESCE(running_qty,0), COALESCE(running_avg_cost,0)
        FROM ledger.cost_movement_results r
        JOIN ledger.inventory_movements m ON m.id=r.movement_id
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        WHERE r.tenant_id=$1 AND r.variant_id=$2 AND b.status='active'
        ORDER BY r.sequence_num DESC LIMIT 1`
	return scanStockPair(q.QueryRow(ctx, query, tenantID, variantID))
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

func (s *Service) WarehouseStock(ctx context.Context, tenantID, variantID, warehouseID string) (decimal.Decimal, error) {
	return s.warehouseStock(ctx, s.store.Pool, tenantID, variantID, warehouseID)
}

func (s *Service) WarehouseStockTx(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	variantID string,
	warehouseID string,
) (decimal.Decimal, error) {
	return s.warehouseStock(ctx, tx, tenantID, variantID, warehouseID)
}

func (s *Service) warehouseStock(
	ctx context.Context,
	q db.DBTX,
	tenantID string,
	variantID string,
	warehouseID string,
) (decimal.Decimal, error) {
	query := `SELECT COALESCE(SUM(CASE WHEN m.direction='IN' THEN m.qty ELSE -m.qty END), 0)
        FROM ledger.inventory_movements m
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        WHERE m.tenant_id=$1 AND m.variant_id=$2 AND m.warehouse_id=$3
            AND b.status='active'`
	return scanWarehouseQty(q.QueryRow(ctx, query, tenantID, variantID, warehouseID))
}

func scanWarehouseQty(row pgx.Row) (decimal.Decimal, error) {
	qty := decimal.Zero
	if err := row.Scan(&qty); err != nil {
		if err == pgx.ErrNoRows {
			return decimal.Zero, nil
		}
		return decimal.Zero, err
	}
	return qty, nil
}

func (s *Service) ProfitByPeriod(ctx context.Context, filter ProfitFilter) (decimal.Decimal, error) {
	query, args := buildProfitPeriodQuery(filter)
	row := s.store.Pool.QueryRow(ctx, query, args...)
	profit := decimal.Zero
	err := row.Scan(&profit)
	return profit, err
}

func buildProfitPeriodQuery(filter ProfitFilter) (string, []any) {
	parts := []string{"m.tenant_id=$1", "b.status='active'", "m.movement_date BETWEEN $2 AND $3"}
	args := []any{filter.TenantID, filter.FromDate, filter.ToDate}
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
        WHERE tenant_id=$1 AND document_item_id=$2`
	row := s.store.Pool.QueryRow(ctx, query, tenantID, documentItemID)
	return scanDecimal(row)
}

func (s *Service) ProfitByDocument(ctx context.Context, tenantID, documentID string) (decimal.Decimal, error) {
	query := `SELECT COALESCE(SUM(gross_profit),0)
        FROM ledger.document_item_financials
        WHERE tenant_id=$1 AND document_id=$2`
	row := s.store.Pool.QueryRow(ctx, query, tenantID, documentID)
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
        m.qty, r.running_qty, r.running_avg_cost, COALESCE(r.cogs_amount,0),
        COALESCE(r.gross_profit,0)
        FROM ledger.inventory_movements m
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        JOIN ledger.cost_movement_results r ON r.movement_id=m.id
        WHERE m.tenant_id=$1 AND m.variant_id=$2 AND b.status='active'
        ORDER BY r.sequence_num`
	rows, err := s.store.Pool.Query(ctx, query, tenantID, variantID)
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
