package storefront

import (
	"slices"

	"smarterp/backend/internal/platform/plugins"
	"smarterp/backend/internal/shared/db"
)

// Recheck the grant in the cart/order transaction so a revoked customization cannot be committed.
func (sh *shop) checkPluginState(tx db.DBTX) error {
	key := sh.service.opts.BundleOptions.Key
	if key == "" {
		return nil
	}
	enabled, err := plugins.Enabled(sh.ctx, tx, sh.tenantID, key)
	if err == nil && enabled != slices.Contains(sh.brand.Plugins, key) {
		return errCartChanged
	}
	return err
}
