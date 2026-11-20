package orders

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
	"smarterp/backend/internal/shared/validation"
)

type Handler struct {
	Store *db.Store
	Media thumbnailSource
}

func (h Handler) transact(w http.ResponseWriter, r *http.Request, fn func(*operation) (any, error)) {
	id := r.PathValue("id")
	if id != "" && !validation.UUID(id) {
		WriteError(w, ErrInvalid)
		return
	}
	ctx := r.Context()
	tid := tenant.FromContext(ctx)
	var output any
	err := h.Store.WithTenantTx(ctx, tid, func(tx pgx.Tx) error {
		var err error
		output, err = fn(&operation{ctx: ctx, tx: tx, tenant: tid, id: id})
		return err
	})
	if err != nil {
		WriteError(w, err)
		return
	}
	h.writeCard(w, r, output)
}

func (h Handler) writeCard(w http.ResponseWriter, r *http.Request, output any) {
	if card, ok := output.(*Card); ok {
		h.thumbnails(r.Context(), tenant.FromContext(r.Context()), card)
	}
	w.Header().Set("Cache-Control", "private, no-store")
	httpx.WriteData(w, 200, output, nil)
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	q := httpx.ParseListQuery(r, httpx.SortConfig{Allowed: []string{"number", "customer_name", "status", "total_with_delivery", "created_at"}, Fallback: "created_at"}, nil)
	ctx := r.Context()
	tid := tenant.FromContext(ctx)
	var items []Summary
	var total int
	err := h.Store.WithTenantTx(ctx, tid, func(tx pgx.Tx) error {
		var err error
		items, total, err = list(ctx, tx, tid, q)
		return err
	})
	if err != nil {
		WriteError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "private, no-store")
	httpx.WriteData(w, 200, items, &httpx.Meta{Page: q.Page, PerPage: q.PerPage, Total: total})
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	h.transact(w, r, func(op *operation) (any, error) { return ReadCard(op.ctx, op.tx, op.tenant, op.id) })
}

func (h Handler) Save(w http.ResponseWriter, r *http.Request) {
	var e Entry
	if !Decode(w, r, &e) {
		return
	}
	h.transact(w, r, func(op *operation) (any, error) { return op.save(e, r.Method == http.MethodPost) })
}

func (h Handler) Quote(w http.ResponseWriter, r *http.Request) {
	var e Entry
	if !Decode(w, r, &e) {
		return
	}
	h.transact(w, r, func(op *operation) (any, error) {
		if op.id != "" {
			if err := op.prepareEdit(e.Revision); err != nil {
				return nil, err
			}
		}
		return op.quote(e)
	})
}

func Decode(w http.ResponseWriter, r *http.Request, target any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := httpx.DecodeJSON(r, target); err != nil {
		WriteError(w, ErrInvalid)
		return false
	}
	return true
}

func WriteError(w http.ResponseWriter, err error) {
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		httpx.WriteError(w, 404, "not_found", "Order or invoice not found.", nil)
	case errors.Is(err, ErrInvalid):
		httpx.WriteError(w, 422, "invalid_order", "Check the order fields, quantities and amounts.", nil)
	case errors.Is(err, ErrConflict), errors.Is(err, ErrIssued):
		httpx.WriteError(w, 409, "order_conflict", err.Error(), nil)
	case errors.Is(err, ErrCurrency), errors.Is(err, ErrSelection):
		httpx.WriteError(w, 422, "order_unavailable", err.Error(), nil)
	case errors.As(err, &pg) && (pg.Code == "23503" || pg.Code == "23514" || pg.Code == "22003"):
		httpx.WriteError(w, 422, "invalid_reference", "Check the selected records and amounts.", nil)
	default:
		httpx.WriteError(w, 500, "internal_error", "Could not complete the order operation.", nil)
	}
}
