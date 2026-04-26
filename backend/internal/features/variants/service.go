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
    if count == 0 {
        return pgx.ErrNoRows
    }
    if count == 1 {
        return ErrLastVariant
    }
    return nil
}

func (s *Service) Components(ctx context.Context, tenantID, variantID string) ([]Component, error) {
    return s.repo.Components(ctx, tenantID, variantID)
}

func (s *Service) SetComponents(ctx context.Context, tenantID, variantID string, components []Component) error {
	components = normalizeComponents(components)
	if err := validateComponents(variantID, components); err != nil {
		return err
	}
	if err := s.validateComponentState(ctx, tenantID, variantID, components); err != nil {
		return err
	}
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.ReplaceComponents(ctx, tx, variantID, components)
	})
}

func (s *Service) validateComponentState(
	ctx context.Context,
	tenantID string,
	variantID string,
	components []Component,
) error {
	isComposite, err := s.repo.VariantComposite(ctx, tenantID, variantID)
	if err != nil {
		return err
	}
	if !validComponentState(isComposite, components) {
		return ErrInvalidComponentState
	}
	return s.ensureComponentsBelong(ctx, tenantID, components)
}

func (s *Service) ensureComponentsBelong(ctx context.Context, tenantID string, components []Component) error {
	for _, component := range components {
		exists, err := s.repo.VariantExists(ctx, tenantID, component.ComponentVariantID)
		if err != nil {
			return err
		}
		if !exists {
			return validation.ErrInvalidData
		}
	}
	return nil
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

func normalizeComponents(items []Component) []Component {
    for index := range items {
        items[index].ComponentVariantID = validation.Clean(items[index].ComponentVariantID)
    }
    return items
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

func validateComponents(variantID string, items []Component) error {
    if !validation.UUID(variantID) {
        return validation.ErrInvalidData
    }
    seen := make(map[string]bool, len(items))
    for _, item := range items {
        if invalidComponent(variantID, item, seen) {
            return validation.ErrInvalidData
        }
        seen[item.ComponentVariantID] = true
    }
    return nil
}

func invalidComponent(variantID string, item Component, seen map[string]bool) bool {
    if !validation.Required(item.ComponentVariantID) || !validation.UUID(item.ComponentVariantID) {
        return true
    }
    if !validation.Positive(item.Qty) {
        return true
    }
    return item.ComponentVariantID == variantID || seen[item.ComponentVariantID]
}
