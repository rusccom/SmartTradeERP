package reports

import "github.com/shopspring/decimal"

type ProfitReport struct {
	Profit decimal.Decimal `json:"profit"`
}

type StockRow struct {
	VariantID string          `json:"variant_id"`
	Name      string          `json:"name"`
	Qty       decimal.Decimal `json:"qty"`
	Avg       decimal.Decimal `json:"avg"`
}

type TopProduct struct {
	ProductID string          `json:"product_id"`
	Name      string          `json:"name"`
	Profit    decimal.Decimal `json:"profit"`
}
