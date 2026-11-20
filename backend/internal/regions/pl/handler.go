package pl

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/erp/orders"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
)

type problem string

func (p problem) Error() string { return string(p) }

type Handler struct {
	store    *db.Store
	key      string
	sessions sessionCache
}

func New(store *db.Store, key string) *Handler { return &Handler{store: store, key: key} }

type scope struct {
	ctx    context.Context
	tenant string
	tx     pgx.Tx
}

func (h *Handler) run(ctx context.Context, fn func(scope) error) error {
	tid := tenant.FromContext(ctx)
	return h.store.WithTenantTx(ctx, tid, func(tx pgx.Tx) error { return fn(scope{ctx, tid, tx}) })
}

func tenantResult[T any](h *Handler, ctx context.Context, fn func(scope) (T, error)) (T, error) {
	var value T
	err := h.run(ctx, func(s scope) error {
		var err error
		value, err = fn(s)
		return err
	})
	return value, err
}

func respond(w http.ResponseWriter, data any, err error) {
	w.Header().Set("Cache-Control", "private, no-store")
	if err == nil {
		httpx.WriteData(w, 200, data, nil)
		return
	}
	var p problem
	if errors.As(err, &p) {
		httpx.WriteError(w, 422, "invoice_configuration", p.Error(), nil)
		return
	}
	orders.WriteError(w, err)
}

func (h *Handler) Settings(w http.ResponseWriter, r *http.Request) {
	profile, err := tenantResult(h, r.Context(), func(s scope) (*Seller, error) { return readSeller(s.ctx, s.tx, s.tenant) })
	data := map[string]any{"seller": profile, "can_manage": auth.ClaimsFromContext(r.Context()).Role == "owner"}
	respond(w, data, err)
}

func (h *Handler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	var u profileUpdate
	if !orders.Decode(w, r, &u) {
		return
	}
	profile, err := tenantResult(h, r.Context(), func(s scope) (*Seller, error) { return h.saveProfile(s.ctx, s.tx, s.tenant, u) })
	respond(w, profile, err)
}

func (h *Handler) Formats(w http.ResponseWriter, r *http.Request) {
	formats := []map[string]string{{"id": standardFormat, "label": "Polish domestic VAT · PLN · FA(3)"}}
	respond(w, formats, nil)
}
