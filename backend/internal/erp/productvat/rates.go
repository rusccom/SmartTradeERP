package productvat

import (
	"context"
	"regexp"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/shared/countries"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/validation"
)

var percent = regexp.MustCompile(`^(100(\.0{1,2})?|[0-9]{1,2}(\.[0-9]{1,2})?)$`)

func ValidRate(value string) bool { return percent.MatchString(value) }

// RatesColumn is shared by catalog product queries, which alias the product as p.
const RatesColumn = `COALESCE((SELECT jsonb_object_agg(r.country,trim_scale(r.rate)::text)
    FROM catalog.product_vat_rates r WHERE (r.tenant_id,r.product_id)=(p.tenant_id,p.id)), '{}'::jsonb)`

// A missing key leaves that country's rate alone. A null value removes it.
type Changes map[string]*string

func (changes Changes) applyValid() bool {
	for country, rate := range changes {
		if !countries.Valid(country) || rate != nil && !ValidRate(*rate) {
			return false
		}
	}
	return true
}

func (changes Changes) Apply(ctx context.Context, tx db.DBTX, tenant, product string) error {
	if len(changes) == 0 {
		return nil
	}
	if !changes.applyValid() {
		return validation.ErrInvalidData
	}
	const query = `WITH changes AS (SELECT key AS country,value FROM jsonb_each_text($3::jsonb)),
	  removed AS (DELETE FROM catalog.product_vat_rates r USING changes c
	    WHERE (r.tenant_id,r.product_id,r.country)=($1,$2,c.country) AND c.value IS NULL)
	  INSERT INTO catalog.product_vat_rates(tenant_id,product_id,country,rate)
	    SELECT $1,$2,country,value::numeric FROM changes WHERE value IS NOT NULL
	    ORDER BY country ON CONFLICT (tenant_id,product_id,country) DO UPDATE SET rate=excluded.rate`
	_, err := tx.Exec(ctx, query, tenant, product, changes)
	return err
}

func ForVariants(ctx context.Context, tx db.DBTX, tenant, country string, variants []string) (map[string]string, error) {
	if country == "" || len(variants) == 0 {
		return nil, nil
	}
	const q = `SELECT v.id::text,trim_scale(r.rate)::text FROM catalog.product_variants v
	  JOIN catalog.product_vat_rates r ON (r.tenant_id,r.product_id)=(v.tenant_id,v.product_id)
	  WHERE v.tenant_id=$1 AND r.country=$2 AND v.id=ANY($3::uuid[])`
	rows, err := tx.Query(ctx, q, tenant, country, variants)
	if err != nil {
		return nil, err
	}
	pairs, err := pgx.CollectRows(rows, pgx.RowToStructByPos[struct{ Variant, Rate string }])
	rates := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		rates[pair.Variant] = pair.Rate
	}
	return rates, err
}
