package bundles

import (
	"context"

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

func (s *Service) List(ctx context.Context, tenantID string, page, perPage int) ([]Bundle, int, error) {
	return s.repo.List(ctx, tenantID, page, perPage)
}

func (s *Service) ByID(ctx context.Context, tenantID, variantID string) (Bundle, error) {
	item, err := s.repo.ByVariantID(ctx, tenantID, variantID)
	if err != nil {
		return Bundle{}, err
	}
	components, err := s.repo.Components(ctx, tenantID, variantID)
	item.Components = components
	return item, err
}

func (s *Service) Components(ctx context.Context, tenantID, variantID string) ([]Component, error) {
	if _, err := s.repo.ByVariantID(ctx, tenantID, variantID); err != nil {
		return nil, err
	}
	return s.repo.Components(ctx, tenantID, variantID)
}

func (s *Service) SetComponents(ctx context.Context, tenantID, variantID string, items []Component) error {
	items = normalizeComponents(items)
	if err := validateComponents(variantID, items); err != nil {
		return err
	}
	if err := s.validateState(ctx, tenantID, variantID, items); err != nil {
		return err
	}
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.ReplaceComponents(ctx, tx, variantID, items)
	})
}

func (s *Service) ResolveComponents(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	variantID string,
) ([]Component, error) {
	items, err := s.repo.ComponentsTx(ctx, tx, tenantID, variantID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrMissingComponents
	}
	return items, nil
}

func (s *Service) SaveSnapshot(ctx context.Context, tx pgx.Tx, input SnapshotInput) error {
	return s.repo.SaveSnapshot(ctx, tx, input)
}

func (s *Service) ProductHasComponents(ctx context.Context, tenantID, productID string) (bool, error) {
	return s.repo.ProductHasComponents(ctx, tenantID, productID)
}

func (s *Service) ProductUsedAsComponent(ctx context.Context, tenantID, productID string) (bool, error) {
	return s.repo.ProductUsedAsComponent(ctx, tenantID, productID)
}

func (s *Service) VariantHasComponents(ctx context.Context, tenantID, variantID string) (bool, error) {
	return s.repo.VariantHasComponents(ctx, tenantID, variantID)
}

func (s *Service) VariantUsedAsComponent(ctx context.Context, tenantID, variantID string) (bool, error) {
	return s.repo.VariantUsedAsComponent(ctx, tenantID, variantID)
}

func (s *Service) validateState(ctx context.Context, tenantID, variantID string, items []Component) error {
	isComposite, err := s.repo.VariantComposite(ctx, tenantID, variantID)
	if err != nil {
		return err
	}
	if !validState(isComposite, items) {
		return ErrInvalidComponentState
	}
	return s.ensureComponentsUsable(ctx, tenantID, items)
}

func validState(isComposite bool, items []Component) bool {
	return isComposite && len(items) > 0
}

func (s *Service) ensureComponentsUsable(ctx context.Context, tenantID string, items []Component) error {
	for _, item := range items {
		valid, err := s.repo.ComponentUsable(ctx, tenantID, item.ComponentVariantID)
		if err != nil {
			return err
		}
		if !valid {
			return validation.ErrInvalidData
		}
	}
	return nil
}

func normalizeComponents(items []Component) []Component {
	for index := range items {
		items[index].ComponentVariantID = validation.Clean(items[index].ComponentVariantID)
	}
	return items
}

func validateComponents(variantID string, items []Component) error {
	if !validation.UUID(variantID) {
		return validation.ErrInvalidData
	}
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		if invalidComponent(variantID, item, seen) {
			return validation.ErrInvalidData
		}
		seen[item.ComponentVariantID] = true
	}
	return nil
}

func invalidComponent(variantID string, item Component, seen map[string]bool) bool {
	if !validation.Required(item.ComponentVariantID) || !validation.UUID(item.ComponentVariantID) {
		return true
	}
	if !validation.Positive(item.Qty) {
		return true
	}
	return item.ComponentVariantID == variantID || seen[item.ComponentVariantID]
}
