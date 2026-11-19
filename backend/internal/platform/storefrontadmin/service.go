package storefrontadmin

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/platform/storefront/themes"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/validation"
)

const previewTTL = 15 * time.Minute

type Service struct {
	store     *db.Store
	repo      *Repository
	registry  *themes.Index
	tokens    *auth.TokenService
	shopsBase string
}

func NewService(store *db.Store, repo *Repository, registry *themes.Index, tokens *auth.TokenService, shopsBase string) *Service {
	return &Service{store: store, repo: repo, registry: registry, tokens: tokens, shopsBase: shopsBase}
}

// Themes lists installed themes for the picker.
func (s *Service) Themes() []themes.Info {
	return s.registry.Items
}

func (s *Service) Get(ctx context.Context, tenantID string) (Settings, error) {
	row, err := s.repo.Load(ctx, tenantID)
	if err != nil {
		return Settings{}, err
	}
	host, err := s.ensureHost(ctx, tenantID)
	if err != nil {
		return Settings{}, err
	}
	if !s.registry.Has(row.ThemeID) {
		row.ThemeID = s.registry.Default
	}
	return Settings{ThemeID: row.ThemeID, DraftThemeID: row.DraftThemeID, ShopURL: shopURL(host)}, nil
}

// SaveDraft stores only the selected theme.
func (s *Service) SaveDraft(ctx context.Context, tenantID string, req DraftRequest) error {
	if !s.registry.Has(req.ThemeID) {
		return validation.ErrInvalidData
	}
	return s.store.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		return s.repo.SaveDraft(ctx, tx, tenantID, req.ThemeID)
	})
}

// Publish promotes the draft theme to the live storefront.
func (s *Service) Publish(ctx context.Context, tenantID string) error {
	return s.store.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		return s.repo.Publish(ctx, tx, tenantID)
	})
}

// Preview mints a short-lived token and the storefront URL that renders the
// draft. URL is empty when the tenant has no active storefront host yet.
func (s *Service) Preview(ctx context.Context, tenantID string) (Preview, error) {
	token, err := s.tokens.IssuePreview(tenantID, previewTTL)
	if err != nil {
		return Preview{}, err
	}
	host, err := s.repo.PrimaryHost(ctx, tenantID)
	if err != nil {
		return Preview{}, err
	}
	return Preview{Token: token, URL: previewURL(host, token)}, nil
}

func previewURL(host, token string) string {
	if host == "" {
		return ""
	}
	return "https://" + host + "/?preview=" + token
}
