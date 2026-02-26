package shifts

import "github.com/shopspring/decimal"

type OpenRequest struct {
	WarehouseID string          `json:"warehouse_id"`
	OpeningCash decimal.Decimal `json:"opening_cash"`
}

type CashOpRequest struct {
	Type   string          `json:"type"`
	Amount decimal.Decimal `json:"amount"`
	Note   string          `json:"note"`
}

type Shift struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	WarehouseID string          `json:"warehouse_id"`
	OpenedAt    string          `json:"opened_at"`
	ClosedAt    string          `json:"closed_at"`
	OpeningCash decimal.Decimal `json:"opening_cash"`
	ClosingCash decimal.Decimal `json:"closing_cash"`
	Status      string          `json:"status"`
}

type CashOp struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Amount    decimal.Decimal `json:"amount"`
	Note      string          `json:"note"`
	CreatedAt string          `json:"created_at"`
}

type ShiftReport struct {
	Shift        Shift           `json:"shift"`
	CashOps      []CashOp        `json:"cash_ops"`
	TotalSales   decimal.Decimal `json:"total_sales"`
	TotalReturns decimal.Decimal `json:"total_returns"`
	SalesCash    decimal.Decimal `json:"sales_cash"`
	SalesCard    decimal.Decimal `json:"sales_card"`
	ReturnsCash  decimal.Decimal `json:"returns_cash"`
	ReturnsCard  decimal.Decimal `json:"returns_card"`
	TotalCashIn  decimal.Decimal `json:"total_cash_in"`
	TotalCashOut decimal.Decimal `json:"total_cash_out"`
	ExpectedCash decimal.Decimal `json:"expected_cash"`
}
