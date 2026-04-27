package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type movementRow struct {
	id        string
	itemType  string
	qty       decimal.Decimal
	unitPrice decimal.Decimal
	revenue   *decimal.Decimal
	warehouse string
}

type movementCalc struct {
	row      movementRow
	seq      int64
	previous calcState
	result   calcResult
}

type rebuildRun struct {
	ctx           context.Context
	tx            pgx.Tx
	tenantID      string
	variantID     string
	version       string
	allowNegative bool
}

func (s *Service) RebuildAffected(ctx context.Context, tx pgx.Tx, tenantID string, affected []VariantSequence) error {
	if err := s.lockTenant(ctx, tx, tenantID); err != nil {
		return err
	}
	allowNegative, err := s.AllowNegativeStock(ctx, tx, tenantID)
	if err != nil {
		return err
	}
	for _, item := range affected {
		if err := s.rebuildVariant(ctx, tx, tenantID, item.VariantID, allowNegative); err != nil {
			return err
		}
	}
	return s.RefreshTenantProjections(ctx, tx, tenantID)
}

func (s *Service) RebuildVariant(ctx context.Context, tx pgx.Tx, tenantID, variantID string) error {
	allowNegative, err := s.AllowNegativeStock(ctx, tx, tenantID)
	if err != nil {
		return err
	}
	return s.rebuildVariant(ctx, tx, tenantID, variantID, allowNegative)
}

func (s *Service) rebuildVariant(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	variantID string,
	allowNegative bool,
) error {
	run := newRebuildRun(ctx, tx, tenantID, variantID, allowNegative)
	if err := s.lockVariant(run); err != nil {
		return err
	}
	if err := s.clearVariant(run); err != nil {
		return err
	}
	return s.rebuildVariantRows(run)
}

func (s *Service) rebuildVariantRows(run rebuildRun) error {
	rows, err := s.activeRows(run)
	if err != nil {
		return err
	}
	defer rows.Close()
	return s.writeResults(run, rows)
}

func newRebuildRun(ctx context.Context, tx pgx.Tx, tenantID, variantID string, allowNegative bool) rebuildRun {
	return rebuildRun{
		ctx: ctx, tx: tx, tenantID: tenantID, variantID: variantID,
		version: uuid.NewString(), allowNegative: allowNegative,
	}
}

func (s *Service) lockTenant(ctx context.Context, tx pgx.Tx, tenantID string) error {
	query := `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`
	_, err := tx.Exec(ctx, query, "tenant:"+tenantID)
	return err
}

func (s *Service) lockVariant(run rebuildRun) error {
	query := `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`
	_, err := run.tx.Exec(run.ctx, query, run.tenantID+":"+run.variantID)
	return err
}

func (s *Service) clearVariant(run rebuildRun) error {
	if err := s.deleteResults(run); err != nil {
		return err
	}
	query := `DELETE FROM ledger.stock_balances
        WHERE tenant_id=$1 AND variant_id=$2`
	_, err := run.tx.Exec(run.ctx, query, run.tenantID, run.variantID)
	return err
}

func (s *Service) deleteResults(run rebuildRun) error {
	query := `DELETE FROM ledger.cost_movement_results
        WHERE tenant_id=$1 AND variant_id=$2`
	_, err := run.tx.Exec(run.ctx, query, run.tenantID, run.variantID)
	return err
}

func (s *Service) activeRows(run rebuildRun) (pgx.Rows, error) {
	query := `SELECT m.id::text, m.direction, m.qty, m.unit_price,
        m.revenue_amount, m.warehouse_id::text
        FROM ledger.inventory_movements m
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        WHERE m.tenant_id=$1 AND m.variant_id=$2 AND b.status='active'
        ORDER BY m.movement_date, b.posted_at, m.posting_order, m.created_at, m.id
        FOR UPDATE OF m`
	return run.tx.Query(run.ctx, query, run.tenantID, run.variantID)
}

func (s *Service) writeResults(run rebuildRun, rows pgx.Rows) error {
	state := zeroState()
	warehouseQty := map[string]decimal.Decimal{}
	seq := int64(1)
	for rows.Next() {
		next, err := s.writeMovementResult(run, rows, warehouseQty, state, seq)
		if err != nil {
			return err
		}
		state = next
		seq++
	}
	return rows.Err()
}

