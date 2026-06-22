package ledger

import (
	"testing"

	"github.com/shopspring/decimal"
)

func dec(value string) decimal.Decimal {
	return decimal.RequireFromString(value)
}

func ptr(value string) *decimal.Decimal {
	d := dec(value)
	return &d
}

func receipt(qty, price, warehouse string) replayEvent {
	return replayEvent{id: warehouse + qty, direction: "IN", reason: "PURCHASE",
		qty: dec(qty), unitPrice: dec(price), warehouse: warehouse}
}

func sale(qty, warehouse string, revenue *decimal.Decimal) replayEvent {
	return replayEvent{id: "sale" + qty, direction: "OUT", reason: "SALE",
		qty: dec(qty), revenue: revenue, warehouse: warehouse}
}

func mustFold(t *testing.T, events []replayEvent, allowNegative bool) []replayResult {
	t.Helper()
	results, err := foldMovements(events, allowNegative)
	if err != nil {
		t.Fatalf("foldMovements failed: %v", err)
	}
	return results
}

func last(results []replayResult) replayResult {
	return results[len(results)-1]
}

func assertEqual(t *testing.T, got, want decimal.Decimal, label string) {
	t.Helper()
	if !got.Equal(want) {
		t.Fatalf("%s: got %s, want %s", label, got, want)
	}
}

func TestFold_WeightedAverage(t *testing.T) {
	events := []replayEvent{
		receipt("10", "8", "W1"),
		receipt("10", "12", "W1"),
		sale("5", "W1", ptr("100")),
	}
	results := mustFold(t, events, false)
	final := last(results)
	assertEqual(t, final.runningAvg, dec("10"), "avg after mixed receipts")
	assertEqual(t, final.runningQty, dec("15"), "qty after sale")
	assertEqual(t, *final.cogs, dec("50"), "cogs at average")
	assertEqual(t, *final.profit, dec("50"), "profit = revenue - cogs")
}

func TestFold_TransferIsCostNeutral(t *testing.T) {
	events := []replayEvent{
		receipt("10", "8", "W1"),
		{id: "out", direction: "OUT", reason: "TRANSFER_OUT", qty: dec("4"), warehouse: "W1"},
		{id: "in", direction: "IN", reason: "TRANSFER_IN", qty: dec("4"), unitPrice: dec("999"), warehouse: "W2"},
	}
	results := mustFold(t, events, false)
	final := last(results)
	assertEqual(t, final.runningAvg, dec("8"), "average unchanged by transfer")
	assertEqual(t, final.runningQty, dec("10"), "global qty unchanged by transfer")
	assertEqual(t, warehouseBalance(results, "W1"), dec("6"), "source warehouse")
	assertEqual(t, warehouseBalance(results, "W2"), dec("4"), "target warehouse")
}

func warehouseBalance(results []replayResult, warehouse string) decimal.Decimal {
	total := decimal.Zero
	for _, result := range results {
		if result.event.warehouse == warehouse {
			total = total.Add(result.qtyDelta)
		}
	}
	return total
}

func TestFold_InventoryCountSetsAbsoluteQty(t *testing.T) {
	events := []replayEvent{
		receipt("10", "8", "W1"),
		{id: "count", direction: "SET", reason: "COUNT", qty: dec("7"), warehouse: "W1"},
	}
	results := mustFold(t, events, false)
	final := last(results)
	assertEqual(t, final.qtyDelta, dec("-3"), "count delta = counted - on hand")
	assertEqual(t, final.runningQty, dec("7"), "running qty equals count")
	assertEqual(t, final.runningAvg, dec("8"), "count keeps average")
}

func TestFold_ReturnEntersAtCurrentAverage(t *testing.T) {
	events := []replayEvent{
		receipt("10", "8", "W1"),
		{id: "ret", direction: "IN", reason: "RETURN_IN", qty: dec("2"), unitPrice: dec("999"),
			revenue: ptr("-30"), warehouse: "W1"},
	}
	final := last(mustFold(t, events, false))
	assertEqual(t, final.unitCost, dec("8"), "return ignores recorded price, uses average")
	assertEqual(t, final.runningAvg, dec("8"), "return at average keeps average")
	assertEqual(t, *final.profit, dec("-14"), "return profit = revenue + cost = -30 + 16")
}

func TestFold_NegativeStockGuard(t *testing.T) {
	events := []replayEvent{
		receipt("5", "8", "W1"),
		sale("8", "W1", ptr("100")),
	}
	if _, err := foldMovements(events, false); err != ErrNegativeStock {
		t.Fatalf("expected ErrNegativeStock, got %v", err)
	}
	if _, err := foldMovements(events, true); err != nil {
		t.Fatalf("allowNegative should permit oversell, got %v", err)
	}
}

func TestFold_IsDeterministic(t *testing.T) {
	events := []replayEvent{
		receipt("10", "8", "W1"),
		receipt("7", "11", "W1"),
		sale("4", "W1", ptr("80")),
	}
	first := mustFold(t, events, false)
	second := mustFold(t, events, false)
	if len(first) != len(second) {
		t.Fatalf("length mismatch: %d vs %d", len(first), len(second))
	}
	for i := range first {
		assertEqual(t, first[i].runningQty, second[i].runningQty, "qty step")
		assertEqual(t, first[i].runningAvg, second[i].runningAvg, "avg step")
	}
}
