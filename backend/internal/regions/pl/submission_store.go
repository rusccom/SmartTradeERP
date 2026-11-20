package pl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"smarterp/backend/internal/erp/orders"
)

type delivery struct {
	ID          string    `json:"-"`
	State       string    `json:"state"`
	Environment string    `json:"environment,omitempty"`
	Number      string    `json:"ksef_number,omitempty"`
	Message     string    `json:"message,omitempty"`
	Updated     time.Time `json:"updated_at"`
	Attempt     string    `json:"-"`
	Session     string    `json:"-"`
	Reference   string    `json:"-"`
	Payload     []byte    `json:"-"`
	SellerNIP   string    `json:"-"`
}

func loadDelivery(s scope, id string, lock bool) (*delivery, error) {
	if lock {
		var key string
		err := s.tx.QueryRow(s.ctx, `SELECT invoice_id FROM region_pl.invoice_payloads WHERE tenant_id=$1 AND invoice_id=$2 FOR UPDATE`, s.tenant, id).Scan(&key)
		if err != nil {
			return nil, err
		}
	}
	const q = `SELECT i.id::text,coalesce(k.state,'ready'),p.environment,coalesce(k.ksef_number,''),coalesce(k.message,''),
      coalesce(k.updated_at,i.issued_at),coalesce(k.attempt_id::text,''),coalesce(k.session_ref,''),coalesce(k.invoice_ref,''),x.payload,i.snapshot#>>'{seller,nip}'
      FROM region_pl.invoice_payloads x JOIN sales.invoices i ON (i.tenant_id,i.id)=(x.tenant_id,x.invoice_id)
      JOIN region_pl.sellers p ON p.tenant_id=i.tenant_id
      LEFT JOIN region_pl.invoice_submissions k ON (k.tenant_id,k.invoice_id,k.environment)=(x.tenant_id,x.invoice_id,p.environment)
      WHERE x.tenant_id=$1 AND x.invoice_id=$2`
	var d delivery
	err := s.tx.QueryRow(s.ctx, q, s.tenant, id).Scan(&d.ID, &d.State, &d.Environment, &d.Number, &d.Message, &d.Updated,
		&d.Attempt, &d.Session, &d.Reference, &d.Payload, &d.SellerNIP)
	return &d, err
}

func (d *delivery) claim(s scope, profile Seller) error {
	if d.State == "accepted" || d.State == "pending" || d.State == "uncertain" {
		return problem("Check the existing submission before sending again.")
	}
	if d.State == "sending" && (d.Session != "" || time.Since(d.Updated) < 3*time.Minute) {
		return problem("A submission is already in progress.")
	}
	if profile.NIP != d.SellerNIP || !profile.HasToken {
		return problem("Configure KSeF access for the seller on this invoice.")
	}
	d.Attempt = uuid.NewString()
	d.Environment = profile.Environment
	d.State = "sending"
	d.Session, d.Reference, d.Message = "", "", ""
	_, err := s.tx.Exec(s.ctx, `INSERT INTO region_pl.invoice_submissions(tenant_id,invoice_id,attempt_id,environment,state)
      VALUES ($1,$2,$3,$4,'sending') ON CONFLICT(tenant_id,invoice_id,environment) DO UPDATE SET
      state='sending',attempt_id=excluded.attempt_id,session_ref=NULL,invoice_ref=NULL,message=NULL,updated_at=now()`, s.tenant, d.ID, d.Attempt, d.Environment)
	return err
}

func (h *Handler) storeDelivery(ctx context.Context, d *delivery) error {
	return h.run(ctx, func(s scope) error {
		const q = `UPDATE region_pl.invoice_submissions SET state=$4,session_ref=NULLIF($5,''),invoice_ref=NULLIF($6,''),
          ksef_number=NULLIF($7,''),message=NULLIF($8,''),updated_at=now()
          WHERE tenant_id=$1 AND invoice_id=$2 AND attempt_id=$3 AND (state NOT IN ('accepted','rejected') OR state=$4)
          AND (session_ref IS NULL OR NULLIF($5,'') IS NOT NULL)`
		result, err := s.tx.Exec(s.ctx, q, s.tenant, d.ID, d.Attempt, d.State, d.Session, d.Reference, d.Number, d.Message)
		if err == nil && result.RowsAffected() != 1 {
			return orders.ErrConflict
		}
		return err
	})
}
