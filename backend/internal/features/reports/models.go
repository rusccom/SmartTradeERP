package reports

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProfitQuery struct {
	TenantID    string
	FromDate   time.Time
	ToDate     time.Time
	WarehouseID string
	VariantID   string
}

type ProfitReport struct {
	Profit decimal.Decimal `json:"profit"`
}

type StockRow struct {
	VariantID string          `json:"variant_id"`
	Name      string          `json:"name"`
	Qty       decimal.Decimal `json:"qty"`
	Avg       decimal.Decimal `json:"avg"`
}

type FullStockRow struct {
	ProductID   string               `json:"product_id"`
	ProductName string               `json:"product_name"`
	VariantID   string               `json:"variant_id"`
	VariantName string               `json:"variant_name"`
	Name        string               `json:"name"`
	SKUCode     string               `json:"sku_code"`
	Barcode     string               `json:"barcode"`
	Unit        string               `json:"unit"`
	GlobalQty   decimal.Decimal      `json:"global_qty"`
	Avg         decimal.Decimal      `json:"avg"`
	StockValue  decimal.Decimal      `json:"stock_value"`
	Warehouses  []FullStockWarehouse `json:"warehouses"`
}

type FullStockWarehouse struct {
	WarehouseID string          `json:"warehouse_id"`
	Warehouse   string          `json:"warehouse"`
	Qty         decimal.Decimal `json:"qty"`
}

type TopProduct struct {
	ProductID string          `json:"product_id"`
	Name      string          `json:"name"`
	Profit    decimal.Decimal `json:"profit"`
}
