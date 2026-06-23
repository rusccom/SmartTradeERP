package storefront

// ResolvedTenant is the outcome of mapping a request Host to a tenant.
type ResolvedTenant struct {
	TenantID string
	Status   string
}
