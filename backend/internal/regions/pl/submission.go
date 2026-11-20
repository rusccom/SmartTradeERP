package pl

import (
	"context"
	"net/http"
	"time"

	"smarterp/backend/internal/erp/orders"
	"smarterp/backend/internal/regions/pl/ksef"
	"smarterp/backend/internal/shared/validation"
)

func (h *Handler) Delivery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !validation.UUID(id) {
		respond(w, nil, orders.ErrInvalid)
		return
	}
	d, err := tenantResult(h, r.Context(), func(s scope) (*delivery, error) { return loadDelivery(s, id, false) })
	respond(w, d, err)
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	var consent struct {
		Confirmed   bool   `json:"confirmed"`
		Environment string `json:"environment"`
	}
	if !orders.Decode(w, r, &consent) {
		return
	}
	id := r.PathValue("id")
	if !consent.Confirmed || !validation.UUID(id) || consent.Environment != "test" && consent.Environment != "production" {
		respond(w, nil, orders.ErrInvalid)
		return
	}
	d, err := h.claim(r.Context(), id, consent.Environment)
	if err == nil {
		err = h.upload(r.Context(), d)
	}
	respond(w, d, err)
}

func (h *Handler) claim(ctx context.Context, id, environment string) (*delivery, error) {
	var d *delivery
	err := h.run(ctx, func(s scope) error {
		var err error
		d, err = loadDelivery(s, id, true)
		if err != nil {
			return err
		}
		return claimProfile(s, d, environment)
	})
	return d, err
}

func claimProfile(s scope, d *delivery, environment string) error {
	p, err := readSeller(s.ctx, s.tx, s.tenant)
	if err != nil {
		return err
	}
	if p == nil {
		return problem("Configure Polish seller settings.")
	}
	if p.Environment != environment {
		return problem("The KSeF environment changed. Reload and confirm the destination.")
	}
	if err := checkStoredVAT(s, d.ID); err != nil {
		return err
	}
	return d.claim(s, *p)
}

func (h *Handler) upload(parent context.Context, d *delivery) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(parent), 65*time.Second)
	defer cancel()
	access, err := h.session(ctx, d)
	if err != nil {
		return h.fail(parent, d, "ready", err)
	}
	up, err := access.Open(ctx, d.Payload)
	if err != nil {
		return h.fail(parent, d, "ready", err)
	}
	d.Session = up.Reference
	if err = h.storeDelivery(ctx, d); err != nil {
		return err
	}
	return h.finishUpload(ctx, access, up, d)
}

func (h *Handler) finishUpload(ctx context.Context, access *ksef.Session, up *ksef.Upload, d *delivery) error {
	var err error
	d.Reference, err = access.Send(ctx, up)
	if err != nil {
		_ = access.Close(ctx, d.Session)
		return h.fail(ctx, d, "uncertain", err)
	}
	d.State = "pending"
	if err = h.storeDelivery(ctx, d); err != nil {
		return err
	}
	_ = access.Close(ctx, d.Session)
	return nil
}

func (h *Handler) fail(parent context.Context, d *delivery, state string, cause error) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(parent), 5*time.Second)
	defer cancel()
	d.State, d.Message = state, cause.Error()
	return h.storeDelivery(ctx, d)
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !validation.UUID(id) {
		respond(w, nil, orders.ErrInvalid)
		return
	}
	d, err := tenantResult(h, r.Context(), func(s scope) (*delivery, error) { return loadDelivery(s, id, false) })
	if err == nil {
		err = h.reconcile(r.Context(), d)
	}
	respond(w, d, err)
}

func (h *Handler) reconcile(ctx context.Context, d *delivery) error {
	if d.State == "sending" && d.Session == "" && time.Since(d.Updated) > 3*time.Minute {
		return h.fail(ctx, d, "ready", problem("The previous attempt ended before the invoice was uploaded."))
	}
	if d.Session == "" || (d.State != "pending" && d.State != "uncertain" && d.State != "sending") {
		return nil
	}
	return h.refresh(ctx, d)
}

func (h *Handler) refresh(ctx context.Context, d *delivery) error {
	session, err := h.session(ctx, d)
	if err != nil {
		return problem(err.Error())
	}
	result, err := session.Check(ctx, d.Session, d.Reference, d.Payload)
	if err != nil {
		return problem(err.Error())
	}
	d.Reference = result.Reference
	d.State, d.Number = outcome(result)
	d.Message = result.Status.Message()
	return h.storeDelivery(ctx, d)
}

func outcome(r ksef.Result) (string, string) {
	if r.Status.Code == 200 && r.Number != "" {
		return "accepted", r.Number
	}
	if r.Status.Code >= 400 {
		return "rejected", ""
	}
	return "pending", ""
}
