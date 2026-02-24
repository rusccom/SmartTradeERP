package ledger

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/shared/db"
)

type Service struct {
    store *db.Store
}

func NewService(store *db.Store) *Service {
    return &Service{store: store}
}

func (s *Service) Append(ctx context.Context, tx pgx.Tx, input EntryInput) (int64, error) {
    seq, err := s.reserveSequence(ctx, tx, input)
    if err != nil {
        return 0, err
    }
    err = s.insertEntry(ctx, tx, input, seq)
    if err != nil {
        return 0, err
    }
    return seq, s.Recalculate(ctx, tx, input.TenantID, input.VariantID, seq)
}

func (s *Service) Recalculate(ctx context.Context, tx pgx.Tx, tenantID, variantID string, fromSeq int64) error {
    state, err := s.loadStateBefore(ctx, tx, tenantID, variantID, fromSeq)
    if err != nil {
        return err
    }
    rows, err := s.loadRowsForRecalculate(ctx, tx, tenantID, variantID, fromSeq)
    if err != nil {
        return err
    }
    defer rows.Close()
    return s.recalculateRows(ctx, tx, rows, state)
}

func (s *Service) reserveSequence(ctx context.Context, tx pgx.Tx, input EntryInput) (int64, error) {
    insertAt, err := s.findInsertSequence(ctx, tx, input.TenantID, input.VariantID, input.Date)
    if err != nil {
        return 0, err
    }
    if insertAt == 0 {
        return s.findNextSequence(ctx, tx, input.TenantID, input.VariantID)
    }
    if err := s.shiftSequence(ctx, tx, input.TenantID, input.VariantID, insertAt); err != nil {
        return 0, err
    }
    return insertAt, nil
}

func (s *Service) findInsertSequence(ctx context.Context, tx pgx.Tx, tenantID, variantID string, date time.Time) (int64, error) {
    query := `SELECT COALESCE(MIN(sequence_num), 0)
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND variant_id=$2 AND date>$3`
    row := tx.QueryRow(ctx, query, tenantID, variantID, date)
    var sequence int64
    return sequence, row.Scan(&sequence)
}

func (s *Service) findNextSequence(ctx context.Context, tx pgx.Tx, tenantID, variantID string) (int64, error) {
    query := `SELECT COALESCE(MAX(sequence_num), 0) + 1
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND variant_id=$2`
    row := tx.QueryRow(ctx, query, tenantID, variantID)
    var sequence int64
    return sequence, row.Scan(&sequence)
}

func (s *Service) shiftSequence(ctx context.Context, tx pgx.Tx, tenantID, variantID string, fromSeq int64) error {
    query := `UPDATE ledger.cost_ledger
        SET sequence_num = sequence_num + 1
        WHERE tenant_id=$1 AND variant_id=$2 AND sequence_num >= $3`
    _, err := tx.Exec(ctx, query, tenantID, variantID, fromSeq)
    return err
}

func (s *Service) insertEntry(ctx context.Context, tx pgx.Tx, input EntryInput, seq int64) error {
    query := `INSERT INTO ledger.cost_ledger
        (tenant_id, variant_id, document_id, document_item_id, warehouse_id, date,
         sequence_num, type, reason, qty, unit_price, total_amount, revenue)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
    _, err := tx.Exec(ctx, query, input.TenantID, input.VariantID, input.DocumentID,
        input.DocumentItemID, input.WarehouseID, input.Date, seq, input.Type, input.Reason,
        input.Qty, input.UnitPrice, input.TotalAmount, input.Revenue)
    return err
}

func (s *Service) loadStateBefore(ctx context.Context, tx pgx.Tx, tenantID, variantID string, seq int64) (calcState, error) {
    query := `SELECT COALESCE(running_qty,0), COALESCE(running_avg,0)
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND variant_id=$2 AND sequence_num < $3
        ORDER BY sequence_num DESC LIMIT 1`
    row := tx.QueryRow(ctx, query, tenantID, variantID, seq)
    state := calcState{}
    if err := row.Scan(&state.qty, &state.avg); err != nil {
        if err == pgx.ErrNoRows {
            return calcState{}, nil
        }
        return calcState{}, err
    }
    return state, nil
}

func (s *Service) loadRowsForRecalculate(ctx context.Context, tx pgx.Tx, tenantID, variantID string, seq int64) (pgx.Rows, error) {
    query := `SELECT id, type, qty, unit_price, revenue
        FROM ledger.cost_ledger
        WHERE tenant_id=$1 AND variant_id=$2 AND sequence_num >= $3
        ORDER BY sequence_num FOR UPDATE`
    return tx.Query(ctx, query, tenantID, variantID, seq)
}

func (s *Service) recalculateRows(ctx context.Context, tx pgx.Tx, rows pgx.Rows, state calcState) error {
    for rows.Next() {
        item, err := scanRow(rows)
        if err != nil {
            return err
        }
        result := calculateRow(state, item)
        if err := s.updateComputed(ctx, tx, item.id, result); err != nil {
            return err
        }
        state = result.state
    }
    return rows.Err()
}

func calculateRow(state calcState, row ledgerRow) calcResult {
    if row.itemType == "IN" {
        return applyIn(state, row.qty, row.unitPrice)
    }
    return applyOut(state, row.qty, row.revenue)
}

func (s *Service) updateComputed(ctx context.Context, tx pgx.Tx, id int64, result calcResult) error {
    query := `UPDATE ledger.cost_ledger
        SET running_qty=$2, running_avg=$3, cogs=$4, profit=$5, updated_at=now()
        WHERE id=$1`
    _, err := tx.Exec(ctx, query, id, result.state.qty, result.state.avg, result.cogs, result.profit)
    return err
}

type ledgerRow struct {
    id       int64
    itemType string
    qty      float64
    unitPrice float64
    revenue  *float64
}

func scanRow(rows pgx.Rows) (ledgerRow, error) {
    row := ledgerRow{}
    err := rows.Scan(&row.id, &row.itemType, &row.qty, &row.unitPrice, &row.revenue)
    return row, err
}
