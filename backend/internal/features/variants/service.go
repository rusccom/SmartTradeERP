package variants

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/shopspring/decimal"

    "smarterp/backend/internal/features/ledger"
    "smarterp/backend/internal/shared/db"
    "smarterp/backend/internal/shared/validation"
)

var ErrLastVariant = errors.New("product must keep at least one variant")
var ErrUsedInBundle = errors.New("variant used in bundle")

type Service struct {
    store       *db.Store
    repo        *Repository
    ledger      *ledger.Service
    bundleState BundleStateReader
}

type BundleStateReader interface {
    VariantHasComponents(ctx context.Context, tenantID string, variantID string) (bool, error)
    VariantUsedAsComponent(ctx context.Context, tenantID string, variantID string) (bool, error)
}

func NewService(
    store *db.Store,
    repo *Repository,
    ledger *ledger.Service,
    bundleState BundleStateReader,
) *Service {
    return &Service{store: store, repo: repo, ledger: ledger, bundleState: bundleState}
}

func (s *Service) List(ctx context.Context, tenantID, productID string, page, perPage int) ([]Variant, int, error) {
    return s.repo.List(ctx, tenantID, productID, page, perPage)
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (string, error) {
    req = normalizeCreate(req)
    if err := validateCreate(req); err != nil {
        return "", err
    }
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
    req = normalizeUpdate(req)
    if err := validateUpdate(req); err != nil {
        return err
    }
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
    if err := s.ensureVariantNotInBundle(ctx, tenantID, id); err != nil {
        return err
    }
    if err := s.ensureNotLastVariant(ctx, tenantID, id); err != nil {
        return err
    }
    return s.repo.Delete(ctx, tenantID, id)
}

func (s *Service) ensureVariantNotInBundle(ctx context.Context, tenantID, id string) error {
    linked, err := s.variantLinkedToBundle(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if linked {
        return ErrUsedInBundle
    }
    return nil
}

func (s *Service) variantLinkedToBundle(ctx context.Context, tenantID, id string) (bool, error) {
    hasComponents, err := s.bundleState.VariantHasComponents(ctx, tenantID, id)
    if err != nil || hasComponents {
        return hasComponents, err
    }
    return s.bundleState.VariantUsedAsComponent(ctx, tenantID, id)
}

func (s *Service) ensureNotLastVariant(ctx context.Context, tenantID, id string) error {
    count, err := s.repo.ProductVariantCount(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if count == 0 {
        return pgx.ErrNoRows
    }
    if count == 1 {
        return ErrLastVariant
    }
    return nil
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

func normalizeCreate(req CreateRequest) CreateRequest {
    req.ProductID = validation.Clean(req.ProductID)
    req.Name = validation.Clean(req.Name)
    req.SKUCode = validation.Clean(req.SKUCode)
    req.Barcode = validation.Clean(req.Barcode)
    req.Unit = validation.Clean(req.Unit)
    req.Option1 = validation.Clean(req.Option1)
    req.Option2 = validation.Clean(req.Option2)
    req.Option3 = validation.Clean(req.Option3)
    return req
}

func normalizeUpdate(req UpdateRequest) UpdateRequest {
    req.Name = validation.Clean(req.Name)
    req.SKUCode = validation.Clean(req.SKUCode)
    req.Barcode = validation.Clean(req.Barcode)
    req.Unit = validation.Clean(req.Unit)
    req.Option1 = validation.Clean(req.Option1)
    req.Option2 = validation.Clean(req.Option2)
    req.Option3 = validation.Clean(req.Option3)
    return req
}

func validateCreate(req CreateRequest) error {
    if !validation.Required(req.ProductID) || !validation.UUID(req.ProductID) {
        return validation.ErrInvalidData
    }
    return validateVariant(req.Name, req.Unit, req.Price)
}

func validateUpdate(req UpdateRequest) error {
    return validateVariant(req.Name, req.Unit, req.Price)
}

func validateVariant(name string, unit string, price decimal.Decimal) error {
    if !validation.Required(name) || !validation.Required(unit) {
        return validation.ErrInvalidData
    }
    if !validation.NonNegative(price) || !validation.Max(unit, 24) {
        return validation.ErrInvalidData
    }
    return nil
}
