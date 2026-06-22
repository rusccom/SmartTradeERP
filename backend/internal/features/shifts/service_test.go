package shifts

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"smarterp/backend/internal/shared/validation"
)

func dec(value string) decimal.Decimal {
	return decimal.RequireFromString(value)
}

const uuidWarehouse = "11111111-1111-1111-1111-111111111111"

func TestExpectedCash(t *testing.T) {
	totals := shiftTotals{
		salesCash:    dec("150"),
		returnsCash:  dec("20"),
		totalCashIn:  dec("50"),
		totalCashOut: dec("30"),
		salesCard:    dec("999"), // card must not affect the cash drawer
	}
	got := expectedCash(dec("100"), totals)
	// 100 + 150 - 20 + 50 - 30 = 250
	if !got.Equal(dec("250")) {
		t.Fatalf("expected cash = %s, want 250", got)
	}
}

func TestValidateOpenRequest(t *testing.T) {
	ok := OpenRequest{WarehouseID: uuidWarehouse, OpeningCash: dec("0")}
	if err := validateOpenRequest(ok); err != nil {
		t.Fatalf("valid open request rejected: %v", err)
	}
	if err := validateOpenRequest(OpenRequest{WarehouseID: "nope", OpeningCash: dec("0")}); !errors.Is(err, validation.ErrInvalidData) {
		t.Fatalf("bad warehouse must be rejected, got %v", err)
	}
	if err := validateOpenRequest(OpenRequest{WarehouseID: uuidWarehouse, OpeningCash: dec("-1")}); !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("negative opening cash must be rejected, got %v", err)
	}
}

func TestValidateCashOp(t *testing.T) {
	if err := validateCashOp(CashOpRequest{Type: "cash_in", Amount: dec("10")}); err != nil {
		t.Fatalf("valid cash_in rejected: %v", err)
	}
	if err := validateCashOp(CashOpRequest{Type: "cash_out", Amount: dec("10")}); err != nil {
		t.Fatalf("valid cash_out rejected: %v", err)
	}
	if err := validateCashOp(CashOpRequest{Type: "withdraw", Amount: dec("10")}); !errors.Is(err, ErrInvalidCashOpType) {
		t.Fatalf("unknown cash op type must be rejected, got %v", err)
	}
	if err := validateCashOp(CashOpRequest{Type: "cash_in", Amount: dec("0")}); !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("zero amount must be rejected, got %v", err)
	}
}
