package admin

import (
	"context"

	"smarterp/backend/internal/shared/auth"
)

type Service struct {
	repo   *Repository
	tokens *auth.TokenService
}

func NewService(repo *Repository, tokens *auth.TokenService) *Service {
	return &Service{repo: repo, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (auth.TokenResponse, error) {
	admin, err := s.repo.FindAdminByEmail(ctx, req.Email)
	if err != nil {
		return auth.TokenResponse{}, err
	}
	if !auth.VerifyPassword(req.Password, admin.PasswordHash) {
		return auth.TokenResponse{}, ErrInvalidCredentials
	}
	return s.tokens.Issue(admin.ID, "", "owner", "admin")
}

func (s *Service) ListTenants(ctx context.Context, page, perPage int) ([]Tenant, int, error) {
	return s.repo.ListTenants(ctx, page, perPage)
}

func (s *Service) TenantByID(ctx context.Context, id string) (Tenant, error) {
	return s.repo.GetTenantByID(ctx, id)
}

func (s *Service) PlatformStats(ctx context.Context) (Stats, error) {
	return s.repo.Stats(ctx)
}
