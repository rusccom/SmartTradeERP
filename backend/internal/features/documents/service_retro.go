package documents

import (
	"context"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/features/ledger"
)

type retroState struct {
	affected          []ledger.VariantSequence
	supersededBatchID string
}

type repostInput struct {
	tenantID          string
	documentID        string
	req               UpdateRequest
	supersededBatchID string
}

func (s *Service) retroUpdate(ctx context.Context, tenantID, id string, req UpdateRequest) error {
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		return s.retroUpdateTx(ctx, tx, tenantID, id, req)
	})
}

func (s *Service) retroUpdateTx(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	id string,
	req UpdateRequest,
) error {
	state, err := s.prepareRetroUpdate(ctx, tx, tenantID, id)
	if err != nil {
		return err
	}
	input := repostInput{tenantID: tenantID, documentID: id, req: req, supersededBatchID: state.supersededBatchID}
	newAffected, err := s.repostUpdatedDocument(ctx, tx, input)
	if err != nil {
		return err
	}
	affected := ledger.MergeAffected(state.affected, newAffected)
	return s.recalculateAffected(ctx, tx, tenantID, affected)
}

func (s *Service) prepareRetroUpdate(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	id string,
) (retroState, error) {
	batchID, err := s.ledger.ActiveBatchID(ctx, tx, tenantID, id)
	if err != nil {
		return retroState{}, err
	}
	affected, err := s.ledger.SupersedeDocument(ctx, tx, tenantID, id)
	if err != nil {
		return retroState{}, err
	}
	state := retroState{affected: affected, supersededBatchID: batchID}
	return state, s.repo.SetStatus(ctx, tx, tenantID, id, "draft")
}

func (s *Service) repostUpdatedDocument(
	ctx context.Context,
	tx pgx.Tx,
	input repostInput,
) ([]ledger.VariantSequence, error) {
	if err := s.updateDraftTx(ctx, tx, input.tenantID, input.documentID, input.req); err != nil {
		return nil, err
	}
	return s.postDocumentVersionTx(ctx, tx, input.tenantID, input.documentID, input.supersededBatchID)
}
