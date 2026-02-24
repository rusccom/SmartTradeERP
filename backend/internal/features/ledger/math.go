package ledger

type calcState struct {
    qty float64
    avg float64
}

type calcResult struct {
    state  calcState
    cogs   *float64
    profit *float64
}

func applyIn(state calcState, qty, unitPrice float64) calcResult {
    totalQty := state.qty + qty
    if totalQty <= 0 {
        return calcResult{state: calcState{qty: 0, avg: 0}}
    }
    weighted := state.qty*state.avg + qty*unitPrice
    next := calcState{qty: totalQty, avg: weighted / totalQty}
    return calcResult{state: next}
}

func applyOut(state calcState, qty float64, revenue *float64) calcResult {
    cogsValue := qty * state.avg
    next := calcState{qty: state.qty - qty, avg: state.avg}
    result := calcResult{state: next, cogs: &cogsValue}
    if revenue == nil {
        return result
    }
    profitValue := *revenue - cogsValue
    result.profit = &profitValue
    return result
}
