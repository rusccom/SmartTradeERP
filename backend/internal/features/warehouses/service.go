package warehouses

import (
    "context"

    "github.com/google/uuid"

    "smarterp/backend/internal/features/ledger"
    "smarterp/backend/internal/shared/validation"
)

type Service struct {
    repo   *Repository
    ledger *ledger.Service
}

func NewService(repo *Repository, ledger *ledger.Service) *Service {
    return &Service{repo: repo, ledger: ledger}
}

func (s *Service) List(ctx context.Context, tenantID string) ([]Warehouse, error) {
    return s.repo.List(ctx, tenantID)
}

func (s *Service) ListWithIncludes(
	ctx context.Context,
	tenantID string,
	include WarehouseListInclude,
) ([]WarehouseListItem, error) {
	return s.repo.ListWithIncludes(ctx, tenantID, include)
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (string, error) {
    req = normalizeCreate(req)
    if err := validateWarehouse(req.Name); err != nil {
        return "", err
    }
    id := uuid.NewString()
    req = applyDefaultActive(req)
    req, err := s.ensureDefaultOnCreate(ctx, tenantID, req)
    if err != nil {
        return "", err
    }
    if err := s.repo.Create(ctx, tenantID, id, req); err != nil {
        return "", err
    }
    return id, nil
}

func applyDefaultActive(req CreateRequest) CreateRequest {
    if !req.IsActive {
        req.IsActive = true
    }
    return req
}

func (s *Service) ensureDefaultOnCreate(ctx context.Context, tenantID string, req CreateRequest) (CreateRequest, error) {
    defaults, err := s.repo.CountDefaults(ctx, tenantID)
    if err != nil {
        return CreateRequest{}, err
    }
    if defaults == 0 {
        req.IsDefault = true
    }
    return req, nil
}

func (s *Service) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
    req = normalizeUpdate(req)
    if err := validateWarehouse(req.Name); err != nil {
        return err
    }
    if err := s.ensureDefaultOnUpdate(ctx, tenantID, id, req.IsDefault); err != nil {
        return err
    }
    return s.repo.Update(ctx, tenantID, id, req)
}

func (s *Service) ensureDefaultOnUpdate(ctx context.Context, tenantID, id string, nextDefault bool) error {
    if nextDefault {
        return nil
    }
    isDefault, err := s.repo.IsDefault(ctx, tenantID, id)
    if err != nil || !isDefault {
        return err
    }
    defaults, err := s.repo.CountDefaults(ctx, tenantID)
    if err != nil {
        return err
    }
    if defaults <= 1 {
        return ErrMustKeepDefault
    }
    return nil
}

func normalizeCreate(req CreateRequest) CreateRequest {
    req.Name = validation.Clean(req.Name)
    req.Address = validation.Clean(req.Address)
    return req
}

func normalizeUpdate(req UpdateRequest) UpdateRequest {
    req.Name = validation.Clean(req.Name)
    req.Address = validation.Clean(req.Address)
    return req
}

func validateWarehouse(name string) error {
    if !validation.Required(name) || !validation.Max(name, 200) {
        return validation.ErrInvalidData
    }
    return nil
}

func (s *Service) Delete(ctx context.Context, tenantID, id string) error {
    hasMovements, err := s.ledger.HasWarehouseMovements(ctx, tenantID, id)
    if err != nil {
        return err
    }
    if hasMovements {
        return ErrHasMovements
    }
    if err := s.ensureDefaultOnDelete(ctx, tenantID, id); err != nil {
        return err
    }
    return s.repo.Delete(ctx, tenantID, id)
}

func (s *Service) ensureDefaultOnDelete(ctx context.Context, tenantID, id string) error {
    isDefault, err := s.repo.IsDefault(ctx, tenantID, id)
    if err != nil || !isDefault {
        return err
    }
    defaults, err := s.repo.CountDefaults(ctx, tenantID)
    if err != nil {
        return err
    }
    if defaults <= 1 {
        return ErrMustKeepDefault
    }
    return nil
}
