package storefront

import "errors"

// ErrTenantNotFound is returned when a request Host maps to no active
// storefront. The resolver translates it into a fail-closed 404 so a disabled
// or unknown host never falls back to another tenant's shop.
var ErrTenantNotFound = errors.New("storefront: no active tenant for host")

// ErrProductNotFound is returned when a slug matches no published product for
// the resolved tenant. Handlers render the themed 404 page with a 404 status.
var ErrProductNotFound = errors.New("storefront: product not found")
