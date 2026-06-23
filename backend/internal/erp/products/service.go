package products

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"smarterp/backend/internal/erp/ledger"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/validation"
)

var ErrHasMovements = errors.New("product has movements")
var ErrCompositeTypeLocked = errors.New("product composite type is locked")
var ErrUsedInBundle = errors.New("product used in bundle")
var ErrSlugTaken = errors.New("product slug already in use")

type Service struct {
	store       *db.Store
	repo        *Repository
	ledger      *ledger.Service
	bundleState ComponentStateReader
	media       MediaService
}

type ComponentStateReader interface {
	ProductHasComponents(ctx context.Context, tenantID string, productID string) (bool, error)
	ProductUsedAsComponent(ctx context.Context, tenantID string, productID string) (bool, error)
}

func NewService(
	store *db.Store,
	repo *Repository,
	ledger *ledger.Service,
	bundleState ComponentStateReader,
) *Service {
	return &Service{store: store, repo: repo, ledger: ledger, bundleState: bundleState}
}

func (s *Service) List(ctx context.Context, tenantID string, query ProductListQuery) ([]Product, int, error) {
	return s.repo.List(ctx, tenantID, query)
}

func (s *Service) ListWithIncludes(
	ctx context.Context,
	tenantID string,
	query ProductListQuery,
	include ProductListInclude,
) ([]ProductListItem, int, error) {
	items, total, err := s.repo.ListWithIncludes(ctx, tenantID, query, include)
	if err != nil {
		return nil, 0, err
	}
	s.attachImages(ctx, tenantID, items, include)
	return items, total, nil
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (string, error) {
	req = normalizeCreate(req)
	if err := validateCreate(req); err != nil {
		return "", err
	}
	slug, err := s.ensureUniqueSlug(ctx, tenantID, req.Slug, "")
	if err != nil {
		return "", err
	}
	req.Slug = slug
	productID := uuid.NewString()
	variantID := uuid.NewString()
	err = s.store.WithTx(ctx, func(tx pgx.Tx) error {
		input := createProductTx{tenantID: tenantID, productID: productID, variantID: variantID, req: req}
		return s.createWithDefaultVariant(ctx, tx, input)
	})
	if err != nil {
		return "", mapProductWriteError(err)
	}
	return productID, nil
}

type createProductTx struct {
	tenantID  string
	productID string
	variantID string
	req       CreateRequest
}

func (s *Service) createWithDefaultVariant(
	ctx context.Context,
	tx pgx.Tx,
	input createProductTx,
) error {
	if err := s.repo.Create(ctx, tx, input.tenantID, input.productID, input.req); err != nil {
		return err
	}
	return s.repo.CreateDefaultVariant(ctx, tx, input)
}

func (s *Service) ByID(ctx context.Context, tenantID, id string) (Product, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *Service) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
	req = normalizeUpdate(req)
	if err := validateUpdate(req); err != nil {
		return err
	}
	if err := s.ensureCompositeChangeAllowed(ctx, tenantID, id, req.IsComposite); err != nil {
		return err
	}
	slug, err := s.ensureUniqueSlug(ctx, tenantID, req.Slug, id)
	if err != nil {
		return err
	}
	req.Slug = slug
	return mapProductWriteError(s.repo.Update(ctx, tenantID, id, req))
}

