package documents

import (
	"context"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/features/ledger"
)

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
	affected, err := s.prepareRetroUpdate(ctx, tx, tenantID, id)
	if err != nil {
		return err
	}
	if err := s.repostUpdatedDocument(ctx, tx, tenantID, id, req); err != nil {
		return err
	}
	return s.recalculateAffected(ctx, tx, tenantID, affected)
}

func (s *Service) prepareRetroUpdate(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	id string,
) ([]ledger.VariantSequence, error) {
	affected, err := s.ledger.DeleteForDocument(ctx, tx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.DeleteItemComponentsByDocument(ctx, tx, id); err != nil {
		return nil, err
	}
	return affected, s.repo.SetStatus(ctx, tx, tenantID, id, "draft")
}

func (s *Service) repostUpdatedDocument(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	id string,
	req UpdateRequest,
) error {
	if err := s.updateDraftTx(ctx, tx, tenantID, id, req); err != nil {
		return err
	}
	return s.postDocumentTx(ctx, tx, tenantID, id)
}
