package pl

import "smarterp/backend/internal/erp/orders"

const domesticDelivery = problem("This Polish domestic VAT invoice requires a delivery address or pickup point in Poland.")

func validateDomesticInvoice(card *orders.Card) error {
	if card.Currency != "PLN" {
		return problem("The Polish standard VAT format currently supports PLN invoices.")
	}
	if card.Grand == nil {
		return problem("Enter the delivery cost before creating an invoice, including zero for free delivery.")
	}
	if card.Details.Destination() == "" {
		return domesticDelivery
	}
	return rejectForeignVAT(card)
}

func rejectForeignVAT(card *orders.Card) error {
	if country := card.Details.Destination(); country != "" && country != "PL" {
		return domesticDelivery
	}
	for _, line := range card.Items {
		if line.VATCountry != "" && line.VATCountry != "PL" {
			return problem("A foreign VAT rate cannot be issued using the Polish domestic VAT invoice format.")
		}
	}
	return nil
}

func checkStoredVAT(s scope, id string) error {
	const q = `SELECT snapshot->'order' FROM sales.invoices WHERE tenant_id=$1 AND id=$2`
	var order orders.Card
	if err := s.tx.QueryRow(s.ctx, q, s.tenant, id).Scan(&order); err != nil {
		return err
	}
	return rejectForeignVAT(&order)
}
