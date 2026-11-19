package storefront

import "smarterp/backend/internal/platform/storefront/catalog"

// Options carries deploy-level configuration the storefront engine needs but
// that is not stored per tenant.
type Options struct {
	BundleOptions catalog.RecipePlugin
	// MediaBaseURL is the R2 public base URL used to build product image URLs.
	MediaBaseURL string
	// JWTSecret verifies preview tokens minted by the ERP API so owners can see
	// their unpublished draft. Empty disables preview (published only).
	JWTSecret string
}
