package basecourier

import (
	"errors"
	"net/http"

	"smarterp/backend/internal/shared/httpx"
)

type problem struct {
	status  int
	code    string
	message string
}

func (p *problem) Error() string { return p.message }

var (
	errInvalid = &problem{400, "invalid_settings", "Enter your Base Courier account email and API key."}
	errVersion = &problem{400, "invalid_version", "Invalid Base Courier connection version."}
	errChanged = &problem{409, "connection_changed", "This connection changed. Reload its settings and try again."}
	errDenied  = &problem{422, "credentials_rejected", "Base Courier rejected these credentials. Check the account email and API key."}
	errLimited = &problem{429, "rate_limited", "Base Courier has limited verification requests. Try again later."}
	errRemote  = &problem{502, "provider_unavailable", "Base Courier could not verify the account right now. Try again."}
	errVault   = &problem{503, "key_storage_unavailable", "Base Courier credential storage is unavailable."}
	errUnknown = &problem{500, "settings_unavailable", "Unable to access the Base Courier connection."}
)

func respond(w http.ResponseWriter, data any, err error) {
	w.Header().Set("Cache-Control", "private, no-store")
	if err == nil {
		httpx.WriteData(w, http.StatusOK, data, nil)
		return
	}
	var p *problem
	if !errors.As(err, &p) {
		p = errUnknown
	}
	httpx.WriteError(w, p.status, "base_courier_"+p.code, p.message, nil)
}
