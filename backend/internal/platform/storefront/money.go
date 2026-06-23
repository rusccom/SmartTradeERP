package storefront

import "github.com/shopspring/decimal"

// currency holds the tenant's base-currency display rules.
type currency struct {
	code   string
	symbol string
	places int32
}

// formatMoney renders a decimal amount for display and keeps the raw amount and
// ISO code for structured data (JSON-LD offers).
func formatMoney(value decimal.Decimal, cur currency) MoneyVM {
	places := cur.places
	if places < 0 || places > 4 {
		places = 2
	}
	amount := value.StringFixed(places)
	return MoneyVM{Display: displayPrice(amount, cur), Amount: amount, Code: cur.code}
}

func displayPrice(amount string, cur currency) string {
	if cur.symbol != "" {
		return cur.symbol + " " + amount
	}
	if cur.code != "" {
		return amount + " " + cur.code
	}
	return amount
}
