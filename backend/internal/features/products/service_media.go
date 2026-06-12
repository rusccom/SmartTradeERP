package products

import (
	"context"

	mediafeature "smarterp/backend/internal/features/media"
	"smarterp/backend/internal/shared/validation"
)

const productMediaOwner = "product"

type MediaService interface {
	List(ctx context.Context, tenantID string, ownerType string, ownerID string) ([]mediafeature.Item, error)
	CreateDirectUpload(ctx context.Context, req mediafeature.DirectUploadRequest) (mediafeature.DirectUpload, error)
	CompleteDirectUpload(ctx context.Context, req mediafeature.CompleteRequest) (mediafeature.Item, error)
}

func (s *Service) SetMediaService(mediaService MediaService) {
	s.media = mediaService
}

func (s *Service) ListMedia(
	ctx context.Context,
	tenantID string,
	productID string,
) ([]mediafeature.Item, error) {
	if !validation.UUID(productID) {
		return nil, mediafeature.ErrInvalidMedia
	}
	if s.media == nil {
		return nil, mediafeature.ErrStorageNotConfigured
	}
	if _, err := s.repo.GetByID(ctx, tenantID, productID); err != nil {
		return nil, err
	}
	return s.media.List(ctx, tenantID, productMediaOwner, productID)
}

func (s *Service) CreateMediaUpload(
	ctx context.Context,
	tenantID string,
	productID string,
	input mediafeature.DirectUploadInput,
) (mediafeature.DirectUpload, error) {
	if !validation.UUID(productID) {
		return mediafeature.DirectUpload{}, mediafeature.ErrInvalidMedia
	}
	if s.media == nil {
		return mediafeature.DirectUpload{}, mediafeature.ErrStorageNotConfigured
	}
	if _, err := s.repo.GetByID(ctx, tenantID, productID); err != nil {
		return mediafeature.DirectUpload{}, err
	}
	req := mediafeature.DirectUploadRequest{
		TenantID: tenantID, OwnerType: productMediaOwner, OwnerID: productID, Input: input,
	}
	return s.media.CreateDirectUpload(ctx, req)
}

func (s *Service) CompleteMediaUpload(
	ctx context.Context,
	tenantID string,
	productID string,
	mediaID string,
) (mediafeature.Item, error) {
	if !validation.UUID(productID) {
		return mediafeature.Item{}, mediafeature.ErrInvalidMedia
	}
	if s.media == nil {
		return mediafeature.Item{}, mediafeature.ErrStorageNotConfigured
	}
	if _, err := s.repo.GetByID(ctx, tenantID, productID); err != nil {
		return mediafeature.Item{}, err
	}
	req := mediafeature.CompleteRequest{
		TenantID: tenantID, OwnerType: productMediaOwner, OwnerID: productID, MediaID: mediaID,
	}
	return s.media.CompleteDirectUpload(ctx, req)
}
