package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const financialsDeleteSQL = `DELETE FROM ledger.document_item_financials
        WHERE tenant_id=$1 AND document_item_id IN (
            SELECT DISTINCT m.document_item_id
            FROM ledger.inventory_movements m
            WHERE m.tenant_id=$1 AND m.variant_id = ANY($2::uuid[])
        )`

const financialsInsertSQL = `INSERT INTO ledger.document_item_financials
        (tenant_id, document_id, document_item_id, revenue_amount, cogs_amount,
         gross_profit, calculation_version)
        SELECT m.tenant_id, m.document_id, m.document_item_id,
            COALESCE(SUM(COALESCE(r.revenue_amount,0)),0),
            COALESCE(SUM(COALESCE(r.cogs_amount,0)),0),
            COALESCE(SUM(COALESCE(r.gross_profit,0)),0), $3
        FROM ledger.inventory_movements m
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        JOIN ledger.cost_movement_results r ON r.movement_id=m.id
        WHERE m.tenant_id=$1 AND b.status='active'
            AND m.document_item_id IN (
                SELECT DISTINCT m2.document_item_id
                FROM ledger.inventory_movements m2
                WHERE m2.tenant_id=$1 AND m2.variant_id = ANY($2::uuid[])
            )
        GROUP BY m.tenant_id, m.document_id, m.document_item_id
        ON CONFLICT (tenant_id, document_item_id) DO UPDATE SET
            revenue_amount=EXCLUDED.revenue_amount,
            cogs_amount=EXCLUDED.cogs_amount,
            gross_profit=EXCLUDED.gross_profit,
            calculation_version=EXCLUDED.calculation_version,
            updated_at=now()`

const dailyMetricsDeleteSQL = `DELETE FROM ledger.daily_variant_metrics
        WHERE tenant_id=$1 AND variant_id = ANY($2::uuid[])`

const dailyMetricsInsertSQL = `WITH rows AS (
            SELECT m.tenant_id, m.movement_date, m.variant_id, m.warehouse_id,
                r.qty_delta, r.revenue_amount, r.cogs_amount,
                r.gross_profit, r.running_qty, r.running_avg_cost, r.sequence_num
            FROM ledger.inventory_movements m
            JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
            JOIN ledger.cost_movement_results r ON r.movement_id=m.id
            WHERE m.tenant_id=$1 AND b.status='active'
                AND m.variant_id = ANY($2::uuid[])
        ), latest AS (
            SELECT DISTINCT ON (tenant_id, movement_date, variant_id)
                tenant_id, movement_date, variant_id, running_qty, running_avg_cost
            FROM rows
            ORDER BY tenant_id, movement_date, variant_id, sequence_num DESC
        )
        INSERT INTO ledger.daily_variant_metrics
            (tenant_id, date, variant_id, warehouse_id, qty_in, qty_out,
             revenue_amount, cogs_amount, gross_profit, ending_qty, ending_avg_cost)
        SELECT r.tenant_id, r.movement_date, r.variant_id, r.warehouse_id,
            COALESCE(SUM(CASE WHEN r.qty_delta > 0 THEN r.qty_delta ELSE 0 END),0),
            COALESCE(SUM(CASE WHEN r.qty_delta < 0 THEN -r.qty_delta ELSE 0 END),0),
            COALESCE(SUM(COALESCE(r.revenue_amount,0)),0),
            COALESCE(SUM(COALESCE(r.cogs_amount,0)),0),
            COALESCE(SUM(COALESCE(r.gross_profit,0)),0),
            l.running_qty, l.running_avg_cost
        FROM rows r
        JOIN latest l ON l.tenant_id=r.tenant_id
            AND l.movement_date=r.movement_date AND l.variant_id=r.variant_id
        GROUP BY r.tenant_id, r.movement_date, r.variant_id, r.warehouse_id,
            l.running_qty, l.running_avg_cost`

func (s *Service) refreshVariantProjections(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	variantIDs []string,
) error {
	if err := s.refreshDocumentFinancials(ctx, tx, tenantID, variantIDs); err != nil {
		return err
	}
	return s.refreshDailyMetrics(ctx, tx, tenantID, variantIDs)
}

func (s *Service) refreshDocumentFinancials(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	variantIDs []string,
) error {
	if _, err := tx.Exec(ctx, financialsDeleteSQL, tenantID, variantIDs); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, financialsInsertSQL, tenantID, variantIDs, uuid.NewString())
	return err
}

func (s *Service) refreshDailyMetrics(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	variantIDs []string,
) error {
	if _, err := tx.Exec(ctx, dailyMetricsDeleteSQL, tenantID, variantIDs); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, dailyMetricsInsertSQL, tenantID, variantIDs)
	return err
}
