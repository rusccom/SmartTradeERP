package products

import (
	"context"

	mediafeature "smarterp/backend/internal/platform/media"
	"smarterp/backend/internal/shared/validation"
)

const productMediaOwner = "product"

type MediaService interface {
	List(ctx context.Context, tenantID string, ownerType string, ownerID string) ([]mediafeature.Item, error)
	CreateDirectUpload(ctx context.Context, req mediafeature.DirectUploadRequest) (mediafeature.DirectUpload, error)
	CompleteDirectUpload(ctx context.Context, req mediafeature.CompleteRequest) (mediafeature.Item, error)
	Delete(ctx context.Context, tenantID string, ownerType string, ownerID string, mediaID string) error
	SetPrimary(ctx context.Context, tenantID string, ownerType string, ownerID string, mediaID string) error
	PrimaryThumbsByOwners(ctx context.Context, tenantID string, ownerType string, ownerIDs []string) (map[string]string, error)
}

const variantImageOwner = "variant"

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

func (s *Service) DeleteMedia(ctx context.Context, tenantID, productID, mediaID string) error {
	if !validation.UUID(productID) {
		return mediafeature.ErrInvalidMedia
	}
	if s.media == nil {
		return mediafeature.ErrStorageNotConfigured
	}
	if _, err := s.repo.GetByID(ctx, tenantID, productID); err != nil {
		return err
	}
	return s.media.Delete(ctx, tenantID, productMediaOwner, productID, mediaID)
}

func (s *Service) SetPrimaryMedia(ctx context.Context, tenantID, productID, mediaID string) error {
	if !validation.UUID(productID) {
		return mediafeature.ErrInvalidMedia
	}
	if s.media == nil {
		return mediafeature.ErrStorageNotConfigured
	}
	if _, err := s.repo.GetByID(ctx, tenantID, productID); err != nil {
		return err
	}
	return s.media.SetPrimary(ctx, tenantID, productMediaOwner, productID, mediaID)
}

func (s *Service) attachImages(
	ctx context.Context,
	tenantID string,
	items []ProductListItem,
	include ProductListInclude,
) {
	if s.media == nil || len(items) == 0 {
		return
	}
	s.attachProductThumbs(ctx, tenantID, items)
	if include.Variants {
		s.attachVariantThumbs(ctx, tenantID, items)
	}
}

func (s *Service) attachProductThumbs(ctx context.Context, tenantID string, items []ProductListItem) {
	ids := make([]string, len(items))
	for i := range items {
		ids[i] = items[i].ID
	}
	thumbs, err := s.media.PrimaryThumbsByOwners(ctx, tenantID, productMediaOwner, ids)
	if err != nil {
		return
	}
	for i := range items {
		items[i].ImageThumbURL = thumbs[items[i].ID]
	}
}

func (s *Service) attachVariantThumbs(ctx context.Context, tenantID string, items []ProductListItem) {
	ids := variantImageIDs(items)
	if len(ids) == 0 {
		return
	}
	thumbs, err := s.media.PrimaryThumbsByOwners(ctx, tenantID, variantImageOwner, ids)
	if err != nil {
		return
	}
	applyVariantThumbs(items, thumbs)
}

func variantImageIDs(items []ProductListItem) []string {
	ids := make([]string, 0)
	for i := range items {
		for j := range items[i].Variants {
			ids = append(ids, items[i].Variants[j].ID)
		}
	}
	return ids
}

func applyVariantThumbs(items []ProductListItem, thumbs map[string]string) {
	for i := range items {
		for j := range items[i].Variants {
			items[i].Variants[j].ImageThumbURL = thumbs[items[i].Variants[j].ID]
		}
	}
}
