package documents

import (
    "context"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/features/ledger"
    "smarterp/backend/internal/shared/db"
)

type Service struct {
    store  *db.Store
    repo   *Repository
    ledger *ledger.Service
}

func NewService(store *db.Store, repo *Repository, ledger *ledger.Service) *Service {
    return &Service{store: store, repo: repo, ledger: ledger}
}

func (s *Service) List(ctx context.Context, tenantID string, filters Filters, page, perPage int) ([]ListItem, int, error) {
    return s.repo.List(ctx, tenantID, filters, page, perPage)
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (string, error) {
    id := uuid.NewString()
    err := s.store.WithTx(ctx, func(tx pgx.Tx) error {
        return s.createDraftTx(ctx, tx, tenantID, id, req)
    })
    if err != nil {
        return "", err
    }
    return id, nil
}

func (s *Service) createDraftTx(
    ctx context.Context,
    tx pgx.Tx,
    tenantID string,
    documentID string,
    req CreateRequest,
) error {
    if err := s.repo.InsertDocument(ctx, tx, tenantID, documentID, req); err != nil {
        return err
    }
    return s.repo.ReplaceItems(ctx, tx, documentID, req.Items)
}

func (s *Service) ByID(ctx context.Context, tenantID, id string) (Document, error) {
    doc, err := s.repo.ByID(ctx, tenantID, id)
    if err != nil {
        return Document{}, err
    }
    items, total, err := s.repo.LoadItemsWithProfit(ctx, tenantID, id)
    if err != nil {
        return Document{}, err
    }
    doc.Items = items
    doc.TotalProfit = total
    return doc, nil
}

func (s *Service) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
    status, err := s.repo.Status(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if status == "draft" {
        return s.updateDraft(ctx, tenantID, id, req)
    }
    if status == "posted" {
        return s.retroUpdate(ctx, tenantID, id, req)
    }
    return ErrStatusConflict
}

func (s *Service) updateDraft(ctx context.Context, tenantID, id string, req UpdateRequest) error {
    return s.store.WithTx(ctx, func(tx pgx.Tx) error {
        return s.updateDraftTx(ctx, tx, tenantID, id, req)
    })
}

func (s *Service) updateDraftTx(ctx context.Context, tx pgx.Tx, tenantID, id string, req UpdateRequest) error {
    if err := s.repo.UpdateDocument(ctx, tx, tenantID, id, req); err != nil {
        return err
    }
    return s.repo.ReplaceItems(ctx, tx, id, req.Items)
}

func (s *Service) Post(ctx context.Context, tenantID, id string) error {
    status, err := s.repo.Status(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if status != "draft" {
        return ErrDraftOnly
    }
    return s.store.WithTx(ctx, func(tx pgx.Tx) error {
        return s.postDocumentTx(ctx, tx, tenantID, id)
    })
}

func (s *Service) Cancel(ctx context.Context, tenantID, id string) error {
    status, err := s.repo.Status(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if status != "posted" {
        return ErrPostedOnly
    }
    return s.store.WithTx(ctx, func(tx pgx.Tx) error {
        return s.cancelTx(ctx, tx, tenantID, id)
    })
}

func (s *Service) cancelTx(ctx context.Context, tx pgx.Tx, tenantID, id string) error {
    affected, err := s.ledger.DeleteForDocument(ctx, tx, tenantID, id)
    if err != nil {
        return err
    }
    if err := s.repo.DeleteItemComponentsByDocument(ctx, tx, id); err != nil {
        return err
    }
    if err := s.recalculateAffected(ctx, tx, tenantID, affected); err != nil {
        return err
    }
    return s.repo.SetStatus(ctx, tx, tenantID, id, "cancelled")
}

func (s *Service) Delete(ctx context.Context, tenantID, id string) error {
    status, err := s.repo.Status(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if status != "draft" {
        return ErrDraftOnly
    }
    return s.repo.DeleteDraft(ctx, tenantID, id)
}

func (s *Service) recalculateAffected(
    ctx context.Context,
    tx pgx.Tx,
    tenantID string,
    affected []ledger.VariantSequence,
) error {
    for _, item := range affected {
        if err := s.ledger.Recalculate(ctx, tx, tenantID, item.VariantID, item.Earliest); err != nil {
            return err
        }
    }
    return nil
}
