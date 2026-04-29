package products

import (
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/shared/httpx"
)

type Product struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    IsComposite bool   `json:"is_composite"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type ProductListItem struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	IsComposite bool                 `json:"is_composite"`
	GlobalQty   decimal.Decimal      `json:"global_qty"`
	StockValue  decimal.Decimal      `json:"stock_value"`
	Variants    []ProductVariantItem `json:"variants"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

type ProductVariantItem struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	SKUCode    string                 `json:"sku_code"`
	Barcode    string                 `json:"barcode"`
	Option1    string                 `json:"option1"`
	Option2    string                 `json:"option2"`
	Option3    string                 `json:"option3"`
	Unit       string                 `json:"unit"`
	Price      decimal.Decimal        `json:"price"`
	GlobalQty  decimal.Decimal        `json:"global_qty"`
	AvgCost    decimal.Decimal        `json:"avg_cost"`
	StockValue decimal.Decimal        `json:"stock_value"`
	Warehouses []ProductWarehouseItem `json:"warehouses"`
}

type ProductWarehouseItem struct {
	WarehouseID string          `json:"warehouse_id"`
	Warehouse   string          `json:"warehouse"`
	Qty         decimal.Decimal `json:"qty"`
}

type ProductListInclude struct {
	Variants   bool
	Stock      bool
	Warehouses bool
}

type ProductListQuery struct {
	List  httpx.ListQuery
	Stock ProductStockFilter
}

type ProductStockFilter struct {
	WarehouseID string
	MinQty      *decimal.Decimal
	MaxQty      *decimal.Decimal
	QtyState    string
}

type CreateRequest struct {
    Name        string          `json:"name"`
    IsComposite bool            `json:"is_composite"`
    Unit        string          `json:"unit"`
    Price       decimal.Decimal `json:"price"`
    SKUCode     string          `json:"sku_code"`
    Barcode     string          `json:"barcode"`
    VariantName string          `json:"variant_name"`
}

type UpdateRequest struct {
    Name        string `json:"name"`
    IsComposite bool   `json:"is_composite"`
}
