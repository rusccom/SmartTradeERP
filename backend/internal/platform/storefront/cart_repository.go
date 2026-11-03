package storefront

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"smarterp/backend/internal/shared/db"
)

var errCartChanged = errors.New("cart changed")

type savedCart struct {
	Items   []cartItemInput `json:"items"`
	Version int64           `json:"version"`
}

func (sh *shop) readCart(id string) (savedCart, error) {
	const query = "SELECT items, version FROM storefront.carts WHERE tenant_id=$1 AND id=$2 AND preview=$3"
	var cart savedCart
	err := sh.service.repo.store.QueryRow(sh.ctx, sh.tenantID, query, sh.tenantID, id, sh.preview).Scan(&cart.Items, &cart.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return savedCart{}, nil
	}
	return cart, err
}

func (sh *shop) saveCart(tx db.DBTX, id string, cart savedCart) error {
	const create = "INSERT INTO storefront.carts(tenant_id,id,preview,items) VALUES($1,$2,$3,$4) ON CONFLICT DO NOTHING"
	const update = "UPDATE storefront.carts SET items=$4, version=version+1 WHERE tenant_id=$1 AND id=$2 AND preview=$3 AND version=$5"
	var tag pgconn.CommandTag
	var err error
	if cart.Version == 0 {
		tag, err = tx.Exec(sh.ctx, create, sh.tenantID, id, sh.preview, cart.Items)
	} else {
		tag, err = tx.Exec(sh.ctx, update, sh.tenantID, id, sh.preview, cart.Items, cart.Version)
	}
	if err == nil && tag.RowsAffected() == 0 {
		return errCartChanged
	}
	return err
}
