package storefront

import (
	"context"
	"encoding/json"

	"github.com/shopspring/decimal"

	"smarterp/backend/internal/erp/orders"
	"smarterp/backend/internal/erp/productvat"
	"smarterp/backend/internal/platform/storefront/catalog"
	"smarterp/backend/internal/shared/db"
)

type orderRecord struct {
	ID            string
	TenantID      string
	Number        string
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	Delivery      orders.Details
	Note          string
	CurrencyCode  string
	Total         decimal.Decimal
}

type orderItemRecord struct {
	VAT         string
	Components  []catalog.Component
	ID          string
	VariantID   string
	ProductName string
	VariantName string
	UnitPrice   decimal.Decimal
	Qty         decimal.Decimal
	Total       decimal.Decimal
}

const insertOrderSQL = `INSERT INTO storefront.orders
        (id, tenant_id, number, status, customer_name, customer_email, customer_phone, shipping_address, currency_code, total_amount, note, details)
    VALUES ($1, $2, $3, 'new', $4, $5, $6, $7, $8, $9, $10, $11)`

const insertOrderItemSQL = `INSERT INTO storefront.order_items
        (id, order_id, variant_id, product_name, variant_name, unit_price, qty, total_amount, tenant_id, bundle_components, vat_rate, vat_country, vat_mode)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NULLIF($11,''), $12, 'catalog')`

func (r *Repository) InsertOrder(ctx context.Context, tx db.DBTX, order orderRecord, items []orderItemRecord) error {
	if err := insertOrderRow(ctx, tx, order); err != nil {
		return err
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.VariantID)
	}
	rates, err := productvat.ForVariants(ctx, tx, order.TenantID, order.Delivery.Destination(), ids)
	if err != nil {
		return err
	}
	for _, item := range items {
		item.VAT = rates[item.VariantID]
		if err := insertOrderItem(ctx, tx, order, item); err != nil {
			return err
		}
	}
	return nil
}

func insertOrderRow(ctx context.Context, tx db.DBTX, o orderRecord) error {
	_, err := tx.Exec(ctx, insertOrderSQL, o.ID, o.TenantID, o.Number, o.CustomerName,
		o.CustomerEmail, o.CustomerPhone, o.Delivery.Shipping.Street, o.CurrencyCode, o.Total, o.Note, o.Delivery)
	return err
}

func insertOrderItem(ctx context.Context, tx db.DBTX, order orderRecord, item orderItemRecord) error {
	if item.Components == nil {
		item.Components = []catalog.Component{}
	}
	components, err := json.Marshal(item.Components)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, insertOrderItemSQL, item.ID, order.ID, item.VariantID,
		item.ProductName, item.VariantName, item.UnitPrice, item.Qty, item.Total, order.TenantID, components, item.VAT, order.Delivery.Destination())
	return err
}
