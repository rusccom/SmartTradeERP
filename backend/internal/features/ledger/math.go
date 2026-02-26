package ledger

import "github.com/shopspring/decimal"

type calcState struct {
	qty decimal.Decimal
	avg decimal.Decimal
}

type calcResult struct {
	state  calcState
	cogs   *decimal.Decimal
	profit *decimal.Decimal
}

func zeroState() calcState {
	return calcState{qty: decimal.Zero, avg: decimal.Zero}
}

func applyIn(state calcState, qty, unitPrice decimal.Decimal, revenue *decimal.Decimal) calcResult {
	totalQty := state.qty.Add(qty)
	if totalQty.LessThanOrEqual(decimal.Zero) {
		return calcResult{state: zeroState()}
	}
	weighted := state.qty.Mul(state.avg).Add(qty.Mul(unitPrice))
	newAvg := weighted.Div(totalQty).Round(4)
	result := calcResult{state: calcState{qty: totalQty, avg: newAvg}}
	if revenue != nil {
		value := revenue.Round(4)
		result.profit = &value
	}
	return result
}

func applyOut(state calcState, qty decimal.Decimal, revenue *decimal.Decimal) calcResult {
	cogsValue := qty.Mul(state.avg).Round(4)
	next := calcState{qty: state.qty.Sub(qty), avg: state.avg}
	result := calcResult{state: next, cogs: &cogsValue}
	if revenue == nil {
		return result
	}
	profitValue := revenue.Sub(cogsValue)
	result.profit = &profitValue
	return result
}
