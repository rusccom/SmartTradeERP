package orders

import (
	"smarterp/backend/internal/erp/productvat"
	"smarterp/backend/internal/shared/countries"
)

func (d Details) Destination() string {
	if d.DeliveryKind == "pickup" {
		if d.Pickup == nil {
			return ""
		}
		return d.Pickup.Country
	}
	return d.Shipping.Country
}

func (line LineInput) validVAT() bool {
	validMode := line.VATMode == "" || line.VATMode == "catalog" || line.VATMode == "manual"
	return validMode && (line.VAT == "" || productvat.ValidRate(line.VAT)) &&
		(line.VATCountry == "" || countries.Valid(line.VATCountry))
}

func (w operation) applyVAT(details Details, lines []Line) error {
	country := details.Destination()
	if country != "" && !countries.Valid(country) {
		return ErrInvalid
	}
	pending, ids := w.catalogPositions(lines, country)
	rates, err := productvat.ForVariants(w.ctx, w.tx, w.tenant, country, ids)
	if err != nil {
		return err
	}
	for _, line := range pending {
		line.VAT = rates[line.VariantID]
	}
	return nil
}

func (w operation) catalogPositions(lines []Line, country string) ([]*Line, []string) {
	previous := w.savedVAT()
	var pending []*Line
	var ids []string
	for i := range lines {
		line := &lines[i]
		if needsCatalogVAT(line, previous[line.ID], country) {
			pending = append(pending, line)
			ids = append(ids, line.VariantID)
		}
	}
	return pending, ids
}

func (w operation) savedVAT() map[string]LineInput {
	previous := map[string]LineInput{}
	if w.current != nil {
		for _, line := range w.current.Items {
			previous[line.ID] = line.LineInput
		}
	}
	return previous
}

func needsCatalogVAT(line *Line, prior LineInput, country string) bool {
	if line.VATMode != "catalog" && line.VAT != "" && (line.VATCountry == "" || line.VATCountry == country) {
		line.VATMode, line.VATCountry = "manual", country
		return false
	}
	if prior.ID != "" && prior.VATMode == "catalog" && prior.VATCountry == country {
		line.VAT, line.VATCountry, line.VATMode = prior.VAT, country, "catalog"
		return false
	}
	line.VAT, line.VATCountry, line.VATMode = "", country, "catalog"
	return true
}
