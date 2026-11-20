package orders

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"smarterp/backend/internal/shared/validation"
)

var (
	ErrInvalid   = errors.New("invalid order data")
	ErrConflict  = errors.New("the order changed; reopen it before saving")
	ErrIssued    = errors.New("this order already has an invoice")
	ErrCurrency  = errors.New("configure the organization's base currency")
	ErrSelection = errors.New("a selected catalog variant is unavailable or has no price")
	countryCode  = regexp.MustCompile(`^[A-Z]{2}$`)
)

func validAddress(a Address) bool {
	return (a.Country == "" || countryCode.MatchString(a.Country)) &&
		within(255, a.Name, a.Company, a.City, a.Region) && within(512, a.Street) && within(32, a.PostalCode, a.TaxID)
}

func within(limit int, values ...string) bool {
	for _, value := range values {
		if utf8.RuneCountInString(value) > limit || strings.ContainsRune(value, 0) {
			return false
		}
	}
	return true
}

func validDetails(d Details) bool {
	if !validAddress(d.Shipping) || !validAddress(d.Billing) || !within(160, d.Nickname, d.PaymentMethod, d.DeliveryMethod) {
		return false
	}
	if d.DeliveryKind == "address" {
		return d.Pickup == nil
	}
	if d.DeliveryKind != "pickup" || d.Pickup == nil {
		return false
	}
	p := d.Pickup
	return strings.TrimSpace(p.PointID) != "" && within(255, p.Name, p.PointID, p.City) && within(512, p.Street) && within(32, p.PostalCode)
}

func validAmount(v decimal.Decimal, digits int32, ceiling int64) bool {
	return v.Exponent() >= -digits && v.Exponent() <= 10 && !v.IsNegative() && v.LessThan(decimal.NewFromInt(ceiling))
}

func (e Entry) validate() error {
	if strings.TrimSpace(e.Name) == "" || !within(160, e.Name) || !within(40, e.Phone) || !within(255, e.Email) {
		return ErrInvalid
	}
	if !validation.Email(e.Email) || !validDetails(e.Details) {
		return ErrInvalid
	}
	if e.BillingCustomerID != "" && !validation.UUID(e.BillingCustomerID) {
		return ErrInvalid
	}
	return e.validateLines()
}

func (e Entry) validateLines() error {
	if len(e.Items) == 0 || len(e.Items) > 300 || !within(12, e.Details.DeliveryVAT) {
		return ErrInvalid
	}
	if e.DeliveryCost != nil && !validAmount(*e.DeliveryCost, 4, 10000000000) {
		return ErrInvalid
	}
	seen := make(map[string]bool)
	for _, line := range e.Items {
		if !validation.UUID(line.VariantID) || !validAmount(line.Qty, 3, 100001) || !line.Qty.IsPositive() || !within(12, line.VAT) {
			return ErrInvalid
		}
		if line.ID != "" && (!validation.UUID(line.ID) || seen[line.ID]) {
			return ErrInvalid
		}
		seen[line.ID] = true
	}
	return nil
}
