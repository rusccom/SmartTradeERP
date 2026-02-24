package reports

import (
    "net/http"
    "time"

    "smarterp/backend/internal/shared/httpx"
    "smarterp/backend/internal/shared/tenant"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Profit(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    fromDate, toDate := parsePeriod(r)
    warehouseID := r.URL.Query().Get("warehouse_id")
    variantID := r.URL.Query().Get("variant_id")
    profit, err := h.service.Profit(r.Context(), tenantID, fromDate, toDate, warehouseID, variantID)
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to load profit", err.Error())
        return
    }
    httpx.WriteData(w, http.StatusOK, ProfitReport{Profit: profit}, nil)
}

func (h *Handler) Stock(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    warehouseID := r.URL.Query().Get("warehouse_id")
    rows, err := h.service.Stock(r.Context(), tenantID, warehouseID)
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to load stock", err.Error())
        return
    }
    httpx.WriteData(w, http.StatusOK, rows, nil)
}

func (h *Handler) TopProducts(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    fromDate, toDate := parsePeriod(r)
    rows, err := h.service.TopProducts(r.Context(), tenantID, fromDate, toDate)
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to load top products", err.Error())
        return
    }
    httpx.WriteData(w, http.StatusOK, rows, nil)
}

func (h *Handler) Movements(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    variantID := r.URL.Query().Get("variant_id")
    rows, err := h.service.Movements(r.Context(), tenantID, variantID)
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to load movements", err.Error())
        return
    }
    httpx.WriteData(w, http.StatusOK, rows, nil)
}

func parsePeriod(r *http.Request) (time.Time, time.Time) {
    fromDate := parseDateOrDefault(r.URL.Query().Get("from"), time.Now().AddDate(0, -1, 0))
    toDate := parseDateOrDefault(r.URL.Query().Get("to"), time.Now())
    return fromDate, toDate
}

func parseDateOrDefault(raw string, fallback time.Time) time.Time {
    if raw == "" {
        return fallback
    }
    value, err := time.Parse("2006-01-02", raw)
    if err != nil {
        return fallback
    }
    return value
}
