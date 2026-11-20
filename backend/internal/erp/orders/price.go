package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/shared/db"
)

type operation struct {
	ctx     context.Context
	tx      db.DBTX
	tenant  string
	id      string
	current *Card
}

type money struct {
	code   string
	places int32
}

func (w operation) currency() (money, error) {
	const base = `SELECT c.code,c.decimal_places FROM platform.tenant_currencies t JOIN platform.currencies c ON c.id=t.currency_id
        WHERE t.tenant_id=$1 AND t.is_base AND t.is_enabled`
	var m money
	row := w.tx.QueryRow
	var err error
	if w.current != nil {
		err = row(w.ctx, `SELECT code,decimal_places FROM platform.currencies WHERE code=$1`, w.current.Currency).Scan(&m.code, &m.places)
	} else {
		err = row(w.ctx, base, w.tenant).Scan(&m.code, &m.places)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return m, ErrCurrency
	}
	return m, err
}

func (w operation) quote(e Entry) (*Quote, error) {
	if err := e.validateLines(); err != nil {
		return nil, err
	}
	m, err := w.currency()
	if err != nil {
		return nil, err
	}
	if e.DeliveryCost != nil && !e.DeliveryCost.Equal(e.DeliveryCost.Round(m.places)) {
		return nil, ErrInvalid
	}
	q, err := w.priceCart(e, m)
	if err != nil {
		return nil, err
	}
	return q, w.applyVAT(e.Details, q.Items)
}

func (w operation) priceCart(e Entry, m money) (*Quote, error) {
	out := &Quote{Currency: m.code}
	for _, input := range e.Items {
		line, err := w.priceLine(input, m)
		if err != nil {
			return nil, err
		}
		out.Items = append(out.Items, line)
		out.Total = out.Total.Add(line.Total)
		out.Quantity = out.Quantity.Add(line.Qty)
	}
	if e.DeliveryCost != nil {
		g := out.Total.Add(*e.DeliveryCost)
		out.Grand = &g
	}
	if !validAmount(out.Total, 4, 10000000000) || out.Grand != nil && !validAmount(*out.Grand, 4, 10000000000) {
		return nil, ErrInvalid
	}
	return out, nil
}

func (w operation) priceLine(in LineInput, currency money) (Line, error) {
	if in.ID != "" {
		return w.existingLine(in, currency.places)
	}
	const q = `SELECT p.name,COALESCE(v.name,''),COALESCE(v.sku_code,''),COALESCE(v.unit,''),v.price
      FROM catalog.product_variants v JOIN catalog.products p ON p.tenant_id=v.tenant_id AND p.id=v.product_id
      WHERE v.tenant_id=$1 AND v.id=$2 AND v.price IS NOT NULL AND NOT p.is_composite
      AND EXISTS (SELECT 1 FROM platform.tenant_currencies t JOIN platform.currencies c ON c.id=t.currency_id
      WHERE t.tenant_id=v.tenant_id AND t.is_base AND t.is_enabled AND c.code=$3)`
	line := Line{LineInput: in}
	err := w.tx.QueryRow(w.ctx, q, w.tenant, in.VariantID, currency.code).Scan(&line.ProductName, &line.VariantName, &line.SKU, &line.Unit, &line.Price)
	if errors.Is(err, pgx.ErrNoRows) {
		return line, ErrSelection
	}
	line.Total = line.Price.Mul(in.Qty).Round(currency.places)
	return line, err
}

func (w operation) existingLine(in LineInput, places int32) (Line, error) {
	if w.current == nil {
		return Line{}, ErrInvalid
	}
	for _, line := range w.current.Items {
		if line.ID != in.ID || line.VariantID != in.VariantID {
			continue
		}
		if !line.Qty.Equal(in.Qty) {
			line.Total = line.Price.Mul(in.Qty).Round(places)
		}
		line.LineInput = in
		return line, nil
	}
	return Line{}, ErrInvalid
}
