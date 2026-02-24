package documents

import "errors"

var ErrStatusConflict = errors.New("document status conflict")
var ErrDraftOnly = errors.New("operation allowed only for draft")
var ErrPostedOnly = errors.New("operation allowed only for posted")
var ErrCompositeWithoutComponents = errors.New("composite variant has no components")
