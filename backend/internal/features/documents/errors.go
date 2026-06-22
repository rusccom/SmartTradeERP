package documents

import "errors"

var ErrStatusConflict = errors.New("document status conflict")
var ErrDraftOnly = errors.New("operation allowed only for draft")
var ErrPostedOnly = errors.New("operation allowed only for posted")
var ErrInvalidDocumentReference = errors.New("invalid document reference")
var ErrPaymentsRequired = errors.New("payments are required")
var ErrPaymentsNotAllowed = errors.New("payments are not allowed")
var ErrInvalidPaymentMethod = errors.New("invalid payment method")
var ErrInvalidPaymentAmount = errors.New("invalid payment amount")
var ErrPaymentTotalMismatch = errors.New("payment total mismatch")
var ErrDocumentNumberConflict = errors.New("document number already exists")
var ErrBundleDocumentType = errors.New("bundle is allowed only for sale and return")
var ErrShiftDocumentLocked = errors.New("shift document cannot be modified")
var ErrShiftClosed = errors.New("shift is closed")
var ErrTypeImmutable = errors.New("document type cannot be changed")
