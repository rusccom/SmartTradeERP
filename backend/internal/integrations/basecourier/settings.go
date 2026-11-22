package basecourier

import (
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type connection struct {
	Email      string    `json:"email"`
	Version    uuid.UUID `json:"version"`
	VerifiedAt time.Time `json:"verified_at"`
}

type update struct {
	Email   string     `json:"email"`
	APIKey  string     `json:"api_key"`
	Version *uuid.UUID `json:"version"`
}

func (u *update) validate() error {
	u.Email = strings.TrimSpace(u.Email)
	u.APIKey = strings.TrimSpace(u.APIKey)
	address, err := mail.ParseAddress(u.Email)
	if err != nil || address.Address != u.Email || len(u.Email) > 254 {
		return errInvalid
	}
	if u.APIKey == "" || len(u.APIKey) > 4096 || strings.IndexFunc(u.APIKey, unicode.IsControl) >= 0 {
		return errInvalid
	}
	return nil
}
