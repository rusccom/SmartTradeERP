package storefront

import (
	"errors"
	"net/http"

	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/validation"
)

type cartPayload struct {
	Items []cartItemInput `json:"items"`
}

func (h *Handler) quote(w http.ResponseWriter, r *http.Request, sh *shop) {
	var payload cartPayload
	if !readPayload(w, r, &payload) {
		return
	}
	cart, err := sh.priceCart(payload.Items)
	if err != nil {
		writeCheckoutError(w, err)
		return
	}
	cur := sh.brand.Currency
	httpx.WriteData(w, 200, cartResult{Lines: cart.lines, Total: formatMoney(cart.total, cur).Display, Currency: cur.code}, nil)
}

func (h *Handler) checkout(w http.ResponseWriter, r *http.Request, sh *shop) {
	var req checkoutRequest
	if !readPayload(w, r, &req) {
		return
	}
	result, err := sh.place(req)
	if err != nil {
		writeCheckoutError(w, err)
		return
	}
	status := http.StatusCreated
	if sh.preview {
		status = http.StatusOK
	}
	httpx.WriteData(w, status, result, nil)
}

func readPayload(w http.ResponseWriter, r *http.Request, target any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 256<<10)
	if err := httpx.DecodeJSON(r, target); err != nil {
		httpx.WriteError(w, 400, "bad_request", "invalid payload", nil)
		return false
	}
	return true
}

func writeCheckoutError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrCheckoutInvalid), errors.Is(err, validation.ErrInvalidData):
		httpx.WriteError(w, 400, "invalid_data", "Invalid checkout data", nil)
	case errors.Is(err, ErrCartEmpty):
		httpx.WriteError(w, 400, "empty_cart", "Cart is empty", nil)
	case errors.Is(err, errStock):
		httpx.WriteError(w, 409, "unavailable", "Some items are unavailable", nil)
	default:
		httpx.WriteError(w, 500, "internal_error", "Store unavailable", nil)
	}
}
