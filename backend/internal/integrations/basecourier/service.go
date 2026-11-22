package basecourier

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type service struct {
	repo   repository
	client *http.Client
	vault  vault
}

func (s service) save(ctx context.Context, tenantID string, u update) (*connection, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}
	sealed, err := s.vault.seal(tenantID, u.Email, u.APIKey)
	if err != nil {
		return nil, err
	}
	if err = verifyProfile(ctx, s.client, u.Email, u.APIKey); err != nil {
		return nil, err
	}
	return s.repo.save(ctx, tenantID, u, sealed)
}

func (s service) check(ctx context.Context, tenantID string, version uuid.UUID) (*connection, error) {
	email, sealed, err := s.repo.credential(ctx, tenantID, version)
	if err != nil {
		return nil, err
	}
	key, err := s.vault.open(tenantID, email, sealed)
	if err != nil {
		return nil, err
	}
	if err = verifyProfile(ctx, s.client, email, key); err != nil {
		return nil, err
	}
	return s.repo.markVerified(ctx, tenantID, version)
}
