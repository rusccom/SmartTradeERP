package variants

import (
	"context"

	mediafeature "smarterp/backend/internal/platform/media"
	"smarterp/backend/internal/shared/validation"
)

const variantMediaOwner = "variant"

type MediaService interface {
	List(ctx context.Context, tenantID string, ownerType string, ownerID string) ([]mediafeature.Item, error)
	CreateDirectUpload(ctx context.Context, req mediafeature.DirectUploadRequest) (mediafeature.DirectUpload, error)
	CompleteDirectUpload(ctx context.Context, req mediafeature.CompleteRequest) (mediafeature.Item, error)
	Delete(ctx context.Context, tenantID string, ownerType string, ownerID string, mediaID string) error
	SetPrimary(ctx context.Context, tenantID string, ownerType string, ownerID string, mediaID string) error
}

func (s *Service) SetMediaService(mediaService MediaService) {
	s.media = mediaService
}

func (s *Service) ListMedia(ctx context.Context, tenantID, variantID string) ([]mediafeature.Item, error) {
	if err := s.ensureVariant(ctx, tenantID, variantID); err != nil {
		return nil, err
	}
	return s.media.List(ctx, tenantID, variantMediaOwner, variantID)
}

func (s *Service) CreateMediaUpload(
	ctx context.Context,
	tenantID string,
	variantID string,
	input mediafeature.DirectUploadInput,
) (mediafeature.DirectUpload, error) {
	if err := s.ensureVariant(ctx, tenantID, variantID); err != nil {
		return mediafeature.DirectUpload{}, err
	}
	req := mediafeature.DirectUploadRequest{
		TenantID: tenantID, OwnerType: variantMediaOwner, OwnerID: variantID, Input: input,
	}
	return s.media.CreateDirectUpload(ctx, req)
}

func (s *Service) CompleteMediaUpload(ctx context.Context, tenantID, variantID, mediaID string) (mediafeature.Item, error) {
	if err := s.ensureVariant(ctx, tenantID, variantID); err != nil {
		return mediafeature.Item{}, err
	}
	req := mediafeature.CompleteRequest{
		TenantID: tenantID, OwnerType: variantMediaOwner, OwnerID: variantID, MediaID: mediaID,
	}
	return s.media.CompleteDirectUpload(ctx, req)
}

func (s *Service) DeleteMedia(ctx context.Context, tenantID, variantID, mediaID string) error {
	if err := s.ensureVariant(ctx, tenantID, variantID); err != nil {
		return err
	}
	return s.media.Delete(ctx, tenantID, variantMediaOwner, variantID, mediaID)
}

func (s *Service) SetPrimaryMedia(ctx context.Context, tenantID, variantID, mediaID string) error {
	if err := s.ensureVariant(ctx, tenantID, variantID); err != nil {
		return err
	}
	return s.media.SetPrimary(ctx, tenantID, variantMediaOwner, variantID, mediaID)
}

func (s *Service) ensureVariant(ctx context.Context, tenantID, variantID string) error {
	if !validation.UUID(variantID) {
		return mediafeature.ErrInvalidMedia
	}
	if s.media == nil {
		return mediafeature.ErrStorageNotConfigured
	}
	if _, err := s.repo.ByID(ctx, tenantID, variantID); err != nil {
		return err
	}
	return nil
}
