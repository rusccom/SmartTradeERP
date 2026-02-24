package ledger

import (
    "context"
    "strconv"
    "strings"
    "time"

    "github.com/jackc/pgx/v5"
)

type VariantSequence struct {
    VariantID string
    Earliest  int64
}

func (s *Service) DeleteForDocument(ctx context.Context, tx pgx.Tx, tenantID, documentID string) ([]VariantSequence, error) {
    variants, err := s.findAffectedVariants(ctx, tx, tenantID, documentID)
    if err != nil {
        return nil, err
    }
    if err := s.deleteDocumentRows(ctx, tx, tenantID, documentID); err != nil {
        return nil, err
    }
    return variants, nil
}

func (s *Service) findAffectedVariants(ctx context.Context, tx pgx.Tx, tenantID, documentID string) ([]VariantSequence, error) {
    query := `SELECT variant_id::text, MIN(sequence_num)
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND document_id=$2
        GROUP BY variant_id`
    rows, err := tx.Query(ctx, query, tenantID, documentID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanVariantSequences(rows)
}

func scanVariantSequences(rows pgx.Rows) ([]VariantSequence, error) {
    items := make([]VariantSequence, 0)
    for rows.Next() {
        item := VariantSequence{}
        if err := rows.Scan(&item.VariantID, &item.Earliest); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func (s *Service) deleteDocumentRows(ctx context.Context, tx pgx.Tx, tenantID, documentID string) error {
    query := `DELETE FROM ledger.cost_ledger WHERE tenant_id=$1 AND document_id=$2`
    _, err := tx.Exec(ctx, query, tenantID, documentID)
    return err
}

func (s *Service) HasVariantMovements(ctx context.Context, tenantID, variantID string) (bool, error) {
    query := `SELECT EXISTS(
        SELECT 1 FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND variant_id=$2
    )`
    row := s.store.Pool.QueryRow(ctx, query, tenantID, variantID)
    var exists bool
    return exists, row.Scan(&exists)
}

func (s *Service) HasWarehouseMovements(ctx context.Context, tenantID, warehouseID string) (bool, error) {
    query := `SELECT EXISTS(
        SELECT 1 FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND warehouse_id=$2
    )`
    row := s.store.Pool.QueryRow(ctx, query, tenantID, warehouseID)
    var exists bool
    return exists, row.Scan(&exists)
}

func (s *Service) HasProductMovements(ctx context.Context, tenantID, productID string) (bool, error) {
    query := `SELECT EXISTS(
        SELECT 1 FROM ledger.cost_ledger l
        JOIN catalog.product_variants v ON v.id = l.variant_id
        JOIN catalog.products p ON p.id = v.product_id
        WHERE l.tenant_id=$1 AND p.id=$2
    )`
    row := s.store.Pool.QueryRow(ctx, query, tenantID, productID)
    var exists bool
    return exists, row.Scan(&exists)
}

func (s *Service) GlobalStock(ctx context.Context, tenantID, variantID string) (float64, float64, error) {
    query := `SELECT COALESCE(running_qty,0), COALESCE(running_avg,0)
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND variant_id=$2
        ORDER BY sequence_num DESC LIMIT 1`
    row := s.store.Pool.QueryRow(ctx, query, tenantID, variantID)
    qty := 0.0
    avg := 0.0
    if err := row.Scan(&qty, &avg); err != nil {
        if err == pgx.ErrNoRows {
            return 0, 0, nil
        }
        return 0, 0, err
    }
    return qty, avg, nil
}

func (s *Service) WarehouseStock(ctx context.Context, tenantID, variantID, warehouseID string) (float64, error) {
    query := `SELECT COALESCE(SUM(CASE WHEN type='IN' THEN qty ELSE -qty END), 0)
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND variant_id=$2 AND warehouse_id=$3`
    row := s.store.Pool.QueryRow(ctx, query, tenantID, variantID, warehouseID)
    qty := 0.0
    return qty, row.Scan(&qty)
}

func (s *Service) ProfitByPeriod(
    ctx context.Context,
    tenantID string,
    fromDate time.Time,
    toDate time.Time,
    warehouseID string,
    variantID string,
) (float64, error) {
    query, args := buildProfitPeriodQuery(tenantID, fromDate, toDate, warehouseID, variantID)
    row := s.store.Pool.QueryRow(ctx, query, args...)
    profit := 0.0
    return profit, row.Scan(&profit)
}

func buildProfitPeriodQuery(
    tenantID string,
    fromDate time.Time,
    toDate time.Time,
    warehouseID string,
    variantID string,
) (string, []any) {
    parts := []string{"tenant_id=$1", "date BETWEEN $2 AND $3", "type='OUT'"}
    args := []any{tenantID, fromDate, toDate}
    parts, args = appendWarehouseFilter(parts, args, warehouseID)
    parts, args = appendVariantFilter(parts, args, variantID)
    query := `SELECT COALESCE(SUM(profit), 0)
        FROM ledger.cost_ledger
        WHERE ` + strings.Join(parts, " AND ")
    return query, args
}

func appendWarehouseFilter(parts []string, args []any, warehouseID string) ([]string, []any) {
    if warehouseID == "" {
        return parts, args
    }
    parts = append(parts, "warehouse_id=$4")
    args = append(args, warehouseID)
    return parts, args
}

func appendVariantFilter(parts []string, args []any, variantID string) ([]string, []any) {
    if variantID == "" {
        return parts, args
    }
    index := len(args) + 1
    parts = append(parts, "variant_id=$"+intToString(index))
    args = append(args, variantID)
    return parts, args
}

func intToString(value int) string {
    return strconv.Itoa(value)
}

func (s *Service) ProfitByDocumentItem(ctx context.Context, tenantID, documentItemID string) (float64, error) {
    query := `SELECT COALESCE(SUM(profit),0)
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND document_item_id=$2`
    row := s.store.Pool.QueryRow(ctx, query, tenantID, documentItemID)
    profit := 0.0
    return profit, row.Scan(&profit)
}

func (s *Service) ProfitByDocument(ctx context.Context, tenantID, documentID string) (float64, error) {
    query := `SELECT COALESCE(SUM(profit),0)
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND document_id=$2`
    row := s.store.Pool.QueryRow(ctx, query, tenantID, documentID)
    profit := 0.0
    return profit, row.Scan(&profit)
}

func (s *Service) Movements(ctx context.Context, tenantID, variantID string) ([]Movement, error) {
    query := `SELECT sequence_num, date::text, type, reason, qty,
        running_qty, running_avg, COALESCE(cogs,0), COALESCE(profit,0)
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND variant_id=$2
        ORDER BY sequence_num`
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
