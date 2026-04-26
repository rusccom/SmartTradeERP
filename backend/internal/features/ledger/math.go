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
	newAvg := nextAverage(state, qty, unitPrice, totalQty)
	cost := qty.Mul(unitPrice).Round(4)
	result := calcResult{state: calcState{qty: totalQty, avg: newAvg}}
	result.profit = inProfit(revenue, cost)
	return result
}

func inProfit(revenue *decimal.Decimal, cost decimal.Decimal) *decimal.Decimal {
	if revenue == nil {
		return nil
	}
	value := revenue.Add(cost).Round(4)
	return &value
}

func nextAverage(state calcState, qty, unitPrice, totalQty decimal.Decimal) decimal.Decimal {
	if state.qty.LessThanOrEqual(decimal.Zero) {
		return unitPrice.Round(4)
	}
	if totalQty.LessThanOrEqual(decimal.Zero) {
		return state.avg
	}
	weighted := state.qty.Mul(state.avg).Add(qty.Mul(unitPrice))
	return weighted.Div(totalQty).Round(4)
}

func applyOut(state calcState, qty decimal.Decimal, revenue *decimal.Decimal) calcResult {
	cogsValue := qty.Mul(state.avg).Round(4)
	next := calcState{qty: state.qty.Sub(qty), avg: state.avg}
	result := calcResult{state: next, cogs: &cogsValue}
	if revenue == nil {
		return result
	}
	profitValue := revenue.Sub(cogsValue).Round(4)
	result.profit = &profitValue
	return result
}
