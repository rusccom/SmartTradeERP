package orders

import (
	"net/http"

	"github.com/shopspring/decimal"
	"smarterp/backend/internal/shared/validation"
)

type receipt struct {
	RequestID string          `json:"request_id"`
	Amount    decimal.Decimal `json:"amount"`
}

func (h Handler) Payment(w http.ResponseWriter, r *http.Request) {
	var p receipt
	if !Decode(w, r, &p) {
		return
	}
	if !validation.UUID(p.RequestID) || !validAmount(p.Amount, 4, 10000000000) || !p.Amount.IsPositive() {
		WriteError(w, ErrInvalid)
		return
	}
	h.transact(w, r, func(op *operation) (any, error) { return op.record(p) })
}

func (w operation) record(p receipt) (*Card, error) {
	if err := Lock(w.ctx, w.tx, w.tenant, w.id); err != nil {
		return nil, err
	}
	card, err := ReadCard(w.ctx, w.tx, w.tenant, w.id)
	if err != nil {
		return nil, err
	}
	w.current = card
	m, err := w.currency()
	if err != nil {
		return nil, err
	}
	if !p.Amount.Equal(p.Amount.Round(m.places)) {
		return nil, ErrInvalid
	}
	return w.storeReceipt(p)
}

func (w operation) storeReceipt(p receipt) (*Card, error) {
	const q = `INSERT INTO sales.order_payments AS saved(tenant_id,id,order_id,amount) VALUES ($1,$2,$3,$4)
      ON CONFLICT (tenant_id,id) DO UPDATE SET id=saved.id WHERE saved.order_id=excluded.order_id AND saved.amount=excluded.amount`
	result, err := w.tx.Exec(w.ctx, q, w.tenant, p.RequestID, w.id, p.Amount)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() != 1 {
		return nil, ErrConflict
	}
	return ReadCard(w.ctx, w.tx, w.tenant, w.id)
}
