package personal

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"smarterp/backend/internal/integrations/allegro/connection"
)

// Allegro issued this token over TLS. These claims restrict the grant to the configured application and permissions.
func validateScopes(token, clientID string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return connection.ReadOnly
	}
	body, err := base64.RawURLEncoding.DecodeString(parts[1])
	var claims struct {
		ClientID string          `json:"client_id"`
		Scope    json.RawMessage `json:"scope"`
	}
	if err != nil || json.Unmarshal(body, &claims) != nil {
		return connection.ReadOnly
	}
	if claims.ClientID != "" && claims.ClientID != clientID {
		return connection.ReadOnly
	}
	return readOnly(claims.Scope)
}

func readOnly(raw json.RawMessage) error {
	var list []string
	if json.Unmarshal(raw, &list) != nil {
		var value string
		if json.Unmarshal(raw, &value) != nil {
			return connection.ReadOnly
		}
		list = strings.Fields(value)
	}
	rights := make(map[string]bool)
	for _, item := range list {
		if strings.HasPrefix(item, "allegro:api:") {
			rights[item] = true
		}
	}
	if len(rights) != 2 || !rights["allegro:api:profile:read"] || !rights["allegro:api:orders:read"] {
		return connection.ReadOnly
	}
	return nil
}
