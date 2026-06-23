package storefront

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/validation"
)

const maxLineQty = 100000

// Commerce prices carts and places orders. All prices come from the database;
// the client is trusted only for variant ids and quantities.
type Commerce struct {
	store *db.Store
	repo  *Repository
}

func NewCommerce(store *db.Store, repo *Repository) *Commerce {
	return &Commerce{store: store, repo: repo}
}

func (c *Commerce) PriceCart(ctx context.Context, tenantID string, items []cartItemInput) (cartResult, error) {
	cur, err := c.repo.TenantCurrency(ctx, tenantID)
	if err != nil {
		return cartResult{}, err
	}
	lines, total, err := c.buildLines(ctx, tenantID, cur, items)
	if err != nil {
		return cartResult{}, err
	}
	return cartResult{Lines: lines, Total: formatMoney(total, cur).Display, Currency: cur.code}, nil
}

func (c *Commerce) buildLines(ctx context.Context, tenantID string, cur currency, items []cartItemInput) ([]cartLine, decimal.Decimal, error) {
	lines := make([]cartLine, 0, len(items))
	total := decimal.Zero
	for _, item := range items {
		line, lineTotal, ok, err := c.lineFor(ctx, tenantID, cur, item)
		if err != nil {
			return nil, decimal.Zero, err
		}
		if ok {
			lines = append(lines, line)
			total = total.Add(lineTotal)
		}
	}
	return lines, total, nil
}

func (c *Commerce) lineFor(ctx context.Context, tenantID string, cur currency, item cartItemInput) (cartLine, decimal.Decimal, bool, error) {
	if !validQty(item.Qty) {
		return cartLine{}, decimal.Zero, false, nil
	}
	price, found, err := c.repo.PriceVariant(ctx, tenantID, item.VariantID)
	if err != nil || !found {
		return cartLine{}, decimal.Zero, false, err
	}
	stock, err := c.repo.VariantStock(ctx, tenantID, item.VariantID)
	if err != nil {
		return cartLine{}, decimal.Zero, false, err
	}
	lineTotal := price.Price.Mul(item.Qty).Round(4)
	return buildCartLine(item.VariantID, price, item.Qty, lineTotal, stock, cur), lineTotal, true, nil
}

func buildCartLine(variantID string, price variantPrice, qty, lineTotal, stock decimal.Decimal, cur currency) cartLine {
	return cartLine{
		VariantID:   variantID,
		ProductName: price.ProductName,
		VariantName: price.VariantName,
		Qty:         qty.String(),
		UnitPrice:   formatMoney(price.Price, cur).Display,
		LineTotal:   formatMoney(lineTotal, cur).Display,
		Available:   stock.String(),
		Purchasable: stock.GreaterThanOrEqual(qty),
	}
}

// Checkout re-prices every line server-side and persists an order. Stock is
// advisory (returned in cart preview); the merchant fulfils from the ERP.
func (c *Commerce) Checkout(ctx context.Context, tenantID string, req checkoutRequest) (orderResult, error) {
	if err := validateCustomer(req.Customer); err != nil {
		return orderResult{}, err
	}
	cur, err := c.repo.TenantCurrency(ctx, tenantID)
	if err != nil {
		return orderResult{}, err
	}
	items, total, err := c.orderItems(ctx, tenantID, req.Items)
	if err != nil {
		return orderResult{}, err
	}
	if len(items) == 0 {
		return orderResult{}, ErrCartEmpty
	}
	return c.persistOrder(ctx, tenantID, cur, req, items, total)
}

func (c *Commerce) orderItems(ctx context.Context, tenantID string, items []cartItemInput) ([]orderItemRecord, decimal.Decimal, error) {
	records := make([]orderItemRecord, 0, len(items))
	total := decimal.Zero
	for _, item := range items {
		record, ok, err := c.orderItem(ctx, tenantID, item)
		if err != nil {
			return nil, decimal.Zero, err
		}
		if ok {
			records = append(records, record)
			total = total.Add(record.Total)
		}
	}
	return records, total, nil
}

func (c *Commerce) orderItem(ctx context.Context, tenantID string, item cartItemInput) (orderItemRecord, bool, error) {
	if !validQty(item.Qty) {
		return orderItemRecord{}, false, nil
	}
	price, found, err := c.repo.PriceVariant(ctx, tenantID, item.VariantID)
	if err != nil || !found {
		return orderItemRecord{}, false, err
	}
	total := price.Price.Mul(item.Qty).Round(4)
	return orderItemRecord{
		ID: uuid.NewString(), VariantID: item.VariantID, ProductName: price.ProductName,
		VariantName: price.VariantName, UnitPrice: price.Price, Qty: item.Qty, Total: total,
	}, true, nil
}

func (c *Commerce) persistOrder(ctx context.Context, tenantID string, cur currency, req checkoutRequest, items []orderItemRecord, total decimal.Decimal) (orderResult, error) {
	order := buildOrderRecord(tenantID, cur, req, total)
	err := c.store.WithTx(ctx, func(tx pgx.Tx) error {
		return c.repo.InsertOrder(ctx, tx, order, items)
	})
	if err != nil {
		return orderResult{}, err
	}
	return orderResult{Number: order.Number, Total: formatMoney(total, cur).Display, Currency: cur.code, Status: "new"}, nil
}

func buildOrderRecord(tenantID string, cur currency, req checkoutRequest, total decimal.Decimal) orderRecord {
	return orderRecord{
		ID: uuid.NewString(), TenantID: tenantID, Number: orderNumber(),
		CustomerName:  strings.TrimSpace(req.Customer.Name),
		CustomerEmail: strings.TrimSpace(req.Customer.Email),
		CustomerPhone: strings.TrimSpace(req.Customer.Phone),
		Address:       strings.TrimSpace(req.Customer.Address),
		Note:          strings.TrimSpace(req.Note),
		CurrencyCode:  cur.code, Total: total,
	}
}

func orderNumber() string {
	return "SF-" + strings.ToUpper(uuid.NewString()[:8])
}

func validateCustomer(cust checkoutCustomer) error {
	name := strings.TrimSpace(cust.Name)
	email := strings.TrimSpace(cust.Email)
	if name == "" || len([]rune(name)) > 160 || email == "" || !validation.Email(email) {
		return ErrCheckoutInvalid
	}
	return nil
}

func validQty(qty decimal.Decimal) bool {
	return qty.GreaterThan(decimal.Zero) && qty.LessThanOrEqual(decimal.NewFromInt(maxLineQty))
}
