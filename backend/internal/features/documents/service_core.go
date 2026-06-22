package documents

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/features/bundles"
	"smarterp/backend/internal/features/ledger"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
)

type Service struct {
	store   *db.Store
	repo    *Repository
	ledger  *ledger.Service
	bundles *bundles.Service
}

func NewService(
	store *db.Store,
	repo *Repository,
	ledger *ledger.Service,
	bundles *bundles.Service,
) *Service {
	return &Service{store: store, repo: repo, ledger: ledger, bundles: bundles}
}

func (s *Service) List(ctx context.Context, tenantID string, query httpx.ListQuery) ([]ListItem, int, error) {
	return s.repo.List(ctx, tenantID, query)
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (CreateResult, error) {
	req = normalizeRequest(req)
	if err := s.validateRequest(req); err != nil {
		return CreateResult{}, err
	}
	id := uuid.NewString()
	result := CreateResult{ID: id}
	err := s.store.WithTx(ctx, func(tx pgx.Tx) error {
		numbered, err := s.createDraftTx(ctx, tx, tenantID, id, req)
		result.Number = numbered.Number
		if err != nil {
			return err
		}
		result.Status, err = s.maybePostShiftDoc(ctx, tx, tenantID, id, req)
		return err
	})
	if err != nil {
		return CreateResult{}, mapDocumentWriteError(err)
	}
	return result, nil
}

func (s *Service) maybePostShiftDoc(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	id string,
	req CreateRequest,
) (string, error) {
	if req.ShiftID == "" {
		return "draft", nil
	}
	if err := s.postDocumentTx(ctx, tx, tenantID, id); err != nil {
		return "", err
	}
	return "posted", nil
}

func (s *Service) createDraftTx(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	documentID string,
	req CreateRequest,
) (CreateRequest, error) {
	if err := s.validateReferences(ctx, tx, tenantID, req); err != nil {
		return req, err
	}
	req, err := s.withDocumentNumber(ctx, tx, tenantID, req)
	if err != nil {
		return req, err
	}
	if err := s.repo.InsertDocument(ctx, tx, tenantID, documentID, req); err != nil {
		return req, err
	}
	if err := s.repo.ReplaceItems(ctx, tx, documentID, req.Items); err != nil {
		return req, err
	}
	return req, s.repo.ReplacePayments(ctx, tx, documentID, req.Payments)
}

func (s *Service) ByID(ctx context.Context, tenantID, id string) (Document, error) {
	doc, err := s.repo.ByID(ctx, tenantID, id)
	if err != nil {
		return Document{}, err
	}
	items, total, err := s.repo.LoadItemsWithProfit(ctx, tenantID, id)
	if err != nil {
		return Document{}, err
	}
	payments, err := s.repo.LoadPayments(ctx, id)
	if err != nil {
		return Document{}, err
	}
	doc.Items = items
	doc.Payments = payments
	doc.TotalProfit = total
	return doc, nil
}

func (s *Service) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
	req = normalizeRequest(req)
	if err := s.validateRequest(req); err != nil {
		return err
	}
	current, err := s.repo.ByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if err := ensureMutable(current, req); err != nil {
		return err
	}
	if current.Status == "draft" {
		return mapDocumentWriteError(s.updateDraft(ctx, tenantID, id, req))
	}
	if current.Status == "posted" {
		return mapDocumentWriteError(s.retroUpdate(ctx, tenantID, id, req))
	}
	return ErrStatusConflict
}

func ensureMutable(current Document, req UpdateRequest) error {
	if current.ShiftID != "" || req.ShiftID != "" {
		return ErrShiftDocumentLocked
	}
	if req.Type != current.Type {
		return ErrTypeImmutable
	}
	return nil
}

func (s *Service) updateDraft(ctx context.Context, tenantID, id string, req UpdateRequest) error {
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		if err := s.lockDocumentStatus(ctx, tx, tenantID, id, "draft"); err != nil {
			return err
		}
		return s.updateDraftTx(ctx, tx, tenantID, id, req)
	})
}

func (s *Service) lockDocumentStatus(ctx context.Context, tx pgx.Tx, tenantID, id, want string) error {
	status, err := s.repo.StatusForUpdate(ctx, tx, tenantID, id)
	if err != nil {
		return err
	}
	if status != want {
		return ErrStatusConflict
	}
	return nil
}

func (s *Service) updateDraftTx(ctx context.Context, tx pgx.Tx, tenantID, id string, req UpdateRequest) error {
	if err := s.validateReferences(ctx, tx, tenantID, req); err != nil {
		return err
	}
	req, err := s.withUpdateNumber(ctx, tx, tenantID, id, req)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateDocument(ctx, tx, tenantID, id, req); err != nil {
		return err
	}
	if err := s.repo.ReplaceItems(ctx, tx, id, req.Items); err != nil {
		return err
	}
	return s.repo.ReplacePayments(ctx, tx, id, req.Payments)
}

func (s *Service) Post(ctx context.Context, tenantID, id string) error {
	doc, err := s.repo.ByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if doc.Status != "draft" {
		return ErrDraftOnly
	}
	if doc.ShiftID != "" {
		return ErrShiftDocumentLocked
	}
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		if err := s.lockDocumentStatus(ctx, tx, tenantID, id, "draft"); err != nil {
			return err
		}
		return s.postDocumentTx(ctx, tx, tenantID, id)
	})
}

func (s *Service) Cancel(ctx context.Context, tenantID, id string) error {
	doc, err := s.repo.ByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if doc.Status != "posted" {
		return ErrPostedOnly
	}
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		if err := s.lockDocumentStatus(ctx, tx, tenantID, id, "posted"); err != nil {
			return err
		}
		if err := s.ensureShiftOpenForCancel(ctx, tx, tenantID, doc.ShiftID); err != nil {
			return err
		}
		return s.cancelTx(ctx, tx, tenantID, id)
	})
}

func (s *Service) ensureShiftOpenForCancel(ctx context.Context, tx pgx.Tx, tenantID, shiftID string) error {
	if shiftID == "" {
		return nil
	}
	open, err := s.repo.OpenShiftExists(ctx, tx, tenantID, shiftID)
	if err != nil {
		return err
	}
	if !open {
		return ErrShiftClosed
	}
	return nil
}

func (s *Service) cancelTx(ctx context.Context, tx pgx.Tx, tenantID, id string) error {
	if err := s.lockDocumentStatus(ctx, tx, tenantID, id, "posted"); err != nil {
		return err
	}
	affected, err := s.ledger.ReverseDocument(ctx, tx, tenantID, id)
	if err != nil {
		return err
	}
	if err := s.recalculateAffected(ctx, tx, tenantID, affected); err != nil {
		return err
	}
	return s.repo.SetStatus(ctx, tx, tenantID, id, "cancelled")
}

func (s *Service) Delete(ctx context.Context, tenantID, id string) error {
	status, err := s.repo.Status(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if status != "draft" {
		return ErrDraftOnly
	}
	return s.repo.DeleteDraft(ctx, tenantID, id)
}

func (s *Service) recalculateAffected(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	affected []ledger.VariantSequence,
) error {
	return s.ledger.RebuildAffected(ctx, tx, tenantID, affected)
}
