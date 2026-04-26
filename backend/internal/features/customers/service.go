package customers

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/validation"
)

var (
	ErrIsDefault    = errors.New("cannot delete default customer")
	ErrHasDocuments = errors.New("customer has documents")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, tenantID string, query httpx.ListQuery) ([]Customer, int, error) {
	return s.repo.List(ctx, tenantID, query)
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (string, error) {
	req = normalizeCreate(req)
	if err := validateCustomer(req.Name, req.Email); err != nil {
		return "", err
	}
	id := uuid.NewString()
	return id, s.repo.Create(ctx, tenantID, id, req)
}

func (s *Service) ByID(ctx context.Context, tenantID, id string) (Customer, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *Service) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
	req = normalizeUpdate(req)
	if err := validateCustomer(req.Name, req.Email); err != nil {
		return err
	}
	return s.repo.Update(ctx, tenantID, id, req)
}

func (s *Service) Delete(ctx context.Context, tenantID, id string) error {
	isDefault, err := s.repo.IsDefault(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if isDefault {
		return ErrIsDefault
	}
	hasDocs, err := s.repo.HasDocuments(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if hasDocs {
		return ErrHasDocuments
	}
	return s.repo.Delete(ctx, tenantID, id)
}

func normalizeCreate(req CreateRequest) CreateRequest {
	req.Name = validation.Clean(req.Name)
	req.Phone = validation.Clean(req.Phone)
	req.Email = validation.Clean(req.Email)
	return req
}

func normalizeUpdate(req UpdateRequest) UpdateRequest {
	req.Name = validation.Clean(req.Name)
	req.Phone = validation.Clean(req.Phone)
	req.Email = validation.Clean(req.Email)
	return req
}

func validateCustomer(name string, email string) error {
	if !validation.Required(name) || !validation.Max(name, 200) {
		return validation.ErrInvalidData
	}
	if !validation.Email(email) || !validation.Max(email, 255) {
		return validation.ErrInvalidData
	}
	return nil
}
