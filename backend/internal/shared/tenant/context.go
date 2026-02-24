package tenant

import "context"

type key string

const tenantKey key = "tenant_id"

func WithTenantID(ctx context.Context, tenantID string) context.Context {
    return context.WithValue(ctx, tenantKey, tenantID)
}

func FromContext(ctx context.Context) string {
    value := ctx.Value(tenantKey)
    tenantID, _ := value.(string)
    return tenantID
}
