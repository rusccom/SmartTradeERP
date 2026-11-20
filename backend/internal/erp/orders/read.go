package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/shared/db"
)

const summaryColumns = `o.id::text,o.number,o.source,o.customer_name,o.status,o.currency_code,o.total_amount,
    o.total_amount+o.delivery_cost,o.created_at`

func (s *Summary) destinations() []any {
	return []any{&s.ID, &s.Number, &s.Source, &s.Name, &s.Status, &s.Currency, &s.Total, &s.Grand, &s.Created}
}

func ReadCard(ctx context.Context, tx db.DBTX, tenantID, id string) (*Card, error) {
	const query = `SELECT ` + summaryColumns + `,o.customer_email,o.customer_phone,o.revision,o.details,o.delivery_cost,
    COALESCE((SELECT sum(p.amount) FROM sales.order_payments p WHERE p.tenant_id=o.tenant_id AND p.order_id=o.id),
        CASE WHEN o.source='manual' THEN 0 END),o.invoice_requested,COALESCE(o.billing_customer_id::text,''),o.shipping_address
    FROM storefront.orders o WHERE o.tenant_id=$1 AND o.id=$2`
	c := new(Card)
	var street string
	fields := append(c.destinations(), &c.Email, &c.Phone, &c.Revision, &c.Details, &c.DeliveryCost, &c.Paid, &c.InvoiceRequested, &c.BillingCustomerID, &street)
	if err := tx.QueryRow(ctx, query, tenantID, id).Scan(fields...); err != nil {
		return nil, err
	}
	if c.Details.Shipping.Street == "" {
		c.Details.Shipping.Street = street
	}
	return loadContents(ctx, tx, tenantID, c)
}

func loadContents(ctx context.Context, tx db.DBTX, tenantID string, c *Card) (*Card, error) {
	var err error
	if c.Items, err = loadLines(ctx, tx, tenantID, c.ID); err != nil {
		return nil, err
	}
	for _, line := range c.Items {
		c.Quantity = c.Quantity.Add(line.Qty)
	}
	c.Invoice, err = loadInvoice(ctx, tx, tenantID, c.ID)
	return c, err
}

func loadLines(ctx context.Context, tx db.DBTX, tenantID, id string) ([]Line, error) {
	const q = `SELECT i.id::text,i.variant_id::text,v.product_id::text,i.qty,COALESCE(i.vat_rate,''),i.product_name,i.variant_name,
        COALESCE(i.sku_code,v.sku_code,''),COALESCE(i.unit,v.unit,''),i.unit_price,i.total_amount FROM storefront.order_items i
        JOIN catalog.product_variants v ON (v.tenant_id,v.id)=(i.tenant_id,i.variant_id)
        WHERE i.tenant_id=$1 AND i.order_id=$2 ORDER BY i.position NULLS LAST,i.id`
	rows, err := tx.Query(ctx, q, tenantID, id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Line, error) {
		var l Line
		err := row.Scan(&l.ID, &l.VariantID, &l.ProductID, &l.Qty, &l.VAT, &l.ProductName, &l.VariantName, &l.SKU, &l.Unit, &l.Price, &l.Total)
		return l, err
	})
}

func loadInvoice(ctx context.Context, tx db.DBTX, tenantID, id string) (*Invoice, error) {
	var i Invoice
	err := tx.QueryRow(ctx, `SELECT id::text,number,format,issued_at FROM sales.invoices WHERE tenant_id=$1 AND order_id=$2`, tenantID, id).Scan(&i.ID, &i.Number, &i.Format, &i.IssuedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &i, err
}

func Lock(ctx context.Context, tx db.DBTX, tenantID, id string) error {
	var rev int64
	return tx.QueryRow(ctx, `SELECT revision FROM storefront.orders WHERE tenant_id=$1 AND id=$2 FOR UPDATE`, tenantID, id).Scan(&rev)
}
