package pl

import (
	"encoding/xml"
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"strings"
	"time"
	"unicode/utf8"

	"smarterp/backend/internal/erp/orders"
)

func renderFA(r invoiceRecord) ([]byte, error) {
	o := r.Order
	if err := validateDomesticInvoice(o); err != nil {
		return nil, err
	}
	buyer, err := buyerParty(o.Details.Billing)
	if err != nil {
		return nil, err
	}
	body, err := invoiceBody(r)
	if err != nil {
		return nil, err
	}
	header := faHeader{formCode{"FA (3)", "1-0E", "FA"}, 3, r.Invoice.IssuedAt.Format(time.RFC3339), "SmartTrade ERP"}
	seller := faParty{Identity: identity{NIP: r.Seller.NIP, Name: r.Seller.Name}, Address: faAddress{"PL", r.Seller.Street, r.Seller.Postal + " " + r.Seller.City}}
	raw, err := xml.Marshal(faDocument{Header: header, Seller: seller, Buyer: buyer, Body: body})
	return append([]byte(xml.Header), raw...), err
}

func buyerParty(a orders.Address) (faParty, error) {
	name := strings.TrimSpace(a.Company)
	if name == "" {
		name = strings.TrimSpace(a.Name)
	}
	country, err := language.ParseRegion(a.Country)
	if name == "" || strings.TrimSpace(a.Street) == "" || strings.TrimSpace(a.City) == "" || strings.TrimSpace(a.PostalCode) == "" || err != nil || !country.IsCountry() {
		return faParty{}, problem("Complete the billing name, street, postal code, city and country.")
	}
	id, err := buyerIdentity(a, name)
	if err != nil {
		return faParty{}, err
	}
	return faParty{Identity: id, Address: faAddress{a.Country, a.Street, a.PostalCode + " " + a.City}, JST: 2, GV: 2}, nil
}

func buyerIdentity(a orders.Address, name string) (identity, error) {
	id := identity{Name: name}
	switch {
	case a.TaxID == "":
		id.NoID = 1
	case a.Country == "PL":
		if !validNIP(a.TaxID) {
			return identity{}, problem("Check the Polish buyer NIP.")
		}
		id.NIP = a.TaxID
	default:
		id.Country = a.Country
		id.ForeignID = a.TaxID
	}
	return id, nil
}

func invoiceBody(r invoiceRecord) (faBody, error) {
	o := r.Order
	b := faBody{Currency: o.Currency, Date: r.Invoice.IssuedAt.In(warsaw).Format("2006-01-02"), Number: r.Invoice.Number,
		Total: o.Grand.StringFixed(2), Notes: annotation{2, 2, 2, 2, 1, 1, 2, 1}, Kind: "VAT"}
	err := b.addLines(o)
	return b, err
}

func (b *faBody) addLines(o *orders.Card) error {
	sums := make(map[string]taxGroup)
	var gross decimal.Decimal
	for _, l := range invoiceLines(o) {
		if !l.Total.Equal(l.Total.Round(2)) {
			return problem("This invoice format requires each gross line total to be stated in grosze.")
		}
		gross = gross.Add(l.Total)
		row, err := taxLine(l, sums)
		if err != nil {
			return err
		}
		row.Position = len(b.Rows) + 1
		b.Rows = append(b.Rows, row)
	}
	if !gross.Equal(*o.Grand) {
		return problem("Invoice line amounts must match the saved order total.")
	}
	b.Sums = taxFields(sums)
	return nil
}

func invoiceLines(o *orders.Card) []orders.Line {
	if !o.DeliveryCost.IsPositive() {
		return o.Items
	}
	ship := orders.Line{ProductName: "Dostawa", Unit: "usł.", Price: *o.DeliveryCost, Total: *o.DeliveryCost}
	ship.VAT = o.Details.DeliveryVAT
	ship.Qty = one
	return append(append([]orders.Line{}, o.Items...), ship)
}

func lineName(l orders.Line) (string, error) {
	name := l.ProductName
	if l.VariantName != "" && l.VariantName != l.ProductName {
		name += " / " + l.VariantName
	}
	if name == "" || utf8.RuneCountInString(name) > 512 {
		return "", problem("The item name does not fit the invoice format.")
	}
	return name, nil
}
