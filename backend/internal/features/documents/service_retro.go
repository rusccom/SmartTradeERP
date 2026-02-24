package documents

import (
    "context"

    "github.com/jackc/pgx/v5"
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
    affected, err := s.ledger.DeleteForDocument(ctx, tx, tenantID, id)
    if err != nil {
        return err
    }
    if err := s.repo.DeleteItemComponentsByDocument(ctx, tx, id); err != nil {
        return err
    }
    if err := s.updateDraftTx(ctx, tx, tenantID, id, req); err != nil {
        return err
    }
    if err := s.postDocumentTx(ctx, tx, tenantID, id); err != nil {
        return err
    }
    return s.recalculateAffected(ctx, tx, tenantID, affected)
}
