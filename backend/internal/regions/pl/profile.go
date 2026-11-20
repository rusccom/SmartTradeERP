package pl

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/shared/db"
)

type Seller struct {
	Name        string `json:"name"`
	NIP         string `json:"nip"`
	Street      string `json:"street"`
	Postal      string `json:"postal_code"`
	City        string `json:"city"`
	Environment string `json:"environment"`
	HasToken    bool   `json:"token_configured"`
	Revision    int64  `json:"revision"`
}

type profileUpdate struct {
	Seller
	Token string `json:"token"`
	Clear bool   `json:"clear_token"`
}

func validNIP(n string) bool {
	if len(n) != 10 || n == "0000000000" {
		return false
	}
	weights := []int{6, 5, 7, 2, 3, 4, 5, 6, 7}
	sum := 0
	for i, c := range []byte(n) {
		if c < '0' || c > '9' {
			return false
		}
		if i < 9 {
			sum += int(c-'0') * weights[i]
		}
	}
	return sum%11 == int(n[9]-'0')
}

func (s Seller) valid() bool {
	if !validNIP(s.NIP) || s.Environment != "test" && s.Environment != "production" {
		return false
	}
	for _, v := range []string{s.Name, s.Street, s.Postal, s.City} {
		if strings.TrimSpace(v) == "" || utf8.RuneCountInString(v) > 256 {
			return false
		}
	}
	return utf8.RuneCountInString(s.Postal+" "+s.City) <= 512
}

func readSeller(ctx context.Context, tx db.DBTX, tenant string) (*Seller, error) {
	const q = `SELECT name,nip,street,postal_code,city,environment,encrypted_token IS NOT NULL,revision FROM region_pl.sellers WHERE tenant_id=$1`
	var s Seller
	err := tx.QueryRow(ctx, q, tenant).Scan(&s.Name, &s.NIP, &s.Street, &s.Postal, &s.City, &s.Environment, &s.HasToken, &s.Revision)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &s, err
}

func (h *Handler) saveProfile(ctx context.Context, tx db.DBTX, tenant string, u profileUpdate) (*Seller, error) {
	if !u.valid() || len(u.Token) > 8192 || u.Clear && u.Token != "" {
		return nil, problem("Check the seller name, NIP, address and environment.")
	}
	sealed, err := h.seal(tenant, u)
	if err != nil {
		return nil, err
	}
	const q = `INSERT INTO region_pl.sellers AS s(tenant_id,name,nip,street,postal_code,city,environment,encrypted_token)
      SELECT $1,$2,$3,$4,$5,$6,$7,$8 WHERE $9::bigint=0 OR EXISTS (SELECT 1 FROM region_pl.sellers WHERE tenant_id=$1 AND revision=$9)
      ON CONFLICT (tenant_id) DO UPDATE SET name=excluded.name,nip=excluded.nip,street=excluded.street,postal_code=excluded.postal_code,
      city=excluded.city,environment=excluded.environment,encrypted_token=CASE WHEN $10 OR s.nip<>excluded.nip OR s.environment<>excluded.environment
      THEN excluded.encrypted_token ELSE COALESCE(excluded.encrypted_token,s.encrypted_token) END,revision=s.revision+1 WHERE s.revision=$9`
	tag, err := tx.Exec(ctx, q, tenant, u.Name, u.NIP, u.Street, u.Postal, u.City, u.Environment, sealed, u.Revision, u.Clear)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() != 1 {
		return nil, problem("Seller settings changed. Reload before saving.")
	}
	return readSeller(ctx, tx, tenant)
}
