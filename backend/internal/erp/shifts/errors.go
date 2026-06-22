package shifts

import "errors"

var ErrShiftAlreadyOpen = errors.New("you already have an open shift")
var ErrNoOpenShift = errors.New("no open shift")
var ErrShiftAlreadyClosed = errors.New("shift is already closed")
var ErrInvalidCashOpType = errors.New("cash operation type must be cash_in or cash_out")
var ErrInvalidAmount = errors.New("amount must be greater than zero")
var ErrInvalidShiftReference = errors.New("invalid shift reference")
