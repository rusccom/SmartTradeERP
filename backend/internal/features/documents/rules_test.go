package documents

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
)

func TestEnsureMutable(t *testing.T) {
	draft := Document{Type: "RECEIPT"}
	clean := CreateRequest{Type: "RECEIPT"}
	if err := ensureMutable(draft, clean); err != nil {
		t.Fatalf("plain draft edit should be allowed: %v", err)
	}

	shiftDoc := Document{Type: "SALE", ShiftID: uuidShift}
	if err := ensureMutable(shiftDoc, CreateRequest{Type: "SALE"}); !errors.Is(err, ErrShiftDocumentLocked) {
		t.Fatalf("shift document must be locked, got %v", err)
	}

	addShift := CreateRequest{Type: "RECEIPT", ShiftID: uuidShift}
	if err := ensureMutable(draft, addShift); !errors.Is(err, ErrShiftDocumentLocked) {
		t.Fatalf("binding a shift on edit must be blocked, got %v", err)
	}

	retype := CreateRequest{Type: "WRITEOFF"}
	if err := ensureMutable(draft, retype); !errors.Is(err, ErrTypeImmutable) {
		t.Fatalf("type change must be blocked, got %v", err)
	}
}

func TestParseGeneratedNumber(t *testing.T) {
	year, seq, ok := parseGeneratedNumber("SALE", "SALE-2026-000004")
	if !ok || year != 2026 || seq != 4 {
		t.Fatalf("got year=%d seq=%d ok=%v", year, seq, ok)
	}
	bad := []struct{ docType, number string }{
		{"SALE", "RECEIPT-2026-000004"}, // prefix mismatch
		{"SALE", "SALE-2026-4"},         // wrong sequence width
		{"SALE", "SALE-1999-000004"},    // year before 2000
		{"SALE", "SALE-2026"},           // wrong shape
	}
	for _, tc := range bad {
		if _, _, ok := parseGeneratedNumber(tc.docType, tc.number); ok {
			t.Fatalf("%q should not parse as a generated number", tc.number)
		}
	}
}

func TestRevenueForType(t *testing.T) {
	if got := revenueForType("SALE", dec("100")); got == nil || !got.Equal(dec("100")) {
		t.Fatalf("sale revenue should equal amount, got %v", got)
	}
	for _, docType := range []string{"RECEIPT", "WRITEOFF", "TRANSFER", "INVENTORY", "RETURN"} {
		if got := revenueForType(docType, dec("100")); got != nil {
			t.Fatalf("%s should carry no revenue, got %v", docType, got)
		}
	}
}

func TestDocMeta(t *testing.T) {
	cases := map[string][2]string{
		"RECEIPT":  {"IN", "PURCHASE"},
		"WRITEOFF": {"OUT", "WRITEOFF"},
		"SALE":     {"OUT", "SALE"},
	}
	for docType, want := range cases {
		gotType, gotReason := docMeta(docType)
		if gotType != want[0] || gotReason != want[1] {
			t.Fatalf("%s -> (%s,%s), want (%s,%s)", docType, gotType, gotReason, want[0], want[1])
		}
	}
}

func TestBuildShares_Proportional(t *testing.T) {
	shares := buildShares(map[string]decimal.Decimal{"a": dec("30"), "b": dec("10")}, dec("40"))
	if !shares["a"].Equal(dec("0.75")) || !shares["b"].Equal(dec("0.25")) {
		t.Fatalf("proportional shares wrong: %v", shares)
	}
}

func TestBuildShares_EqualWhenNoCost(t *testing.T) {
	shares := buildShares(map[string]decimal.Decimal{"a": dec("0"), "b": dec("0")}, dec("0"))
	if !shares["a"].Equal(dec("0.5")) || !shares["b"].Equal(dec("0.5")) {
		t.Fatalf("zero-cost bundle should split equally: %v", shares)
	}
}
