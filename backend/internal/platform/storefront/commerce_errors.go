package storefront

import "errors"

var (
	// ErrCartEmpty means no submitted item resolved to a purchasable variant.
	ErrCartEmpty = errors.New("storefront: no purchasable items")
	// ErrCheckoutInvalid means the customer details failed validation.
	ErrCheckoutInvalid = errors.New("storefront: invalid checkout data")
)
