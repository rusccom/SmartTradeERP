package products

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/features/ledger"
    "smarterp/backend/internal/shared/db"
)

var ErrHasMovements = errors.New("product has movements")

type Service struct {
    store  *db.Store
    repo   *Repository
    ledger *ledger.Service
}

func NewService(store *db.Store, repo *Repository, ledger *ledger.Service) *Service {
    return &Service{store: store, repo: repo, ledger: ledger}
}

func (s *Service) List(ctx context.Context, tenantID, isComposite string, page, perPage int) ([]Product, int, error) {
    return s.repo.List(ctx, tenantID, isComposite, page, perPage)
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (string, error) {
    productID := uuid.NewString()
    variantID := uuid.NewString()
    err := s.store.WithTx(ctx, func(tx pgx.Tx) error {
        return s.createWithDefaultVariant(ctx, tx, tenantID, productID, variantID, req)
    })
    if err != nil {
        return "", err
    }
    return productID, nil
}

func (s *Service) createWithDefaultVariant(
    ctx context.Context,
    tx pgx.Tx,
    tenantID string,
    productID string,
    variantID string,
    req CreateRequest,
) error {
    if err := s.repo.Create(ctx, tx, tenantID, productID, req); err != nil {
        return err
    }
    return s.repo.CreateDefaultVariant(ctx, tx, variantID, productID, req)
}

func (s *Service) ByID(ctx context.Context, tenantID, id string) (Product, error) {
    return s.repo.GetByID(ctx, tenantID, id)
}

func (s *Service) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
    return s.repo.Update(ctx, tenantID, id, req)
}

func (s *Service) Delete(ctx context.Context, tenantID, id string) error {
    hasMovements, err := s.ledger.HasProductMovements(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if hasMovements {
        return ErrHasMovements
    }
    return s.repo.Delete(ctx, tenantID, id)
}
