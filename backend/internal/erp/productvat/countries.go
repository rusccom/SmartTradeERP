package productvat

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"smarterp/backend/internal/shared/countries"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
)

type Countries struct{ Store *db.Store }

func (h Countries) List(w http.ResponseWriter, r *http.Request) {
	id := tenant.FromContext(r.Context())
	rows, err := h.Store.Query(r.Context(), id, "SELECT country FROM catalog.vat_countries WHERE tenant_id=$1 ORDER BY country", id)
	if err != nil {
		writeError(w, err)
		return
	}
	codes, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "private, no-store")
	httpx.WriteData(w, http.StatusOK, codes, nil)
}

func (h Countries) Change(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("country")
	if !countries.Valid(code) {
		httpx.WriteError(w, 422, "invalid_country", "Choose a country from the list.", nil)
		return
	}
	id := tenant.FromContext(r.Context())
	query := "INSERT INTO catalog.vat_countries (tenant_id,country) VALUES ($1,$2) ON CONFLICT DO NOTHING"
	if r.Method == http.MethodDelete {
		query = "DELETE FROM catalog.vat_countries WHERE tenant_id=$1 AND country=$2"
	}
	if _, err := h.Store.Exec(r.Context(), id, query, id, code); err != nil {
		writeError(w, err)
		return
	}
	httpx.WriteData(w, http.StatusOK, map[string]string{"country": code}, nil)
}

func writeError(w http.ResponseWriter, err error) {
	var pg *pgconn.PgError
	if errors.As(err, &pg) && pg.Code == "23503" {
		httpx.WriteError(w, 409, "country_in_use", "Remove this country's rates from products before removing the country.", nil)
		return
	}
	httpx.WriteError(w, 500, "vat_settings_failed", "Could not access the catalog VAT settings.", nil)
}
