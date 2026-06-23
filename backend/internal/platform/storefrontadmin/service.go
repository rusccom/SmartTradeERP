package storefrontadmin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/platform/storefront"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/validation"
)

const previewTTL = 15 * time.Minute

type Service struct {
	store    *db.Store
	repo     *Repository
	registry *storefront.Registry
	tokens   *auth.TokenService
}

func NewService(store *db.Store, repo *Repository, registry *storefront.Registry, tokens *auth.TokenService) *Service {
	return &Service{store: store, repo: repo, registry: registry, tokens: tokens}
}

// Themes lists installed themes for the picker.
func (s *Service) Themes() []storefront.ThemeInfo {
	return s.registry.Themes()
}

func (s *Service) Get(ctx context.Context, tenantID string) (Settings, error) {
	row, err := s.repo.Load(ctx, tenantID)
	if err != nil {
		return Settings{}, err
	}
	return Settings{
		ThemeID:         row.ThemeID,
		DraftThemeID:    row.DraftThemeID,
		PublishedTokens: decodeTokens(row.PublishedTokens),
		DraftTokens:     decodeTokens(row.DraftTokens),
		Sections:        decodeSections(row.Sections),
		DraftSections:   normalizeSections(decodeSections(row.DraftSections)),
	}, nil
}

// SaveDraft validates the theme, drops unknown/unsafe token overrides, and
// normalizes sections before storing the draft. The live storefront is untouched.
func (s *Service) SaveDraft(ctx context.Context, tenantID string, req DraftRequest) error {
	if !s.registry.Has(req.ThemeID) {
		return validation.ErrInvalidData
	}
	tokens, err := json.Marshal(s.registry.SanitizeOverrides(req.ThemeID, req.Tokens))
	if err != nil {
		return err
	}
	sections, err := json.Marshal(normalizeSections(req.Sections))
	if err != nil {
		return err
	}
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.SaveDraft(ctx, tx, tenantID, req.ThemeID, tokens, sections)
	})
}

// Publish promotes the draft theme, tokens and sections to the live storefront.
func (s *Service) Publish(ctx context.Context, tenantID string) error {
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
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

func decodeTokens(raw []byte) map[string]string {
	out := map[string]string{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &out)
	}
	return out
}

func decodeSections(raw []byte) []SectionInput {
	out := []SectionInput{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &out)
	}
	return out
}

// normalizeSections keeps known keys in the given order (deduped), then appends
// any missing known keys as enabled, so the stored config is always complete.
func normalizeSections(input []SectionInput) []SectionInput {
	known := storefront.HomeSectionKeys()
	out := make([]SectionInput, 0, len(known))
	seen := map[string]bool{}
	for _, item := range input {
		if containsKey(known, item.Key) && !seen[item.Key] {
			out = append(out, SectionInput{Key: item.Key, Enabled: item.Enabled})
			seen[item.Key] = true
		}
	}
	return appendMissing(out, known, seen)
}

func appendMissing(out []SectionInput, known []string, seen map[string]bool) []SectionInput {
	for _, key := range known {
		if !seen[key] {
			out = append(out, SectionInput{Key: key, Enabled: true})
		}
	}
	return out
}

func containsKey(keys []string, key string) bool {
	for _, item := range keys {
		if item == key {
			return true
		}
	}
	return false
}
