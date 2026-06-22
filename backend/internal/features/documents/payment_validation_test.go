package documents

import (
	"errors"
	"testing"
)

func saleItems() []ItemInput {
	return []ItemInput{{VariantID: uuidVariant, Qty: dec("2"), UnitPrice: dec("15")}}
}

func TestValidatePayments_SaleMatches(t *testing.T) {
	payments := []PaymentInput{
		{Method: "cash", Amount: dec("20")},
		{Method: "card", Amount: dec("10")},
	}
	if err := validatePayments("SALE", saleItems(), payments); err != nil {
		t.Fatalf("mixed payment matching total rejected: %v", err)
	}
}

func TestValidatePayments_Errors(t *testing.T) {
	cases := []struct {
		name     string
		docType  string
		payments []PaymentInput
		want     error
	}{
		{"sale requires payments", "SALE", nil, ErrPaymentsRequired},
		{"receipt forbids payments", "RECEIPT", []PaymentInput{{Method: "cash", Amount: dec("30")}}, ErrPaymentsNotAllowed},
		{"bad method", "SALE", []PaymentInput{{Method: "crypto", Amount: dec("30")}}, ErrInvalidPaymentMethod},
		{"zero amount", "SALE", []PaymentInput{{Method: "cash", Amount: dec("0")}}, ErrInvalidPaymentAmount},
		{"total mismatch", "SALE", []PaymentInput{{Method: "cash", Amount: dec("29")}}, ErrPaymentTotalMismatch},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePayments(tc.docType, saleItems(), tc.payments)
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestValidatePayments_ReceiptWithoutPaymentsOK(t *testing.T) {
	if err := validatePayments("RECEIPT", saleItems(), nil); err != nil {
		t.Fatalf("receipt without payments should pass: %v", err)
	}
}

func TestSumPaymentsTotal(t *testing.T) {
	payments := []PaymentInput{
		{Method: "cash", Amount: dec("10.5")},
		{Method: "card", Amount: dec("4.25")},
	}
	if got := sumPaymentsTotal(payments); !got.Equal(dec("14.75")) {
		t.Fatalf("payments total = %s, want 14.75", got)
	}
}
