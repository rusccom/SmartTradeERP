package media

import (
	"context"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/validation"
)

func (s *Service) Delete(ctx context.Context, tenantID, ownerType, ownerID, mediaID string) error {
	if s.objects == nil {
		return ErrStorageNotConfigured
	}
	if !validScope(tenantID, ownerType, ownerID) || !validation.UUID(mediaID) {
		return ErrInvalidMedia
	}
	objectKey, err := s.deleteRow(ctx, tenantID, ownerType, ownerID, mediaID)
	if err != nil {
		return err
	}
	_ = s.objects.Delete(ctx, objectKey)
	return nil
}

func (s *Service) deleteRow(ctx context.Context, tenantID, ownerType, ownerID, mediaID string) (string, error) {
	objectKey := ""
	err := s.store.WithTx(ctx, func(tx pgx.Tx) error {
		if err := s.repo.LockOwner(ctx, tx, tenantID, ownerType, ownerID); err != nil {
			return err
		}
		key, wasPrimary, err := s.repo.DeleteReady(ctx, tx, tenantID, ownerType, ownerID, mediaID)
		if err != nil {
			return err
		}
		objectKey = key
		if wasPrimary {
			return s.repo.PromoteNextPrimary(ctx, tx, tenantID, ownerType, ownerID)
		}
		return nil
	})
	return objectKey, err
}

func (s *Service) SetPrimary(ctx context.Context, tenantID, ownerType, ownerID, mediaID string) error {
	if s.objects == nil {
		return ErrStorageNotConfigured
	}
	if !validScope(tenantID, ownerType, ownerID) || !validation.UUID(mediaID) {
		return ErrInvalidMedia
	}
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		if err := s.repo.LockOwner(ctx, tx, tenantID, ownerType, ownerID); err != nil {
			return err
		}
		if err := s.repo.ClearPrimary(ctx, tx, tenantID, ownerType, ownerID); err != nil {
			return err
		}
		return s.repo.SetPrimaryByID(ctx, tx, tenantID, ownerType, ownerID, mediaID)
	})
}
