package currencies

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/validation"
)

type Service struct {
	store *db.Store
	repo  *Repository
}

func NewService(store *db.Store, repo *Repository) *Service {
	return &Service{store: store, repo: repo}
}

func (s *Service) List(ctx context.Context, tenantID string, page, perPage int) ([]Currency, int, error) {
	return s.repo.List(ctx, tenantID, page, perPage)
}

func (s *Service) Options(ctx context.Context, page, perPage int) ([]CurrencyOption, int, error) {
	return s.repo.Options(ctx, page, perPage)
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (string, error) {
	req = normalizeRequest(req)
	if err := validateRequest(req); err != nil {
		return "", err
	}
	id := uuid.NewString()
	err := s.store.WithTx(ctx, func(tx pgx.Tx) error {
		return s.insertCurrency(ctx, tx, tenantID, id, req)
	})
	if err != nil {
		return "", mapCurrencyWriteError(err)
	}
	return id, nil
}

func (s *Service) SetBase(ctx context.Context, tenantID string, req BaseRequest) error {
	req = normalizeBaseRequest(req)
	if err := validateBaseRequest(req); err != nil {
		return err
	}
	err := s.store.WithTx(ctx, func(tx pgx.Tx) error {
		return s.setBaseCurrency(ctx, tx, tenantID, req)
	})
	if err != nil {
		return mapCurrencyWriteError(err)
	}
	return nil
}

func (s *Service) insertCurrency(ctx context.Context, tx pgx.Tx, tenantID, id string, req CreateRequest) error {
	count, err := s.repo.Count(ctx, tx, tenantID)
	if err != nil {
		return err
	}
	req.IsBase = req.IsBase || count == 0
	if !req.IsBase {
		return s.repo.Create(ctx, tx, tenantID, id, req)
	}
	return s.createBase(ctx, tx, tenantID, id, req)
}

func (s *Service) setBaseCurrency(ctx context.Context, tx pgx.Tx, tenantID string, req BaseRequest) error {
	if err := s.repo.ClearBase(ctx, tx, tenantID); err != nil {
		return err
	}
	id := uuid.NewString()
	if err := s.repo.UpsertBase(ctx, tx, tenantID, id, req); err != nil {
		return err
	}
	return s.repo.SaveBaseSetting(ctx, tx, tenantID, req.CurrencyID)
}

func (s *Service) createBase(ctx context.Context, tx pgx.Tx, tenantID, id string, req CreateRequest) error {
	if err := s.repo.ClearBase(ctx, tx, tenantID); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, tx, tenantID, id, req); err != nil {
		return err
	}
	return s.repo.SaveBaseSetting(ctx, tx, tenantID, req.CurrencyID)
}

func normalizeRequest(req CreateRequest) CreateRequest {
	req.CurrencyID = validation.Clean(req.CurrencyID)
	req.DisplaySymbol = validation.Clean(req.DisplaySymbol)
	return req
}

func normalizeBaseRequest(req BaseRequest) BaseRequest {
	req.CurrencyID = validation.Clean(req.CurrencyID)
	req.DisplaySymbol = validation.Clean(req.DisplaySymbol)
	return req
}

func validateRequest(req CreateRequest) error {
	if !validation.UUID(req.CurrencyID) {
		return validation.ErrInvalidData
	}
	if !validation.Max(req.DisplaySymbol, 8) {
		return validation.ErrInvalidData
	}
	return nil
}

func validateBaseRequest(req BaseRequest) error {
	if !validation.UUID(req.CurrencyID) {
		return validation.ErrInvalidData
	}
	if !validation.Max(req.DisplaySymbol, 8) {
		return validation.ErrInvalidData
	}
	return nil
}