func (s *Service) ensureCompositeChangeAllowed(ctx context.Context, tenantID, id string, next bool) error {
	current, err := s.repo.CompositeFlag(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if current == next {
		return nil
	}
	return s.ensureCompositeUnlocked(ctx, tenantID, id)
}

func (s *Service) ensureCompositeUnlocked(ctx context.Context, tenantID, id string) error {
	blocked, err := s.productHasBundleState(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if blocked {
		return ErrCompositeTypeLocked
	}
	return nil
}

func (s *Service) productHasBundleState(ctx context.Context, tenantID, id string) (bool, error) {
	hasMovements, err := s.ledger.HasProductMovements(ctx, tenantID, id)
	if err != nil || hasMovements {
		return hasMovements, err
	}
	return s.productLinkedToBundle(ctx, tenantID, id)
}

func (s *Service) Delete(ctx context.Context, tenantID, id string) error {
	hasMovements, err := s.ledger.HasProductMovements(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if hasMovements {
		return ErrHasMovements
	}
	if err := s.ensureProductNotInBundle(ctx, tenantID, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, tenantID, id)
}

func (s *Service) ensureProductNotInBundle(ctx context.Context, tenantID, id string) error {
	linked, err := s.productLinkedToBundle(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if linked {
		return ErrUsedInBundle
	}
	return nil
}

func (s *Service) productLinkedToBundle(ctx context.Context, tenantID, id string) (bool, error) {
	hasComponents, err := s.bundleState.ProductHasComponents(ctx, tenantID, id)
	if err != nil || hasComponents {
		return hasComponents, err
	}
	return s.bundleState.ProductUsedAsComponent(ctx, tenantID, id)
}

func normalizeCreate(req CreateRequest) CreateRequest {
	req.Name = validation.Clean(req.Name)
	req.Unit = validation.Clean(req.Unit)
	req.SKUCode = validation.Clean(req.SKUCode)
	req.Barcode = validation.Clean(req.Barcode)
	req.VariantName = validation.Clean(req.VariantName)
	req.Slug = resolveSlug(req.Slug, req.Name)
	req.SEOTitle = validation.Clean(req.SEOTitle)
	req.SEODescription = validation.Clean(req.SEODescription)
	return req
}

func normalizeUpdate(req UpdateRequest) UpdateRequest {
	req.Name = validation.Clean(req.Name)
	req.Slug = resolveSlug(req.Slug, req.Name)
	req.SEOTitle = validation.Clean(req.SEOTitle)
	req.SEODescription = validation.Clean(req.SEODescription)
	return req
}

// resolveSlug normalizes the handle, falling back to the product name so every
// product gets a products/<handle> URL even when the user leaves it blank.
func resolveSlug(slug, name string) string {
	normalized := normalizeSlug(slug)
	if normalized != "" {
		return normalized
	}
	return normalizeSlug(name)
}

func validateCreate(req CreateRequest) error {
	if validateName(req.Name) != nil || !validation.Required(req.Unit) {
		return validation.ErrInvalidData
	}
	if !validation.NonNegative(req.Price) || !validation.Max(req.Unit, 24) {
		return validation.ErrInvalidData
	}
	if !validateSEO(req.Slug, req.SEOTitle, req.SEODescription) {
		return validation.ErrInvalidData
	}
	return nil
}

func validateUpdate(req UpdateRequest) error {
	if err := validateName(req.Name); err != nil {
		return err
	}
	if !validateSEO(req.Slug, req.SEOTitle, req.SEODescription) {
		return validation.ErrInvalidData
	}
	return nil
}

func validateSEO(slug, title, description string) bool {
	return validation.Max(slug, 200) && validation.Max(title, 255) && validation.Max(description, 320)
}

var cyrillicMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "i", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
}

// normalizeSlug lowercases the handle, transliterates Cyrillic to Latin, and
// keeps only [a-z0-9], collapsing every other run into a single hyphen so even
// a Cyrillic product name yields a URL-safe handle.
func normalizeSlug(value string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if latin, ok := cyrillicMap[r]; ok {
			b.WriteString(latin)
			dash = false
			continue
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			dash = false
			continue
		}
		if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func mapProductWriteError(err error) error {
	pgErr := &pgconn.PgError{}
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrSlugTaken
	}
	return err
}

// ensureUniqueSlug keeps the handle unique per tenant by appending -2, -3, ...
// when the base is already taken, so two products with the same name still get
// distinct products/<handle> URLs instead of failing.
func (s *Service) ensureUniqueSlug(ctx context.Context, tenantID, base, excludeID string) (string, error) {
	if base == "" {
		return "", nil
	}
	candidate := base
	for i := 2; i < 1000; i++ {
		taken, err := s.repo.SlugExists(ctx, tenantID, candidate, excludeID)
		if err != nil {
			return "", err
		}
		if !taken {
			return candidate, nil
		}
		candidate = base + "-" + strconv.Itoa(i)
	}
	return candidate, nil
}

func validateName(name string) error {
	if !validation.Required(name) || !validation.Max(name, 200) {
		return validation.ErrInvalidData
	}
	return nil
}
