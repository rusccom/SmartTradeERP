package documents

import "github.com/shopspring/decimal"

type ItemInput struct {
	VariantID string          `json:"variant_id"`
	Qty       decimal.Decimal `json:"qty"`
	UnitPrice decimal.Decimal `json:"unit_price"`
}

type CreateRequest struct {
	Type              string      `json:"type"`
	Date              string      `json:"date"`
	Number            string      `json:"number"`
	WarehouseID       string      `json:"warehouse_id"`
	SourceWarehouseID string      `json:"source_warehouse_id"`
	TargetWarehouseID string      `json:"target_warehouse_id"`
	Note              string      `json:"note"`
	Items             []ItemInput `json:"items"`
}

type UpdateRequest = CreateRequest

type Document struct {
	ID                string          `json:"id"`
	Type              string          `json:"type"`
	Date              string          `json:"date"`
	Number            string          `json:"number"`
	Status            string          `json:"status"`
	WarehouseID       string          `json:"warehouse_id"`
	SourceWarehouseID string          `json:"source_warehouse_id"`
	TargetWarehouseID string          `json:"target_warehouse_id"`
	Note              string          `json:"note"`
	Items             []DocumentItem  `json:"items"`
	TotalProfit       decimal.Decimal `json:"total_profit"`
}

type DocumentItem struct {
	ID          string          `json:"id"`
	VariantID   string          `json:"variant_id"`
	Qty         decimal.Decimal `json:"qty"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Profit      decimal.Decimal `json:"profit"`
}

type ListItem struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Date   string `json:"date"`
	Number string `json:"number"`
	Status string `json:"status"`
	Note   string `json:"note"`
}

type Filters struct {
	Type   string
	Status string
	Date   string
}

type postingItem struct {
	ID          string
	VariantID   string
	Qty         decimal.Decimal
	UnitPrice   decimal.Decimal
	TotalAmount decimal.Decimal
	IsComposite bool
}

type variantComponent struct {
	ComponentVariantID string
	QtyPerUnit         decimal.Decimal
}
