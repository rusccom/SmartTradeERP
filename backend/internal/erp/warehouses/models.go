package warehouses

import "github.com/shopspring/decimal"

type Warehouse struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Address   string `json:"address"`
    IsDefault bool   `json:"is_default"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
}

type WarehouseListItem struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Address    string               `json:"address"`
	IsDefault  bool                 `json:"is_default"`
	IsActive   bool                 `json:"is_active"`
	StockValue decimal.Decimal      `json:"stock_value"`
	Stock      []WarehouseStockItem `json:"stock"`
	CreatedAt  string               `json:"created_at"`
}

type WarehouseStockItem struct {
	ProductID   string          `json:"product_id"`
	ProductName string          `json:"product_name"`
	VariantID   string          `json:"variant_id"`
	VariantName string          `json:"variant_name"`
	Name        string          `json:"name"`
	SKUCode     string          `json:"sku_code"`
	Barcode     string          `json:"barcode"`
	Unit        string          `json:"unit"`
	Qty         decimal.Decimal `json:"qty"`
	AvgCost     decimal.Decimal `json:"avg_cost"`
	StockValue  decimal.Decimal `json:"stock_value"`
}

type WarehouseListInclude struct {
	Stock bool
}

type CreateRequest struct {
    Name      string `json:"name"`
    Address   string `json:"address"`
    IsDefault bool   `json:"is_default"`
    IsActive  bool   `json:"is_active"`
}

type UpdateRequest struct {
    Name      string `json:"name"`
    Address   string `json:"address"`
    IsDefault bool   `json:"is_default"`
    IsActive  bool   `json:"is_active"`
}
