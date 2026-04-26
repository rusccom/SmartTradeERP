package documents

import "github.com/shopspring/decimal"

func validatePayments(documentType string, items []ItemInput, payments []PaymentInput) error {
	if !allowsPayments(documentType) {
		return validateNoPayments(payments)
	}
	if len(payments) == 0 {
		return ErrPaymentsRequired
	}
	if err := validatePaymentRows(payments); err != nil {
		return err
	}
	return validatePaymentTotal(items, payments)
}

func validateNoPayments(payments []PaymentInput) error {
	if len(payments) > 0 {
		return ErrPaymentsNotAllowed
	}
	return nil
}

func validatePaymentTotal(items []ItemInput, payments []PaymentInput) error {
	itemsTotal := sumItemsTotal(items)
	paymentsTotal := sumPaymentsTotal(payments)
	if !itemsTotal.Equal(paymentsTotal) {
		return ErrPaymentTotalMismatch
	}
	return nil
}

func validatePaymentRows(payments []PaymentInput) error {
	for _, payment := range payments {
		if !isAllowedPaymentMethod(payment.Method) {
			return ErrInvalidPaymentMethod
		}
		if invalidPaymentAmount(payment) {
			return ErrInvalidPaymentAmount
		}
	}
	return nil
}

func invalidPaymentAmount(payment PaymentInput) bool {
	return !payment.Amount.GreaterThan(decimal.Zero)
}

func allowsPayments(documentType string) bool {
	return documentType == "SALE" || documentType == "RETURN"
}

func isAllowedPaymentMethod(method string) bool {
	return method == "cash" || method == "card" || method == "transfer"
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
