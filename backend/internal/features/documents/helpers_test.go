package documents

import "github.com/shopspring/decimal"

const (
	uuidWarehouse  = "11111111-1111-1111-1111-111111111111"
	uuidWarehouse2 = "22222222-2222-2222-2222-222222222222"
	uuidVariant    = "33333333-3333-3333-3333-333333333333"
	uuidVariant2   = "44444444-4444-4444-4444-444444444444"
	uuidShift      = "55555555-5555-5555-5555-555555555555"
)

func dec(value string) decimal.Decimal {
	return decimal.RequireFromString(value)
}

func saleRequest() CreateRequest {
	return CreateRequest{
		Type:        "SALE",
		Date:        "2026-06-22",
		WarehouseID: uuidWarehouse,
		Items:       []ItemInput{{VariantID: uuidVariant, Qty: dec("1"), UnitPrice: dec("10")}},
		Payments:    []PaymentInput{{Method: "cash", Amount: dec("10")}},
	}
}
