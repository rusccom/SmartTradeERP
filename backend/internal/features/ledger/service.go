package ledger

import (
	"smarterp/backend/internal/shared/db"
)

type Service struct {
	store *db.Store
}

func NewService(store *db.Store) *Service {
	return &Service{store: store}
}
