package documents

import "github.com/shopspring/decimal"

func (s *Service) validateRequest(req CreateRequest) error {
	return validatePayments(req.Type, req.Items, req.Payments)
}

func validatePayments(documentType string, items []ItemInput, payments []PaymentInput) error {
	if len(payments) == 0 {
		return validateRequiredPayments(documentType)
	}
	if err := validatePaymentRows(payments); err != nil {
		return err
	}
	itemsTotal := sumItemsTotal(items)
	paymentsTotal := sumPaymentsTotal(payments)
	if !itemsTotal.Equal(paymentsTotal) {
		return ErrPaymentTotalMismatch
	}
	return nil
}

func validateRequiredPayments(documentType string) error {
	if documentType == "SALE" || documentType == "RETURN" {
		return ErrPaymentsRequired
	}
	return nil
}

func validatePaymentRows(payments []PaymentInput) error {
	for _, payment := range payments {
		if !isAllowedPaymentMethod(payment.Method) {
			return ErrInvalidPaymentMethod
		}
		if !payment.Amount.GreaterThan(decimal.Zero) {
			return ErrInvalidPaymentAmount
		}
	}
	return nil
}

func isAllowedPaymentMethod(method string) bool {
	switch method {
	case "cash", "card", "transfer":
		return true
	default:
		return false
	}
}

func sumItemsTotal(items []ItemInput) decimal.Decimal {
	total := decimal.Zero
	for _, item := range items {
		amount := item.Qty.Mul(item.UnitPrice).Round(4)
		total = total.Add(amount)
	}
	return total
}

func sumPaymentsTotal(payments []PaymentInput) decimal.Decimal {
	total := decimal.Zero
	for _, payment := range payments {
		total = total.Add(payment.Amount)
	}
	return total.Round(4)
}
