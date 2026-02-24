package variants

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/features/ledger"
    "smarterp/backend/internal/shared/db"
)

var ErrLastVariant = errors.New("product must keep at least one variant")
var ErrInvalidComponentState = errors.New("invalid component state")

type Service struct {
    store  *db.Store
    repo   *Repository
    ledger *ledger.Service
}

func NewService(store *db.Store, repo *Repository, ledger *ledger.Service) *Service {
    return &Service{store: store, repo: repo, ledger: ledger}
}

func (s *Service) List(ctx context.Context, tenantID, productID string, page, perPage int) ([]Variant, int, error) {
    return s.repo.List(ctx, tenantID, productID, page, perPage)
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (string, error) {
    id := uuid.NewString()
    if err := s.repo.Create(ctx, tenantID, id, req); err != nil {
        return "", err
    }
    return id, nil
}

func (s *Service) ByID(ctx context.Context, tenantID, id string) (Variant, error) {
    return s.repo.ByID(ctx, tenantID, id)
}

func (s *Service) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
    return s.repo.Update(ctx, tenantID, id, req)
}

func (s *Service) Delete(ctx context.Context, tenantID, id string) error {
    blocked, err := s.ledger.HasVariantMovements(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if blocked {
        return ErrHasMovements
    }
    if err := s.ensureNotLastVariant(ctx, tenantID, id); err != nil {
        return err
    }
    return s.repo.Delete(ctx, tenantID, id)
}

func (s *Service) ensureNotLastVariant(ctx context.Context, tenantID, id string) error {
    count, err := s.repo.ProductVariantCount(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if count <= 1 {
        return ErrLastVariant
    }
    return nil
}

func (s *Service) Components(ctx context.Context, tenantID, variantID string) ([]Component, error) {
    return s.repo.Components(ctx, tenantID, variantID)
}

func (s *Service) SetComponents(ctx context.Context, tenantID, variantID string, components []Component) error {
    isComposite, err := s.repo.VariantComposite(ctx, tenantID, variantID)
    if err != nil {
        return err
    }
    if !validComponentState(isComposite, components) {
        return ErrInvalidComponentState
    }
    return s.store.WithTx(ctx, func(tx pgx.Tx) error {
        return s.repo.ReplaceComponents(ctx, tx, variantID, components)
    })
}

func validComponentState(isComposite bool, components []Component) bool {
    if isComposite && len(components) == 0 {
        return false
    }
    if !isComposite && len(components) > 0 {
        return false
    }
    return true
}

func (s *Service) Stock(ctx context.Context, tenantID, variantID string) (Stock, error) {
    qty, avg, err := s.ledger.GlobalStock(ctx, tenantID, variantID)
    if err != nil {
        return Stock{}, err
    }
    byWarehouse, err := s.repo.WarehouseStock(ctx, tenantID, variantID)
    if err != nil {
        return Stock{}, err
    }
    stock := Stock{GlobalQty: qty, RunningAvg: avg, WarehousesStock: byWarehouse}
    return stock, nil
}
