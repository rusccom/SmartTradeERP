package documents

import (
	"context"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/features/ledger"
)

type postingRun struct {
	ctx      context.Context
	tx       pgx.Tx
	tenantID string
	batchID  string
	doc      Document
}

func (s *Service) newPostingRun(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	doc Document,
	supersedes string,
) (postingRun, error) {
	input := ledger.BatchInput{
		TenantID:          tenantID,
		DocumentID:        doc.ID,
		EffectiveDate:     mustDate(doc.Date),
		SupersedesBatchID: supersedes,
		Reason:            "document_posting",
	}
	batchID, err := s.ledger.BeginPosting(ctx, tx, input)
	run := postingRun{ctx: ctx, tx: tx, tenantID: tenantID, batchID: batchID, doc: doc}
	return run, err
}

func (s *Service) appendEntries(run postingRun, entries []ledger.EntryInput) ([]ledger.VariantSequence, error) {
	for _, entry := range entries {
		if _, err := s.ledger.Append(run.ctx, run.tx, run.batchID, entry); err != nil {
			return nil, err
		}
	}
	return ledger.AffectedFromEntries(entries), nil
}
