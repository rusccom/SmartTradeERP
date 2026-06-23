package storefront

import (
	"errors"
	"net/http"

	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
	"smarterp/backend/internal/shared/validation"
)

type CommerceHandler struct {
	commerce *Commerce
}

func NewCommerceHandler(commerce *Commerce) *CommerceHandler {
	return &CommerceHandler{commerce: commerce}
}

type cartPayload struct {
	Items []cartItemInput `json:"items"`
}

func (h *CommerceHandler) Cart(w http.ResponseWriter, r *http.Request) {
	payload := cartPayload{}
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
		return
	}
	result, err := h.commerce.PriceCart(r.Context(), tenant.FromContext(r.Context()), payload.Items)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to price cart", err.Error())
		return
	}
	httpx.WriteData(w, http.StatusOK, result, nil)
}

func (h *CommerceHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	req := checkoutRequest{}
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
		return
	}
	result, err := h.commerce.Checkout(r.Context(), tenant.FromContext(r.Context()), req)
	if err != nil {
		writeCheckoutError(w, err)
		return
	}
	httpx.WriteData(w, http.StatusCreated, result, nil)
}

func writeCheckoutError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrCheckoutInvalid), errors.Is(err, validation.ErrInvalidData):
		httpx.WriteError(w, http.StatusBadRequest, "invalid_data", "invalid checkout data", nil)
	case errors.Is(err, ErrCartEmpty):
		httpx.WriteError(w, http.StatusBadRequest, "empty_cart", "no purchasable items", nil)
	default:
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "checkout failed", err.Error())
	}
}