func (s *Service) writeMovementResult(
	run rebuildRun,
	rows pgx.Rows,
	warehouseQty map[string]decimal.Decimal,
	state calcState,
	seq int64,
) (calcState, error) {
	calc, err := nextMovementCalc(rows, state, seq)
	if err != nil {
		return state, err
	}
	if err := applyStockPolicy(run, warehouseQty, calc); err != nil {
		return state, err
	}
	if err := s.insertResult(run, calc); err != nil {
		return state, err
	}
	return calc.result.state, s.upsertBalance(run, calc)
}

func applyStockPolicy(
	run rebuildRun,
	warehouseQty map[string]decimal.Decimal,
	calc movementCalc,
) error {
	nextQty := warehouseQty[calc.row.warehouse].Add(qtyDelta(calc))
	if !run.allowNegative && nextQty.LessThan(decimal.Zero) {
		return ErrNegativeStock
	}
	if !run.allowNegative && calc.result.state.qty.LessThan(decimal.Zero) {
		return ErrNegativeStock
	}
	warehouseQty[calc.row.warehouse] = nextQty
	return nil
}

func nextMovementCalc(rows pgx.Rows, state calcState, seq int64) (movementCalc, error) {
	row, err := scanMovementRow(rows)
	if err != nil {
		return movementCalc{}, err
	}
	return movementCalc{row: row, seq: seq, previous: state, result: calculateRow(state, row)}, nil
}

func scanMovementRow(rows pgx.Rows) (movementRow, error) {
	row := movementRow{}
	err := rows.Scan(&row.id, &row.itemType, &row.qty, &row.unitPrice, &row.revenue, &row.warehouse)
	return row, err
}

func calculateRow(state calcState, row movementRow) calcResult {
	if row.itemType == "IN" {
		return applyIn(state, row.qty, row.unitPrice, row.revenue)
	}
	return applyOut(state, row.qty, row.revenue)
}

func (s *Service) insertResult(run rebuildRun, calc movementCalc) error {
	query := `INSERT INTO ledger.cost_movement_results
        (movement_id, tenant_id, variant_id, sequence_num, qty_delta,
         unit_cost, movement_cost, revenue_amount, cogs_amount, gross_profit,
         running_qty, running_avg_cost, running_inventory_value, calculation_version)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := run.tx.Exec(run.ctx, query, calc.row.id, run.tenantID, run.variantID,
		calc.seq, qtyDelta(calc), unitCost(calc), movementCost(calc), calc.row.revenue,
		calc.result.cogs, calc.result.profit, calc.result.state.qty,
		calc.result.state.avg, inventoryValue(calc), run.version)
	return err
}

func (s *Service) upsertBalance(run rebuildRun, calc movementCalc) error {
	query := `INSERT INTO ledger.stock_balances
        (tenant_id, variant_id, warehouse_id, qty)
        VALUES ($1,$2,$3,$4)
        ON CONFLICT (tenant_id, variant_id, warehouse_id)
        DO UPDATE SET qty=ledger.stock_balances.qty + EXCLUDED.qty, updated_at=now()`
	_, err := run.tx.Exec(run.ctx, query, run.tenantID, run.variantID,
		calc.row.warehouse, qtyDelta(calc))
	return err
}

func qtyDelta(calc movementCalc) decimal.Decimal {
	if calc.row.itemType == "IN" {
		return calc.row.qty
	}
	return calc.row.qty.Neg()
}

func unitCost(calc movementCalc) decimal.Decimal {
	if calc.row.itemType == "OUT" {
		return calc.previous.avg.Round(4)
	}
	return calc.row.unitPrice.Round(4)
}

func movementCost(calc movementCalc) decimal.Decimal {
	if calc.result.cogs != nil {
		return *calc.result.cogs
	}
	return calc.row.qty.Mul(calc.row.unitPrice).Round(4)
}

func inventoryValue(calc movementCalc) decimal.Decimal {
	return calc.result.state.qty.Mul(calc.result.state.avg).Round(4)
}
