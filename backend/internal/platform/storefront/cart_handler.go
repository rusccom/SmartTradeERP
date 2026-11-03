package storefront

import (
	"maps"
	"net/http"
	"slices"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/validation"
)

func (h *Handler) currentCart(w http.ResponseWriter, r *http.Request, sh *shop) {
	cart, err := sh.readCart(cartSession(w, r, sh.preview))
	if err != nil {
		writeCheckoutError(w, err)
		return
	}
	quote, err := sh.priceCart(cart.Items)
	if err != nil {
		writeCheckoutError(w, err)
		return
	}
	sh.writeCart(w, quote, cart.Version)
}

func (h *Handler) updateCart(w http.ResponseWriter, r *http.Request, sh *shop) {
	var cart savedCart
	if !readPayload(w, r, &cart) {
		return
	}
	id := cartSession(w, r, sh.preview)
	quote, err := sh.replaceCart(id, cart)
	if err != nil {
		writeCheckoutError(w, err)
		return
	}
	sh.writeCart(w, quote, cart.Version+1)
}

func (sh *shop) replaceCart(id string, cart savedCart) (pricedCart, error) {
	if cart.Items == nil || cart.Version < 0 {
		return pricedCart{}, validation.ErrInvalidData
	}
	quote, err := sh.priceCart(cart.Items)
	if err != nil {
		return quote, err
	}
	if err = sh.checkUnavailable(id, cart, quote.lines); err != nil {
		return quote, err
	}
	err = sh.service.repo.store.WithTenantTx(sh.ctx, sh.tenantID, func(tx pgx.Tx) error {
		return sh.saveCart(tx, id, cart)
	})
	return quote, err
}

func (sh *shop) checkUnavailable(id string, cart savedCart, lines []cartLine) error {
	if !slices.ContainsFunc(lines, func(l cartLine) bool { return l.UnitPrice == "" }) {
		return nil
	}
	previous, err := sh.readCart(id)
	if err != nil {
		return err
	}
	for i, line := range lines {
		if line.UnitPrice == "" && !cartContained(previous.Items, cart.Items[i]) {
			return validation.ErrInvalidData
		}
	}
	return nil
}

func cartContained(items []cartItemInput, wanted cartItemInput) bool {
	return slices.ContainsFunc(items, func(item cartItemInput) bool {
		return item.VariantID == wanted.VariantID && maps.Equal(item.Picks, wanted.Picks) && item.Qty.GreaterThanOrEqual(wanted.Qty)
	})
}

func (sh *shop) writeCart(w http.ResponseWriter, cart pricedCart, version int64) {
	cur := sh.brand.Currency
	httpx.WriteData(w, 200, cartResult{Lines: cart.lines, Total: formatMoney(cart.total, cur).Display, Currency: cur.code, Version: version}, nil)
}
