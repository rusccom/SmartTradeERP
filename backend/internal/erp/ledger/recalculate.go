package ledger

import (
	"context"
	"sort"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type rebuildRun struct {
	ctx           context.Context
	tx            pgx.Tx
	tenantID      string
	variantID     string
	version       string
	allowNegative bool
}

func (s *Service) RebuildAffected(ctx context.Context, tx pgx.Tx, tenantID string, affected []VariantSequence) error {
	variantIDs := sortedVariantIDs(affected)
	if len(variantIDs) == 0 {
		return nil
	}
	allowNegative, err := s.AllowNegativeStock(ctx, tx, tenantID)
	if err != nil {
		return err
	}
	for _, variantID := range variantIDs {
		if err := s.rebuildVariant(ctx, tx, tenantID, variantID, allowNegative); err != nil {
			return err
		}
	}
	return s.refreshVariantProjections(ctx, tx, tenantID, variantIDs)
}

func sortedVariantIDs(affected []VariantSequence) []string {
	ids := make([]string, 0, len(affected))
	for _, item := range affected {
		ids = append(ids, item.VariantID)
	}
	sort.Strings(ids)
	return ids
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
	events, err := s.loadEvents(run)
	if err != nil {
		return err
	}
	results, err := foldMovements(events, run.allowNegative)
	if err != nil {
		return err
	}
	return s.persistResults(run, results)
}

func newRebuildRun(ctx context.Context, tx pgx.Tx, tenantID, variantID string, allowNegative bool) rebuildRun {
	return rebuildRun{
		ctx: ctx, tx: tx, tenantID: tenantID, variantID: variantID,
		version: uuid.NewString(), allowNegative: allowNegative,
	}
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

func (s *Service) loadEvents(run rebuildRun) ([]replayEvent, error) {
	query := `SELECT m.id::text, m.direction, m.reason, m.qty, m.unit_price,
        m.revenue_amount, m.warehouse_id::text
        FROM ledger.inventory_movements m
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        WHERE m.tenant_id=$1 AND m.variant_id=$2 AND b.status='active'
        ORDER BY m.movement_date, b.posted_at, m.posting_order, m.created_at, m.id
        FOR UPDATE OF m`
	rows, err := run.tx.Query(run.ctx, query, run.tenantID, run.variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvents(rows)
}

func scanEvents(rows pgx.Rows) ([]replayEvent, error) {
	events := make([]replayEvent, 0)
	for rows.Next() {
		event := replayEvent{}
		err := rows.Scan(&event.id, &event.direction, &event.reason, &event.qty,
			&event.unitPrice, &event.revenue, &event.warehouse)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (s *Service) persistResults(run rebuildRun, results []replayResult) error {
	for _, result := range results {
		if err := s.insertResult(run, result); err != nil {
			return err
		}
		if err := s.upsertBalance(run, result); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) insertResult(run rebuildRun, result replayResult) error {
	query := `INSERT INTO ledger.cost_movement_results
        (movement_id, tenant_id, variant_id, sequence_num, qty_delta,
         unit_cost, movement_cost, revenue_amount, cogs_amount, gross_profit,
         running_qty, running_avg_cost, running_inventory_value, calculation_version)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := run.tx.Exec(run.ctx, query, result.event.id, run.tenantID, run.variantID,
		result.seq, result.qtyDelta, result.unitCost, movementCost(result), result.event.revenue,
		result.cogs, result.profit, result.runningQty,
		result.runningAvg, inventoryValue(result), run.version)
	return err
}

func (s *Service) upsertBalance(run rebuildRun, result replayResult) error {
	query := `INSERT INTO ledger.stock_balances
        (tenant_id, variant_id, warehouse_id, qty)
        VALUES ($1,$2,$3,$4)
        ON CONFLICT (tenant_id, variant_id, warehouse_id)
        DO UPDATE SET qty=ledger.stock_balances.qty + EXCLUDED.qty, updated_at=now()`
	_, err := run.tx.Exec(run.ctx, query, run.tenantID, run.variantID,
		result.event.warehouse, result.qtyDelta)
	return err
}
