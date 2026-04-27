package ledger

import "errors"

var ErrNegativeStock = errors.New("negative stock is not allowed")
var ErrActiveBatchRequired = errors.New("active posting batch is required")
