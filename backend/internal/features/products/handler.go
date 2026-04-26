package products

import (
    "errors"
    "net/http"

    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/shared/httpx"
    "smarterp/backend/internal/shared/tenant"
    "smarterp/backend/internal/shared/validation"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.FromContext(r.Context())
	query, err := productListQuery(r)
	if err != nil {
		h.writeListError(w, err)
		return
	}
	includes := productIncludes(r)
	if includes.Variants {
		h.listWithIncludes(w, r, tenantID, query, includes)
		return
	}
	data, total, err := h.service.List(r.Context(), tenantID, query)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list products", err.Error())
		return
	}
	meta := &httpx.Meta{Page: query.List.Page, PerPage: query.List.PerPage, Total: total}
	httpx.WriteData(w, http.StatusOK, data, meta)
}

func (h *Handler) listWithIncludes(
	w http.ResponseWriter,
	r *http.Request,
	tenantID string,
	query ProductListQuery,
	includes ProductListInclude,
) {
	data, total, err := h.service.ListWithIncludes(r.Context(), tenantID, query, includes)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list products", err.Error())
		return
	}
	meta := &httpx.Meta{Page: query.List.Page, PerPage: query.List.PerPage, Total: total}
	httpx.WriteData(w, http.StatusOK, data, meta)
}

func (h *Handler) writeListError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrInvalidProductFilter) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_filter", "invalid product filter", nil)
		return
	}
	httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list products", err.Error())
}

func productIncludes(r *http.Request) ProductListInclude {
	values := httpx.ParseIncludes(r)
	include := ProductListInclude{Variants: httpx.HasInclude(values, "variants")}
	include.Stock = httpx.HasInclude(values, "stock")
	include.Warehouses = httpx.HasInclude(values, "warehouses")
	return normalizeProductIncludes(include)
}

func normalizeProductIncludes(include ProductListInclude) ProductListInclude {
	if include.Stock || include.Warehouses {
		include.Variants = true
	}
	if include.Warehouses {
		include.Stock = true
	}
	return include
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    req := CreateRequest{}
    if err := httpx.DecodeJSON(r, &req); err != nil {
        httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
        return
    }
    tenantID := tenant.FromContext(r.Context())
    id, err := h.service.Create(r.Context(), tenantID, req)
    if err != nil {
        h.writeMutationError(w, err, "failed to create product")
        return
    }
    httpx.WriteData(w, http.StatusCreated, map[string]string{"id": id}, nil)
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    item, err := h.service.ByID(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeQueryError(w, err)
        return
    }
    httpx.WriteData(w, http.StatusOK, item, nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    req := UpdateRequest{}
    if err := httpx.DecodeJSON(r, &req); err != nil {
        httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
        return
    }
    tenantID := tenant.FromContext(r.Context())
    err := h.service.Update(r.Context(), tenantID, r.PathValue("id"), req)
    if err != nil {
        h.writeMutationError(w, err, "failed to update product")
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "updated"}, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    err := h.service.Delete(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeDeleteError(w, err)
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "deleted"}, nil)
}

func (h *Handler) writeQueryError(w http.ResponseWriter, err error) {
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "product not found", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to read product", err.Error())
}

func (h *Handler) writeDeleteError(w http.ResponseWriter, err error) {
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "product not found", nil)
        return
    }
    if errors.Is(err, ErrHasMovements) {
        httpx.WriteError(w, http.StatusConflict, "has_movements", "product has ledger movements", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to delete product", err.Error())
}

func (h *Handler) writeMutationError(w http.ResponseWriter, err error, message string) {
    if errors.Is(err, validation.ErrInvalidData) {
        httpx.WriteError(w, http.StatusBadRequest, "invalid_data", "invalid product data", nil)
        return
    }
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "product not found", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}
