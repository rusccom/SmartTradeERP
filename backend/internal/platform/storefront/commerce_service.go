package storefront

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/integrations/smartbundle"
	"smarterp/backend/internal/platform/storefront/catalog"
	"smarterp/backend/internal/shared/validation"
)

const maxLineQty = 100000

type pricedCart struct {
	lines   []cartLine
	items   []orderItemRecord
	total   decimal.Decimal
	blocked bool
}

type catalogEntry struct {
	name    string
	variant catalog.Variant
}

func (sh *shop) priceCart(input []cartItemInput) (pricedCart, error) {
	if len(input) == 0 {
		return pricedCart{lines: []cartLine{}}, nil
	}
	ids, err := cartIDs(input)
	if err != nil {
		return pricedCart{}, err
	}
	products, err := sh.catalog(catalog.Filter{VariantIDs: ids})
	if err != nil {
		return pricedCart{}, err
	}
	return sh.quoteItems(input, indexCatalog(products)), nil
}

func (sh *shop) quoteItems(input []cartItemInput, index map[string]catalogEntry) pricedCart {
	cart := pricedCart{lines: make([]cartLine, 0, len(input))}
	needs := make([][]smartbundle.Component, 0, len(input))
	for _, item := range input {
		line, record, components := sh.priceLine(item, index[item.VariantID])
		cart.lines = append(cart.lines, line)
		cart.items = append(cart.items, record)
		cart.total = cart.total.Add(record.Total)
		needs = append(needs, components)
	}
	cart.checkStock(needs)
	return cart
}

func cartIDs(items []cartItemInput) ([]string, error) {
	if len(items) > 300 {
		return nil, validation.ErrInvalidData
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		if !validation.UUID(item.VariantID) || !item.Qty.IsPositive() || item.Qty.GreaterThan(decimal.NewFromInt(maxLineQty)) {
			return nil, validation.ErrInvalidData
		}
		ids = append(ids, item.VariantID)
	}
	return ids, nil
}

func indexCatalog(products []catalog.Product) map[string]catalogEntry {
	out := make(map[string]catalogEntry)
	for _, product := range products {
		for _, v := range product.Variants {
			out[v.ID] = catalogEntry{product.Name, v}
		}
	}
	return out
}

func (sh *shop) priceLine(item cartItemInput, entry catalogEntry) (cartLine, orderItemRecord, []smartbundle.Component) {
	v, cur := entry.variant, sh.brand.Currency
	line := cartLine{VariantID: item.VariantID, Picks: item.Picks, Qty: item.Qty.String(), ProductName: entry.name, VariantName: v.Name}
	q, err := resolveVariant(v, item.Picks, cur.places)
	if err != nil {
		return line, orderItemRecord{}, nil
	}
	total := q.Price.Mul(item.Qty).Round(cur.places)
	line.UnitPrice, line.LineTotal = formatMoney(q.Price, cur).Display, formatMoney(total, cur).Display
	line.Available = availableUnits(q.Components).String()
	line.Purchasable = true
	if v.Bundle != nil {
		line.Components = q.Components
	}
	record := orderItemRecord{ID: uuid.NewString(), VariantID: v.ID, ProductName: entry.name, VariantName: v.Name,
		Qty: item.Qty, UnitPrice: q.Price, Total: total, Components: line.Components}
	return line, record, q.Components
}

func resolveVariant(v catalog.Variant, picks map[string]string, places int32) (smartbundle.Quote, error) {
	if v.Bundle != nil {
		return v.Bundle.Resolve(v.Price, picks, places)
	}
	if v.Price == nil || len(picks) > 0 {
		return smartbundle.Quote{}, smartbundle.ErrSelection
	}
	return smartbundle.Quote{Price: *v.Price, Components: []smartbundle.Component{
		{VariantID: v.ID, Qty: decimal.NewFromInt(1), Stock: v.Stock},
	}}, nil
}

func (cart *pricedCart) checkStock(needs [][]smartbundle.Component) {
	demand := map[string]decimal.Decimal{}
	for i, components := range needs {
		for _, c := range components {
			demand[c.VariantID] = demand[c.VariantID].Add(c.Qty.Mul(cart.items[i].Qty))
		}
	}
	for i, components := range needs {
		for _, c := range components {
			if demand[c.VariantID].GreaterThan(c.Stock) {
				cart.lines[i].Purchasable = false
			}
		}
		if !cart.lines[i].Purchasable {
			cart.blocked = true
		}
	}
}

func availableUnits(components []smartbundle.Component) decimal.Decimal {
	amounts := map[string]decimal.Decimal{}
	for _, c := range components {
		amounts[c.VariantID] = amounts[c.VariantID].Add(c.Qty)
	}
	var result decimal.Decimal
	for i, c := range components {
		units := c.Stock.Div(amounts[c.VariantID])
		if i == 0 || units.LessThan(result) {
			result = units
		}
	}
	return result
}

var errStock = errors.New("cart contains unavailable items")

func (sh *shop) place(id string, req checkoutRequest) (orderResult, error) {
	if err := validateCustomer(req.Customer); err != nil {
		return orderResult{}, err
	}
	saved, err := sh.readCart(id)
	if err != nil {
		return orderResult{}, err
	}
	if saved.Version != req.Version {
		return orderResult{}, errCartChanged
	}
	cart, err := sh.checkoutCart(saved.Items)
	if err != nil {
		return orderResult{}, err
	}
	return sh.finishOrder(id, req, cart)
}

func (sh *shop) checkoutCart(items []cartItemInput) (pricedCart, error) {
	if len(items) == 0 {
		return pricedCart{}, ErrCartEmpty
	}
	cart, err := sh.priceCart(items)
	if err == nil && cart.blocked {
		err = errStock
	}
	return cart, err
}

func (sh *shop) finishOrder(id string, req checkoutRequest, cart pricedCart) (orderResult, error) {
	cur := sh.brand.Currency
	result := orderResult{Total: formatMoney(cart.total, cur).Display, Currency: cur.code, Status: "preview"}
	err := sh.service.repo.store.WithTenantTx(sh.ctx, sh.tenantID, func(tx pgx.Tx) error {
		err := sh.saveCart(tx, id, savedCart{Version: req.Version, Items: []cartItemInput{}})
		if err != nil || sh.preview {
			return err
		}
		order := sh.order(req, cart.total)
		result.Number, result.Status = order.Number, "new"
		return sh.service.repo.InsertOrder(sh.ctx, tx, order, cart.items)
	})
	return result, err
}

func (sh *shop) order(req checkoutRequest, total decimal.Decimal) orderRecord {
	return orderRecord{
		ID: uuid.NewString(), TenantID: sh.tenantID, Number: "SF-" + strings.ToUpper(uuid.NewString()[:8]),
		CustomerName: strings.TrimSpace(req.Customer.Name), CustomerEmail: strings.TrimSpace(req.Customer.Email),
		CustomerPhone: strings.TrimSpace(req.Customer.Phone), Address: strings.TrimSpace(req.Customer.Address),
		Note: strings.TrimSpace(req.Note), CurrencyCode: sh.brand.Currency.code, Total: total,
	}
}

func validateCustomer(c checkoutCustomer) error {
	name := strings.TrimSpace(c.Name)
	if name == "" || len([]rune(name)) > 160 || !validation.Email(strings.TrimSpace(c.Email)) {
		return ErrCheckoutInvalid
	}
	return nil
}
