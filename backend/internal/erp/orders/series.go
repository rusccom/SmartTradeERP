package orders

import (
	"context"
	"smarterp/backend/internal/shared/db"
)

type NumberSeries struct {
	Ctx    context.Context
	Tx     db.DBTX
	Tenant string
	Prefix string
	Year   int
}
