package reports

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"smarterp/backend/internal/features/ledger"
	"smarterp/backend/internal/shared/httpx"
)

type Service struct {
	repo   *Repository
	ledger *ledger.Service
}

func NewService(repo *Repository, ledger *ledger.Service) *Service {
	return &Service{repo: repo, ledger: ledger}
}

func (s *Service) Profit(ctx context.Context, query ProfitQuery) (decimal.Decimal, error) {
	filter := ledger.ProfitFilter{
		TenantID:    query.TenantID,
		FromDate:   query.FromDate,
		ToDate:     query.ToDate,
		WarehouseID: query.WarehouseID,
		VariantID:   query.VariantID,
	}
	return s.ledger.ProfitByPeriod(ctx, filter)
}

func (s *Service) Stock(ctx context.Context, tenantID, warehouseID string) ([]StockRow, error) {
	return s.repo.StockRows(ctx, tenantID, warehouseID)
}

func (s *Service) FullStock(ctx context.Context, tenantID string, query httpx.ListQuery) ([]FullStockRow, int, error) {
	return s.repo.FullStockRows(ctx, tenantID, query)
}

func (s *Service) TopProducts(ctx context.Context, tenantID string, fromDate, toDate time.Time) ([]TopProduct, error) {
	return s.repo.TopProducts(ctx, tenantID, fromDate, toDate)
}

func (s *Service) Movements(ctx context.Context, tenantID, variantID string) ([]ledger.Movement, error) {
	return s.ledger.Movements(ctx, tenantID, variantID)
}
