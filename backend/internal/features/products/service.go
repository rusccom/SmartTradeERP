package products

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/features/ledger"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/validation"
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

func (s *Service) List(ctx context.Context, tenantID string, query ProductListQuery) ([]Product, int, error) {
	return s.repo.List(ctx, tenantID, query)
}

func (s *Service) ListWithIncludes(
	ctx context.Context,
	tenantID string,
	query ProductListQuery,
	include ProductListInclude,
) ([]ProductListItem, int, error) {
	return s.repo.ListWithIncludes(ctx, tenantID, query, include)
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (string, error) {
	req = normalizeCreate(req)
	if err := validateCreate(req); err != nil {
		return "", err
	}
    productID := uuid.NewString()
    variantID := uuid.NewString()
    err := s.store.WithTx(ctx, func(tx pgx.Tx) error {
        input := createProductTx{tenantID: tenantID, productID: productID, variantID: variantID, req: req}
        return s.createWithDefaultVariant(ctx, tx, input)
    })
    if err != nil {
        return "", err
    }
    return productID, nil
}

type createProductTx struct {
	tenantID  string
	productID string
	variantID string
	req       CreateRequest
}

func (s *Service) createWithDefaultVariant(
    ctx context.Context,
    tx pgx.Tx,
    input createProductTx,
) error {
    if err := s.repo.Create(ctx, tx, input.tenantID, input.productID, input.req); err != nil {
        return err
    }
    return s.repo.CreateDefaultVariant(ctx, tx, input)
}

func (s *Service) ByID(ctx context.Context, tenantID, id string) (Product, error) {
    return s.repo.GetByID(ctx, tenantID, id)
}

func (s *Service) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
	req = normalizeUpdate(req)
	if err := validateUpdate(req); err != nil {
		return err
	}
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

func normalizeCreate(req CreateRequest) CreateRequest {
	req.Name = validation.Clean(req.Name)
	req.Unit = validation.Clean(req.Unit)
	req.SKUCode = validation.Clean(req.SKUCode)
	req.Barcode = validation.Clean(req.Barcode)
	return req
}

func normalizeUpdate(req UpdateRequest) UpdateRequest {
	req.Name = validation.Clean(req.Name)
	return req
}

func validateCreate(req CreateRequest) error {
	if validateName(req.Name) != nil || !validation.Required(req.Unit) {
		return validation.ErrInvalidData
	}
	if !validation.NonNegative(req.Price) || !validation.Max(req.Unit, 24) {
		return validation.ErrInvalidData
	}
	return nil
}

func validateUpdate(req UpdateRequest) error {
	return validateName(req.Name)
}

func validateName(name string) error {
	if !validation.Required(name) || !validation.Max(name, 200) {
		return validation.ErrInvalidData
	}
	return nil
}
