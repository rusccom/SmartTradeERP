package shifts

import "errors"

var ErrShiftAlreadyOpen = errors.New("у вас уже есть открытая смена")
var ErrNoOpenShift = errors.New("нет открытой смены")
var ErrShiftAlreadyClosed = errors.New("смена уже закрыта")
var ErrInvalidCashOpType = errors.New("тип операции должен быть cash_in или cash_out")
var ErrInvalidAmount = errors.New("сумма должна быть больше нуля")
