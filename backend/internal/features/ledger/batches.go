package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) BeginPosting(ctx context.Context, tx pgx.Tx, input BatchInput) (string, error) {
	id := uuid.NewString()
	query := `INSERT INTO ledger.posting_batches
        (id, tenant_id, document_id, effective_date, supersedes_batch_id, reason)
        VALUES ($1,$2,$3,$4,NULLIF($5,'')::uuid,$6)`
	_, err := tx.Exec(ctx, query, id, input.TenantID, input.DocumentID,
		input.EffectiveDate, input.SupersedesBatchID, postingReason(input))
	return id, err
}

func postingReason(input BatchInput) string {
	if input.Reason == "" {
		return "posting"
	}
	return input.Reason
}

func (s *Service) SupersedeDocument(ctx context.Context, tx pgx.Tx, tenantID, documentID string) ([]VariantSequence, error) {
	return s.deactivateDocument(ctx, tx, tenantID, documentID, "superseded")
}

func (s *Service) ReverseDocument(ctx context.Context, tx pgx.Tx, tenantID, documentID string) ([]VariantSequence, error) {
	return s.deactivateDocument(ctx, tx, tenantID, documentID, "reversed")
}

func (s *Service) ActiveBatchID(ctx context.Context, tx pgx.Tx, tenantID, documentID string) (string, error) {
	query := `SELECT id::text FROM ledger.posting_batches
        WHERE tenant_id=$1 AND document_id=$2 AND status='active'
        ORDER BY posted_at DESC LIMIT 1`
	id := ""
	err := tx.QueryRow(ctx, query, tenantID, documentID).Scan(&id)
	if err == pgx.ErrNoRows {
		return "", ErrActiveBatchRequired
	}
	return id, err
}

func (s *Service) deactivateDocument(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	documentID string,
	status string,
) ([]VariantSequence, error) {
	affected, err := s.activeDocumentVariants(ctx, tx, tenantID, documentID)
	if err != nil {
		return nil, err
	}
	return affected, s.updateBatchStatus(ctx, tx, tenantID, documentID, status)
}

func (s *Service) updateBatchStatus(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	documentID string,
	status string,
) error {
	query := `UPDATE ledger.posting_batches
        SET status=$3, updated_at=now()
        WHERE tenant_id=$1 AND document_id=$2 AND status='active'`
	tag, err := tx.Exec(ctx, query, tenantID, documentID, status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrActiveBatchRequired
	}
	return nil
}

func (s *Service) activeDocumentVariants(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	documentID string,
) ([]VariantSequence, error) {
	query := `SELECT DISTINCT m.variant_id::text
        FROM ledger.inventory_movements m
        JOIN ledger.posting_batches b ON b.id=m.posting_batch_id
        WHERE m.tenant_id=$1 AND m.document_id=$2 AND b.status='active'`
	rows, err := tx.Query(ctx, query, tenantID, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVariantIDs(rows)
}

func scanVariantIDs(rows pgx.Rows) ([]VariantSequence, error) {
	items := make([]VariantSequence, 0)
	for rows.Next() {
		item := VariantSequence{}
		if err := rows.Scan(&item.VariantID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
