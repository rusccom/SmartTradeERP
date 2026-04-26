package shifts

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"

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

func (s *Service) Open(ctx context.Context, tenantID, userID string, req OpenRequest) (string, error) {
	req.WarehouseID = validation.Clean(req.WarehouseID)
	if err := validateOpenRequest(req); err != nil {
		return "", err
	}
	if err := s.ensureWarehouse(ctx, tenantID, req.WarehouseID); err != nil {
		return "", err
	}
	if err := s.ensureNoOpenShift(ctx, tenantID, userID); err != nil {
		return "", err
	}
	shift := makeOpenShift(userID, req)
	err := s.store.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.Insert(ctx, tx, tenantID, shift)
	})
	if err != nil {
		return "", mapOpenInsertError(err)
	}
	return shift.ID, nil
}

func makeOpenShift(userID string, req OpenRequest) Shift {
	return Shift{
		ID:          uuid.NewString(),
		UserID:      userID,
		WarehouseID: req.WarehouseID,
		OpeningCash: req.OpeningCash,
	}
}

func validateOpenRequest(req OpenRequest) error {
	if !validation.UUID(req.WarehouseID) {
		return validation.ErrInvalidData
	}
	if !req.OpeningCash.GreaterThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	return nil
}

func (s *Service) ensureWarehouse(ctx context.Context, tenantID string, id string) error {
	exists, err := s.repo.WarehouseExists(ctx, tenantID, id)
	if err != nil || !exists {
		return shiftReferenceError(err)
	}
	return nil
}

func shiftReferenceError(err error) error {
	if err != nil {
		return err
	}
	return ErrInvalidShiftReference
}

func (s *Service) ensureNoOpenShift(ctx context.Context, tenantID, userID string) error {
	_, err := s.repo.FindOpen(ctx, tenantID, userID)
	if err == nil {
		return ErrShiftAlreadyOpen
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	return err
}

func mapOpenInsertError(err error) error {
	if !isUniqueViolation(err) {
		return err
	}
	return ErrShiftAlreadyOpen
}

func isUniqueViolation(err error) bool {
	pgErr := &pgconn.PgError{}
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (s *Service) Current(ctx context.Context, tenantID, userID string) (Shift, error) {
	shift, err := s.repo.FindOpen(ctx, tenantID, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Shift{}, ErrNoOpenShift
	}
	if err != nil {
		return Shift{}, err
	}
	return shift, nil
}

func (s *Service) CashOp(ctx context.Context, tenantID, userID string, req CashOpRequest) error {
	req = normalizeCashOp(req)
	shift, err := s.Current(ctx, tenantID, userID)
	if err != nil {
		return err
	}
	if err := validateCashOp(req); err != nil {
		return err
	}
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.InsertCashOp(ctx, tx, shift.ID, req)
	})
}

func validateCashOp(req CashOpRequest) error {
	if req.Type != "cash_in" && req.Type != "cash_out" {
		return ErrInvalidCashOpType
	}
	if !req.Amount.GreaterThan(decimal.Zero) {
		return ErrInvalidAmount
	}
	if !validation.Max(req.Note, 1000) {
		return validation.ErrInvalidData
	}
	return nil
}

func normalizeCashOp(req CashOpRequest) CashOpRequest {
	req.Type = validation.Clean(req.Type)
	req.Note = validation.Clean(req.Note)
	return req
}

func (s *Service) Close(ctx context.Context, tenantID, userID string) (ShiftReport, error) {
	shift, err := s.Current(ctx, tenantID, userID)
	if err != nil {
		return ShiftReport{}, err
	}
	if err := s.closeWithExpectedCash(ctx, tenantID, shift); err != nil {
		return ShiftReport{}, err
	}
	return s.Report(ctx, tenantID, shift.ID)
}

func (s *Service) closeWithExpectedCash(ctx context.Context, tenantID string, shift Shift) error {
	return s.store.WithTx(ctx, func(tx pgx.Tx) error {
		totals, err := s.repo.ShiftTotals(ctx, tenantID, shift.ID)
		if err != nil {
			return err
		}
		expected := expectedCash(shift.OpeningCash, totals)
		return s.repo.Close(ctx, tx, tenantID, shift.ID, expected)
	})
}

func (s *Service) Report(ctx context.Context, tenantID, shiftID string) (ShiftReport, error) {
	shift, err := s.repo.ByID(ctx, tenantID, shiftID)
	if err != nil {
		return ShiftReport{}, err
	}
	cashOps, err := s.repo.CashOps(ctx, shiftID)
	if err != nil {
		return ShiftReport{}, err
	}
	totals, err := s.repo.ShiftTotals(ctx, tenantID, shiftID)
	if err != nil {
		return ShiftReport{}, err
	}
	return buildReport(shift, cashOps, totals), nil
}

func buildReport(shift Shift, cashOps []CashOp, totals shiftTotals) ShiftReport {
	return ShiftReport{
		Shift:        shift,
		CashOps:      cashOps,
		TotalSales:   totals.totalSales,
		TotalReturns: totals.totalReturns,
		SalesCash:    totals.salesCash,
		SalesCard:    totals.salesCard,
		ReturnsCash:  totals.returnsCash,
		ReturnsCard:  totals.returnsCard,
		TotalCashIn:  totals.totalCashIn,
		TotalCashOut: totals.totalCashOut,
		ExpectedCash: expectedCash(shift.OpeningCash, totals),
	}
}

func expectedCash(openingCash decimal.Decimal, totals shiftTotals) decimal.Decimal {
	value := openingCash
	value = value.Add(totals.salesCash)
	value = value.Sub(totals.returnsCash)
	value = value.Add(totals.totalCashIn)
	value = value.Sub(totals.totalCashOut)
	return value.Round(4)
}
