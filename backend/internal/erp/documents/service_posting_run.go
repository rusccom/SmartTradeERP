package documents

import (
	"context"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/erp/bundles"
	"smarterp/backend/internal/erp/ledger"
)

type postingRun struct {
	ctx       context.Context
	tx        pgx.Tx
	tenantID  string
	batchID   string
	doc       Document
	snapshots map[string][]bundles.Component
}

func (s *Service) newPostingRun(
	ctx context.Context,
	tx pgx.Tx,
	input postingVersionInput,
	doc Document,
) (postingRun, error) {
	batchInput := ledger.BatchInput{
		TenantID:          input.tenantID,
		DocumentID:        doc.ID,
		EffectiveDate:     mustDate(doc.Date),
		SupersedesBatchID: input.supersedes,
		Reason:            "document_posting",
		PostedBy:          actorID(ctx),
	}
	batchID, err := s.ledger.BeginPosting(ctx, tx, batchInput)
	run := postingRun{
		ctx: ctx, tx: tx, tenantID: input.tenantID,
		batchID: batchID, doc: doc, snapshots: input.snapshots,
	}
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
