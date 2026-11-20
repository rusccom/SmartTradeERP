package orders

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/shared/validation"
)

func (w *operation) prepareEdit(revision int64) error {
	if err := Lock(w.ctx, w.tx, w.tenant, w.id); err != nil {
		return err
	}
	c, err := ReadCard(w.ctx, w.tx, w.tenant, w.id)
	if err != nil {
		return err
	}
	if c.Invoice != nil {
		return ErrIssued
	}
	if c.Revision != revision {
		return ErrConflict
	}
	w.current = c
	return nil
}

func (w *operation) save(e Entry, create bool) (*Card, error) {
	if err := e.validate(); err != nil {
		return nil, err
	}
	if create {
		duplicate, err := w.beginOrder(e)
		if err != nil || duplicate {
			return w.current, err
		}
	} else if err := w.prepareEdit(e.Revision); err != nil {
		return nil, err
	}
	return w.persist(e, create)
}

func (w *operation) persist(e Entry, create bool) (*Card, error) {
	q, err := w.quote(e)
	if err != nil {
		return nil, err
	}
	if err = w.writeHeader(e, q, create); err != nil {
		return nil, err
	}
	if err = w.writeLines(q.Items); err != nil {
		return nil, err
	}
	return ReadCard(w.ctx, w.tx, w.tenant, w.id)
}

func requestHash(e Entry) []byte {
	data, _ := json.Marshal(e)
	hash := sha256.Sum256(data)
	return hash[:]
}

func (w *operation) beginOrder(e Entry) (bool, error) {
	if !validation.UUID(e.RequestID) {
		return false, ErrInvalid
	}
	w.id = e.RequestID
	if _, err := w.tx.Exec(w.ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1,0))`, w.tenant+w.id); err != nil {
		return false, err
	}
	var hash []byte
	err := w.tx.QueryRow(w.ctx, `SELECT entry_hash FROM storefront.orders WHERE tenant_id=$1 AND id=$2`, w.tenant, w.id).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !bytes.Equal(hash, requestHash(e)) {
		return false, ErrConflict
	}
	w.current, err = ReadCard(w.ctx, w.tx, w.tenant, w.id)
	return true, err
}

func NextNumber(w NumberSeries) (string, error) {
	const q = `INSERT INTO sales.sequences AS seq (tenant_id,series,year,last_number) VALUES ($1,$2,$3,1)
        ON CONFLICT (tenant_id,series,year) DO UPDATE SET last_number=seq.last_number+1 RETURNING last_number`
	var next int64
	err := w.Tx.QueryRow(w.Ctx, q, w.Tenant, w.Prefix, w.Year).Scan(&next)
	return fmt.Sprintf("%s/%d/%06d", w.Prefix, w.Year, next), err
}

func (w operation) writeHeader(e Entry, q *Quote, create bool) error {
	if create {
		if err := w.insertHeader(e, q.Currency); err != nil {
			return err
		}
	}
	const update = `UPDATE storefront.orders SET customer_name=$3,customer_email=$4,customer_phone=$5,shipping_address=$6,
      invoice_requested=$7,delivery_cost=$8,details=$9,billing_customer_id=NULLIF($10,'')::uuid,total_amount=$11,
      revision=revision+$12 WHERE tenant_id=$1 AND id=$2`
	increment := 1
	if create {
		increment = 0
	}
	_, err := w.tx.Exec(w.ctx, update, w.tenant, w.id, e.Name, e.Email, e.Phone, e.Details.Shipping.Street,
		e.InvoiceRequested, e.DeliveryCost, e.Details, e.BillingCustomerID, q.Total, increment)
	return err
}

func (w operation) insertHeader(e Entry, currency string) error {
	num, err := NextNumber(NumberSeries{w.ctx, w.tx, w.tenant, "ORD", time.Now().UTC().Year()})
	if err != nil {
		return err
	}
	const insert = `INSERT INTO storefront.orders(tenant_id,id,number,source,entry_hash,customer_name,customer_email,currency_code,created_at)
	VALUES ($1,$2,$3,'manual',$4,$5,$6,$7,timezone('UTC',now()))`
	_, err = w.tx.Exec(w.ctx, insert, w.tenant, w.id, num, requestHash(e), e.Name, e.Email, currency)
	return err
}

func (w operation) writeLines(lines []Line) error {
	ids := make([]string, 0, len(lines))
	for position, line := range lines {
		if line.ID == "" {
			line.ID = uuid.NewString()
		}
		if err := w.writeLine(line, position); err != nil {
			return err
		}
		ids = append(ids, line.ID)
	}
	_, err := w.tx.Exec(w.ctx, `DELETE FROM storefront.order_items WHERE tenant_id=$1 AND order_id=$2 AND NOT(id::text=ANY($3::text[]))`, w.tenant, w.id, ids)
	return err
}

func (w operation) writeLine(l Line, position int) error {
	const q = `INSERT INTO storefront.order_items (tenant_id,order_id,id,variant_id,product_name,variant_name,unit_price,qty,total_amount,vat_rate,sku_code,unit,position)
      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NULLIF($10,''),$11,$12,$13)
      ON CONFLICT(id) DO UPDATE SET qty=excluded.qty,total_amount=excluded.total_amount,vat_rate=excluded.vat_rate,position=excluded.position
      WHERE storefront.order_items.tenant_id=excluded.tenant_id AND storefront.order_items.order_id=excluded.order_id`
	result, err := w.tx.Exec(w.ctx, q, w.tenant, w.id, l.ID, l.VariantID, l.ProductName, l.VariantName, l.Price, l.Qty, l.Total, l.VAT, l.SKU, l.Unit, position)
	if err == nil && result.RowsAffected() != 1 {
		return ErrConflict
	}
	return err
}
