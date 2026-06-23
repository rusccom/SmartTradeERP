package storefront

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/shared/db"
)

type variantPrice struct {
	ProductName string
	VariantName string
	Price       decimal.Decimal
}

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
	ID          string
	VariantID   string
	ProductName string
	VariantName string
	UnitPrice   decimal.Decimal
	Qty         decimal.Decimal
	Total       decimal.Decimal
}

const priceVariantSQL = `SELECT p.name, COALESCE(v.name, 'Default'), COALESCE(v.price, 0)
    FROM catalog.product_variants v
    JOIN catalog.products p ON p.id = v.product_id
    WHERE v.tenant_id = $1 AND v.id = $2 AND v.is_published AND p.is_published
    LIMIT 1`

// PriceVariant returns the server-authoritative price for a published variant.
// found is false when the variant is missing, unpublished, or cross-tenant.
func (r *Repository) PriceVariant(ctx context.Context, tenantID, variantID string) (variantPrice, bool, error) {
	out := variantPrice{}
	err := r.store.Pool.QueryRow(ctx, priceVariantSQL, tenantID, variantID).Scan(&out.ProductName, &out.VariantName, &out.Price)
	if errors.Is(err, pgx.ErrNoRows) {
		return variantPrice{}, false, nil
	}
	if err != nil {
		return variantPrice{}, false, err
	}
	return out, true, nil
}

const variantStockSQL = `SELECT COALESCE(running_qty, 0)
    FROM ledger.cost_movement_results r
    JOIN ledger.inventory_movements m ON m.id = r.movement_id
    JOIN ledger.posting_batches b ON b.id = m.posting_batch_id
    WHERE r.tenant_id = $1 AND r.variant_id = $2 AND b.status = 'active'
    ORDER BY r.sequence_num DESC LIMIT 1`

func (r *Repository) VariantStock(ctx context.Context, tenantID, variantID string) (decimal.Decimal, error) {
	qty := decimal.Zero
	err := r.store.Pool.QueryRow(ctx, variantStockSQL, tenantID, variantID).Scan(&qty)
	if errors.Is(err, pgx.ErrNoRows) {
		return decimal.Zero, nil
	}
	return qty, err
}

const tenantCurrencySQL = `SELECT COALESCE(c.code, ''), COALESCE(c.symbol, ''), COALESCE(c.decimal_places, 2)
    FROM platform.tenant_settings ts
    LEFT JOIN platform.currencies c ON c.id = ts.base_currency_id
    WHERE ts.tenant_id = $1`

func (r *Repository) TenantCurrency(ctx context.Context, tenantID string) (currency, error) {
	var code, symbol string
	var places int
	err := r.store.Pool.QueryRow(ctx, tenantCurrencySQL, tenantID).Scan(&code, &symbol, &places)
	if errors.Is(err, pgx.ErrNoRows) {
		return currency{places: 2}, nil
	}
	if err != nil {
		return currency{}, err
	}
	return currency{code: code, symbol: symbol, places: int32(places)}, nil
}

const insertOrderSQL = `INSERT INTO storefront.orders
        (id, tenant_id, number, status, customer_name, customer_email, customer_phone, shipping_address, currency_code, total_amount, note)
    VALUES ($1, $2, $3, 'new', $4, $5, $6, $7, $8, $9, $10)`

const insertOrderItemSQL = `INSERT INTO storefront.order_items
        (id, order_id, variant_id, product_name, variant_name, unit_price, qty, total_amount)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

func (r *Repository) InsertOrder(ctx context.Context, tx db.DBTX, order orderRecord, items []orderItemRecord) error {
	if err := insertOrderRow(ctx, tx, order); err != nil {
		return err
	}
	for _, item := range items {
		if err := insertOrderItem(ctx, tx, order.ID, item); err != nil {
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

func insertOrderItem(ctx context.Context, tx db.DBTX, orderID string, item orderItemRecord) error {
	_, err := tx.Exec(ctx, insertOrderItemSQL, item.ID, orderID, item.VariantID,
		item.ProductName, item.VariantName, item.UnitPrice, item.Qty, item.Total)
	return err
}
