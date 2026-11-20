package pl

import (
	"net/http"
	"time"
	_ "time/tzdata"

	"github.com/google/uuid"
	"smarterp/backend/internal/erp/orders"
	"smarterp/backend/internal/shared/validation"
)

const standardFormat = "pl-fa3-vat"

var warsaw, _ = time.LoadLocation("Europe/Warsaw")

type issueRequest struct {
	Format   string `json:"format"`
	Revision int64  `json:"revision"`
}
type invoiceRecord struct {
	Seller  Seller         `json:"seller"`
	Order   *orders.Card   `json:"order"`
	Invoice orders.Invoice `json:"invoice"`
}

func (h *Handler) Issue(w http.ResponseWriter, r *http.Request) {
	var req issueRequest
	if !orders.Decode(w, r, &req) {
		return
	}
	id := r.PathValue("id")
	if req.Format != standardFormat || !validation.UUID(id) {
		respond(w, nil, orders.ErrInvalid)
		return
	}
	invoice, err := tenantResult(h, r.Context(), func(s scope) (*orders.Invoice, error) { return issue(s, id, req) })
	respond(w, invoice, err)
}

func issue(s scope, id string, req issueRequest) (*orders.Invoice, error) {
	if err := orders.Lock(s.ctx, s.tx, s.tenant, id); err != nil {
		return nil, err
	}
	order, err := orders.ReadCard(s.ctx, s.tx, s.tenant, id)
	if err != nil {
		return nil, err
	}
	if order.Invoice != nil {
		return order.Invoice, nil
	}
	if order.Revision != req.Revision {
		return nil, orders.ErrConflict
	}
	if !order.InvoiceRequested {
		return nil, problem("Select the invoice request in Payment first.")
	}
	return issueForSeller(s, order, req.Format)
}

func issueForSeller(s scope, order *orders.Card, format string) (*orders.Invoice, error) {
	seller, err := readSeller(s.ctx, s.tx, s.tenant)
	if err != nil {
		return nil, err
	}
	if seller == nil {
		return nil, problem("Complete the Polish seller profile in order settings before issuing an invoice.")
	}
	now := time.Now().UTC()
	number, err := orders.NextNumber(orders.NumberSeries{Ctx: s.ctx, Tx: s.tx, Tenant: s.tenant, Prefix: "FV", Year: now.In(warsaw).Year()})
	if err != nil {
		return nil, err
	}
	doc := orders.Invoice{ID: uuid.NewString(), Number: number, Format: format, IssuedAt: now}
	return storeInvoice(s, invoiceRecord{*seller, order, doc})
}

func storeInvoice(s scope, record invoiceRecord) (*orders.Invoice, error) {
	payload, err := renderFA(record)
	if err != nil {
		return nil, err
	}
	i := record.Invoice
	const q = `INSERT INTO sales.invoices(tenant_id,id,order_id,number,format,issued_at,snapshot) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	if _, err = s.tx.Exec(s.ctx, q, s.tenant, i.ID, record.Order.ID, i.Number, i.Format, i.IssuedAt, record); err != nil {
		return nil, err
	}
	_, err = s.tx.Exec(s.ctx, `INSERT INTO region_pl.invoice_payloads(tenant_id,invoice_id,payload) VALUES ($1,$2,$3)`, s.tenant, i.ID, payload)
	return &i, err
}
