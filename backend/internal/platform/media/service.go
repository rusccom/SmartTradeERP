package media

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/storage"
	"smarterp/backend/internal/shared/validation"
)

const directUploadTTL = 15 * time.Minute

type Service struct {
	store   *db.Store
	repo    *Repository
	objects storage.ObjectStore
}

func NewService(store *db.Store, repo *Repository, objects storage.ObjectStore) *Service {
	return &Service{store: store, repo: repo, objects: objects}
}

func (s *Service) List(
	ctx context.Context,
	tenantID string,
	ownerType string,
	ownerID string,
) ([]Item, error) {
	if s.objects == nil {
		return nil, ErrStorageNotConfigured
	}
	if !validScope(tenantID, ownerType, ownerID) {
		return nil, ErrInvalidMedia
	}
	items, err := s.repo.List(ctx, tenantID, ownerType, ownerID)
	return s.attachURLs(items), err
}

func (s *Service) CreateDirectUpload(
	ctx context.Context,
	req DirectUploadRequest,
) (DirectUpload, error) {
	if s.objects == nil {
		return DirectUpload{}, ErrStorageNotConfigured
	}
	if err := validateDirectUpload(req); err != nil {
		return DirectUpload{}, err
	}
	return s.createPendingUpload(ctx, req)
}

func (s *Service) CompleteDirectUpload(
	ctx context.Context,
	req CompleteRequest,
) (Item, error) {
	if s.objects == nil {
		return Item{}, ErrStorageNotConfigured
	}
	if !validCompleteRequest(req) {
		return Item{}, ErrInvalidMedia
	}
	item, err := s.repo.GetPending(ctx, req.TenantID, req.OwnerType, req.OwnerID, req.MediaID)
	if err != nil {
		return Item{}, err
	}
	return s.confirmUploaded(ctx, req, item)
}

func validateDirectUpload(req DirectUploadRequest) error {
	if !validScope(req.TenantID, req.OwnerType, req.OwnerID) {
		return ErrInvalidMedia
	}
	if req.Input.SizeBytes <= 0 || req.Input.SizeBytes > MaxUploadBytes {
		return ErrInvalidMedia
	}
	if !validContentType(req.Input.ContentType) {
		return ErrInvalidMedia
	}
	return nil
}

func validCompleteRequest(req CompleteRequest) bool {
	return validScope(req.TenantID, req.OwnerType, req.OwnerID) && validation.UUID(req.MediaID)
}

func validScope(tenantID string, ownerType string, ownerID string) bool {
	return validation.UUID(tenantID) && validOwnerType(ownerType) && validation.UUID(ownerID)
}

func (s *Service) createPendingUpload(
	ctx context.Context,
	req DirectUploadRequest,
) (DirectUpload, error) {
	item := newPendingItem(req)
	signed, err := s.presign(item)
	if err != nil {
		return DirectUpload{}, err
	}
	if err := s.repo.CreatePending(ctx, req.TenantID, item, signed.ExpiresAt); err != nil {
		return DirectUpload{}, err
	}
	return newDirectUpload(item.ID, signed), nil
}

func (s *Service) presign(item Item) (storage.PresignedPut, error) {
	return s.objects.PresignPut(storage.PresignPutRequest{
		Key: item.ObjectKey, ContentType: item.ContentType, Expires: directUploadTTL,
	})
}

func (s *Service) confirmUploaded(ctx context.Context, req CompleteRequest, item Item) (Item, error) {
	if err := s.ensureUploaded(ctx, item); err != nil {
		_ = s.objects.Delete(ctx, item.ObjectKey)
		_ = s.repo.Delete(ctx, req.TenantID, item.ID)
		return Item{}, err
	}
	if err := s.markReady(ctx, req, item); err != nil {
		return Item{}, err
	}
	item.Status = "ready"
	item.IsPrimary = true
	item.URL = s.objects.PublicURL(item.ObjectKey)
	return item, nil
}

func (s *Service) ensureUploaded(ctx context.Context, item Item) error {
	info, err := s.objects.Head(ctx, item.ObjectKey)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
			return ErrInvalidMedia
		}
		return err
	}
	if info.SizeBytes != item.SizeBytes || info.ContentType != item.ContentType {
		return ErrInvalidMedia
	}
	return nil
}

func (s *Service) markReady(ctx context.Context, req CompleteRequest, item Item) error {
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		if err := s.repo.ClearPrimary(ctx, tx, req.TenantID, req.OwnerType, req.OwnerID); err != nil {
			return err
		}
		return s.repo.MarkReady(ctx, tx, req.TenantID, item)
	})
}

func (s *Service) attachURLs(items []Item) []Item {
	for index := range items {
		items[index].URL = s.objects.PublicURL(items[index].ObjectKey)
	}
	return items
}

func newPendingItem(req DirectUploadRequest) Item {
	id := uuid.NewString()
	return Item{
		ID: id, OwnerType: req.OwnerType, OwnerID: req.OwnerID,
		ObjectKey: mediaObjectKey(req, id), FileName: cleanFileName(req.Input.FileName),
		ContentType: req.Input.ContentType, SizeBytes: req.Input.SizeBytes,
		Status: "pending",
	}
}

func newDirectUpload(id string, signed storage.PresignedPut) DirectUpload {
	return DirectUpload{
		ID: id, UploadURL: signed.URL, Method: signed.Method,
		Headers: signed.Headers, ExpiresAt: formatTime(signed.ExpiresAt),
	}
}

func mediaObjectKey(req DirectUploadRequest, mediaID string) string {
	return "tenants/" + req.TenantID + "/" + req.OwnerType + "/" +
		req.OwnerID + "/" + mediaID + fileExtension(req.Input.ContentType)
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}
