package bundles

import "errors"

var ErrInvalidComponentState = errors.New("invalid component state")
var ErrMissingComponents = errors.New("bundle has no components")
