package personal

import (
	"context"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"smarterp/backend/internal/integrations/allegro/connection"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/httpx"
)

const Path = connection.Path + "/private"

type input struct {
	Version string `json:"version"`
	ID      string `json:"id"`
	Refresh string `json:"refresh_token"`
}

type operation func(context.Context, actor, input) (any, error)

func (s *Service) handle(w http.ResponseWriter, r *http.Request, op operation) {
	claims := auth.ClaimsFromContext(r.Context())
	if !s.enabled() || claims.TenantID != s.cfg.TenantID || claims.Role != "owner" {
		connection.Respond(w, nil, forbidden)
		return
	}
	var in input
	r.Body = http.MaxBytesReader(w, r.Body, 32<<10)
	if httpx.DecodeJSON(r, &in) != nil {
		connection.Respond(w, nil, badInput)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()
	result, err := op(ctx, actor{claims.TenantID, claims.UserID}, in)
	connection.Respond(w, result, err)
}

func (s *Service) Start(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r, func(ctx context.Context, who actor, in input) (any, error) {
		if in.Version != "" && uuid.Validate(in.Version) != nil {
			return nil, connection.Changed
		}
		return s.start(ctx, who, in.Version)
	})
}

func (s *Service) Poll(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r, func(ctx context.Context, who actor, in input) (any, error) {
		if uuid.Validate(in.ID) != nil {
			return nil, connection.State
		}
		return s.poll(ctx, who, in.ID)
	})
}

func (s *Service) Import(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r, func(ctx context.Context, who actor, in input) (any, error) {
		token := strings.TrimSpace(in.Refresh)
		if token == "" || len(token) > 16384 || strings.IndexFunc(token, unicode.IsSpace) >= 0 {
			return nil, badInput
		}
		return nil, s.importRefresh(ctx, who, in.Version, token)
	})
}

var badInput = &connection.Problem{Code: "private_input", Status: 400, Message: "Enter a valid private application authorization."}
