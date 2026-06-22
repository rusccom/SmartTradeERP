package documents

import (
	"context"

	"smarterp/backend/internal/shared/auth"
)

// actorID returns the user behind the current request, taken from JWT claims
// that the auth middleware stored in the context. Empty when no claims exist.
func actorID(ctx context.Context) string {
	return auth.ClaimsFromContext(ctx).UserID
}
