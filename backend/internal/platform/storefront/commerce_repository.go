package storefront

import (
	"context"
	"encoding/json"

	"github.com/shopspring/decimal"

	"smarterp/backend/internal/integrations/smartbundle"
	"smarterp/backend/internal/shared/db"
)

type orderRecord struct {
	ID            string
	TenantID      string
	Number        string
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	Address       string
	Note          string
	CurrencyCode  string
	Total         decimal.Decimal
}

type orderItemRecord struct {
	Components  []smartbundle.Component
	ID          string
	VariantID   string
	ProductName string
	VariantName string
	UnitPrice   decimal.Decimal
	Qty         decimal.Decimal
	Total       decimal.Decimal
}

const insertOrderSQL = `INSERT INTO storefront.orders
        (id, tenant_id, number, status, customer_name, customer_email, customer_phone, shipping_address, currency_code, total_amount, note)
    VALUES ($1, $2, $3, 'new', $4, $5, $6, $7, $8, $9, $10)`

const insertOrderItemSQL = `INSERT INTO storefront.order_items
        (id, order_id, variant_id, product_name, variant_name, unit_price, qty, total_amount, tenant_id, bundle_components)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

func (r *Repository) InsertOrder(ctx context.Context, tx db.DBTX, order orderRecord, items []orderItemRecord) error {
	if err := insertOrderRow(ctx, tx, order); err != nil {
		return err
	}
	for _, item := range items {
		if err := insertOrderItem(ctx, tx, order, item); err != nil {
			return err
		}
	}
	return nil
}

func insertOrderRow(ctx context.Context, tx db.DBTX, o orderRecord) error {
	_, err := tx.Exec(ctx, insertOrderSQL, o.ID, o.TenantID, o.Number, o.CustomerName,
		o.CustomerEmail, o.CustomerPhone, o.Address, o.CurrencyCode, o.Total, o.Note)
	return err
}

func insertOrderItem(ctx context.Context, tx db.DBTX, order orderRecord, item orderItemRecord) error {
	if item.Components == nil {
		item.Components = []smartbundle.Component{}
	}
	components, err := json.Marshal(item.Components)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, insertOrderItemSQL, item.ID, order.ID, item.VariantID,
		item.ProductName, item.VariantName, item.UnitPrice, item.Qty, item.Total, order.TenantID, components)
	return err
}
