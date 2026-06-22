package ledger

import "github.com/shopspring/decimal"

// replayEvent is one inventory movement as the cost engine sees it: a fact
// (direction, qty, purchase price for receipts) with no derived cost.
type replayEvent struct {
	id        string
	direction string
	reason    string
	qty       decimal.Decimal
	unitPrice decimal.Decimal
	revenue   *decimal.Decimal
	warehouse string
}

// replayResult is the cost outcome of one event after the running average is
// applied. movement_cost and inventory value are derived from these fields.
type replayResult struct {
	event      replayEvent
	seq        int64
	qtyDelta   decimal.Decimal
	unitCost   decimal.Decimal
	cogs       *decimal.Decimal
	profit     *decimal.Decimal
	runningQty decimal.Decimal
	runningAvg decimal.Decimal
}

// foldMovements replays an ordered event stream into per-event cost results.
// Pure: no database, deterministic for a given input. This is the single
// source of truth for inventory costing.
func foldMovements(events []replayEvent, allowNegative bool) ([]replayResult, error) {
	state := zeroState()
	warehouseQty := map[string]decimal.Decimal{}
	results := make([]replayResult, 0, len(events))
	for index, event := range events {
		result, next, err := foldOne(event, state, warehouseQty, int64(index+1), allowNegative)
		if err != nil {
			return nil, err
		}
		state = next
		results = append(results, result)
	}
	return results, nil
}

func foldOne(
	event replayEvent,
	state calcState,
	warehouseQty map[string]decimal.Decimal,
	seq int64,
	allowNegative bool,
) (replayResult, calcState, error) {
	delta := movementQtyDelta(event, warehouseQty[event.warehouse])
	price := effectiveUnitPrice(state, event)
	calc := applyDelta(state, delta, price, event.revenue)
	if err := checkStock(warehouseQty[event.warehouse], delta, calc.state.qty, allowNegative); err != nil {
		return replayResult{}, state, err
	}
	warehouseQty[event.warehouse] = warehouseQty[event.warehouse].Add(delta)
	return newReplayResult(event, calc, delta, price, seq), calc.state, nil
}

func newReplayResult(event replayEvent, calc calcResult, delta, price decimal.Decimal, seq int64) replayResult {
	return replayResult{
		event: event, seq: seq, qtyDelta: delta, unitCost: price,
		cogs: calc.cogs, profit: calc.profit,
		runningQty: calc.state.qty, runningAvg: calc.state.avg,
	}
}

func movementQtyDelta(event replayEvent, warehouseQty decimal.Decimal) decimal.Decimal {
	if event.direction == "IN" {
		return event.qty
	}
	if event.direction == "SET" {
		return event.qty.Sub(warehouseQty)
	}
	return event.qty.Neg()
}

// effectiveUnitPrice keeps the recorded purchase price only for real receipts.
// Transfers and returns re-enter at the average current at this point in the
// replay, so cross-warehouse moves stay cost-neutral under retro edits.
func effectiveUnitPrice(state calcState, event replayEvent) decimal.Decimal {
	if event.direction == "IN" && event.reason != "TRANSFER_IN" && event.reason != "RETURN_IN" {
		return event.unitPrice.Round(4)
	}
	return state.avg.Round(4)
}

func applyDelta(state calcState, delta, price decimal.Decimal, revenue *decimal.Decimal) calcResult {
	if delta.IsPositive() {
		return applyIn(state, delta, price, revenue)
	}
	if delta.IsNegative() {
		return applyOut(state, delta.Neg(), revenue)
	}
	return calcResult{state: state}
}

func checkStock(warehouseQty, delta, globalQty decimal.Decimal, allowNegative bool) error {
	if allowNegative {
		return nil
	}
	if warehouseQty.Add(delta).LessThan(decimal.Zero) {
		return ErrNegativeStock
	}
	if globalQty.LessThan(decimal.Zero) {
		return ErrNegativeStock
	}
	return nil
}

func movementCost(result replayResult) decimal.Decimal {
	if result.cogs != nil {
		return *result.cogs
	}
	return result.qtyDelta.Abs().Mul(result.unitCost).Round(4)
}

func inventoryValue(result replayResult) decimal.Decimal {
	return result.runningQty.Mul(result.runningAvg).Round(4)
}
