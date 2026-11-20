package pl

import (
	"encoding/xml"
	"github.com/shopspring/decimal"
	"smarterp/backend/internal/erp/orders"
	"unicode/utf8"
)

type taxGroup struct{ net, tax decimal.Decimal }

var one = decimal.NewFromInt(1)
var hundred = decimal.NewFromInt(100)
var standardRates = []struct{ rate, base, tax string }{
	{"23", "P_13_1", "P_14_1"}, {"8", "P_13_2", "P_14_2"}, {"5", "P_13_3", "P_14_3"}, {"0", "P_13_6_1", ""},
}

func taxLine(l orders.Line, sums map[string]taxGroup) (faLine, error) {
	name, err := lineName(l)
	if err != nil {
		return faLine{}, err
	}
	if utf8.RuneCountInString(l.SKU) > 50 {
		l.SKU = ""
	} // Indeks is optional and limited to 50 characters in FA(3).
	rate, err := decimal.NewFromString(l.VAT)
	if err != nil || rate.Exponent() < -4 || rate.Exponent() > 2 {
		return faLine{}, problem("Choose a standard VAT rate for this invoice.")
	}
	for _, band := range standardRates {
		if band.rate != rate.String() {
			continue
		}
		return groupedLine(l, name, band.rate, rate, sums), nil
	}
	return faLine{}, problem("Choose an explicit standard Polish VAT rate for every item and any delivery charge.")
}

func groupedLine(l orders.Line, name, code string, rate decimal.Decimal, sums map[string]taxGroup) faLine {
	vat := l.Total.Mul(rate).DivRound(hundred.Add(rate), 2)
	g := sums[code]
	sums[code] = taxGroup{g.net.Add(l.Total.Sub(vat)), g.tax.Add(vat)}
	if code == "0" {
		code = "0 KR"
	}
	return faLine{Name: name, SKU: l.SKU, Unit: l.Unit, Quantity: l.Qty.String(), Price: l.Price.String(),
		Gross: l.Total.StringFixed(2), Tax: vat.StringFixed(2), Rate: code}
}

func taxFields(sums map[string]taxGroup) []taxValue {
	fields := []taxValue{}
	for _, band := range standardRates {
		g, ok := sums[band.rate]
		if !ok {
			continue
		}
		fields = append(fields, taxValue{xml.Name{Local: band.base}, g.net.StringFixed(2)})
		if band.tax != "" {
			fields = append(fields, taxValue{xml.Name{Local: band.tax}, g.tax.StringFixed(2)})
		}
	}
	return fields
}
