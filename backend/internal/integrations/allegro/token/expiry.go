package token

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// The provider verifies the bearer token. Its optional exp claim is only used for display.
func expiry(token string) *time.Time {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil
	}
	body, err := base64.RawURLEncoding.DecodeString(parts[1])
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err != nil || json.Unmarshal(body, &claims) != nil || claims.Exp <= 0 || claims.Exp > 253402300799 {
		return nil
	}
	until := time.Unix(claims.Exp, 0).UTC()
	return &until
}
