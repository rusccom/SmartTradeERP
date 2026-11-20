package personal

import (
	"context"
	"log"
	"time"

	"smarterp/backend/internal/integrations/allegro/connection"
)

func (s *Service) Run(ctx context.Context) {
	if !s.enabled() {
		return
	}
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		s.maintain(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *Service) maintain(parent context.Context) {
	ctx, cancel := context.WithTimeout(parent, 25*time.Second)
	defer cancel()
	err := s.MaintainPrivate(ctx, s.Renewal())
	if err != nil && err != connection.Missing && ctx.Err() == nil {
		log.Printf("Allegro private renewal: %s", connection.PublicError(err).Code)
	}
}
