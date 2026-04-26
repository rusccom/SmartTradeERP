package documents

import "errors"

var ErrStatusConflict = errors.New("document status conflict")
var ErrDraftOnly = errors.New("operation allowed only for draft")
var ErrPostedOnly = errors.New("operation allowed only for posted")
var ErrCompositeWithoutComponents = errors.New("composite variant has no components")
var ErrInvalidDocumentReference = errors.New("invalid document reference")
var ErrPaymentsRequired = errors.New("payments are required")
var ErrPaymentsNotAllowed = errors.New("payments are not allowed")
var ErrInvalidPaymentMethod = errors.New("invalid payment method")
var ErrInvalidPaymentAmount = errors.New("invalid payment amount")
var ErrPaymentTotalMismatch = errors.New("payment total mismatch")
var ErrDocumentNumberConflict = errors.New("document number already exists")
