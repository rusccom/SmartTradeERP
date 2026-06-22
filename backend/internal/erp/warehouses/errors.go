package warehouses

import "errors"

var ErrHasMovements = errors.New("warehouse has movements")
var ErrMustKeepDefault = errors.New("tenant must have default warehouse")
