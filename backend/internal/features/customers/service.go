package customers

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"smarterp/backend/internal/shared/httpx"
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
	id := uuid.NewString()
	return id, s.repo.Create(ctx, tenantID, id, req)
}

func (s *Service) ByID(ctx context.Context, tenantID, id string) (Customer, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *Service) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
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
