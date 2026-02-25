package reports

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"smarterp/backend/internal/features/ledger"
)

type Service struct {
	repo   *Repository
	ledger *ledger.Service
}

func NewService(repo *Repository, ledger *ledger.Service) *Service {
	return &Service{repo: repo, ledger: ledger}
}

func (s *Service) Profit(
	ctx context.Context,
	tenantID string,
	fromDate time.Time,
	toDate time.Time,
	warehouseID string,
	variantID string,
) (decimal.Decimal, error) {
	return s.ledger.ProfitByPeriod(ctx, tenantID, fromDate, toDate, warehouseID, variantID)
}

func (s *Service) Stock(ctx context.Context, tenantID, warehouseID string) ([]StockRow, error) {
	return s.repo.StockRows(ctx, tenantID, warehouseID)
}

func (s *Service) TopProducts(ctx context.Context, tenantID string, fromDate, toDate time.Time) ([]TopProduct, error) {
	return s.repo.TopProducts(ctx, tenantID, fromDate, toDate)
}

func (s *Service) Movements(ctx context.Context, tenantID, variantID string) ([]ledger.Movement, error) {
	return s.ledger.Movements(ctx, tenantID, variantID)
}
