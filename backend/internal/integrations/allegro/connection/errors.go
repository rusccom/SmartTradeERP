package connection

import (
	"errors"
	"net/http"

	"smarterp/backend/internal/shared/httpx"
)

type Problem struct {
	Code    string
	Status  int
	Message string
}

func (p *Problem) Error() string { return p.Message }

var (
	Setup        = &Problem{"setup", 503, "This connection method has not been configured."}
	State        = &Problem{"expired", 409, "This authorization expired or was replaced. Start a new connection."}
	Denied       = &Problem{"cancelled", 422, "Allegro access was not granted."}
	Unauthorized = &Problem{"invalid_token", 422, "Allegro rejected the access token. It may have expired or been revoked."}
	Scope        = &Problem{"permissions", 422, "Allegro did not grant read access to both the seller profile and order history."}
	Different    = &Problem{"different", 409, "A different seller or method is connected. Disconnect it before switching."}
	Changed      = &Problem{"changed", 409, "The connection changed. Reload the settings and try again."}
	Missing      = &Problem{"missing", 404, "No Allegro account is connected."}
	Remote       = &Problem{"unavailable", 502, "Allegro could not verify the connection. Please try again."}
	Limit        = &Problem{"rate_limit", 429, "Allegro has limited requests. Please try again later."}
	Pending      = &Problem{"pending", 409, "Confirm this seller on Allegro to finish authorization."}
	Slow         = &Problem{"slow_down", 429, "Wait before checking the Allegro authorization again."}
	ReadOnly     = &Problem{"scope_mismatch", 422, "Authorize this application for seller profile and order read access only."}
	Stored       = &Problem{"storage", 500, "The saved credential cannot be opened."}
	Invalid      = &Problem{"invalid", 400, "Enter an Allegro seller access token."}
	internal     = &Problem{"failed", 500, "Could not update the Allegro connection."}
)

func PublicError(err error) *Problem {
	var p *Problem
	if errors.As(err, &p) {
		return p
	}
	return internal
}

func Respond(w http.ResponseWriter, data any, err error) {
	w.Header().Set("Cache-Control", "private, no-store")
	if err != nil {
		p := PublicError(err)
		httpx.WriteError(w, p.Status, "allegro_"+p.Code, p.Message, nil)
		return
	}
	httpx.WriteData(w, http.StatusOK, data, nil)
}
