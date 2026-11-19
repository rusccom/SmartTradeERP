package inpost

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	errInvalid     = errors.New("invalid connection settings")
	errToken       = errors.New("a token is required for this account")
	errCredentials = errors.New("ShipX denied access to the organization")
	errCarrier     = errors.New("ShipX verification unavailable")
	errThrottled   = errors.New("ShipX rate limit reached")
	errVault       = errors.New("credential storage unavailable")
	errConflict    = errors.New("settings changed")
)

type account struct {
	Environment    string `json:"environment"`
	OrganizationID int64  `json:"organization_id,string"`
}

type update struct {
	account
	Token    string `json:"token"`
	Revision int64  `json:"revision"`
}

type connection struct {
	account
	OrganizationName string    `json:"organization_name"`
	Services         []string  `json:"services"`
	VerifiedAt       time.Time `json:"verified_at"`
	Revision         int64     `json:"revision"`
}

func (a account) shipxHost() string {
	switch a.Environment {
	case "production":
		return "https://api-shipx-pl.easypack24.net"
	case "sandbox":
		return "https://sandbox-api-shipx-pl.easypack24.net"
	default:
		return ""
	}
}

func (a account) organizationID() string { return strconv.FormatInt(a.OrganizationID, 10) }

func (u update) validate() error {
	if u.shipxHost() == "" || u.OrganizationID < 1 || u.Revision < 0 {
		return errInvalid
	}
	if u.Revision == 0 && u.Token == "" {
		return errToken
	}
	if len(u.Token) > 4096 || strings.IndexFunc(u.Token, func(r rune) bool { return r < 33 || r > 126 }) >= 0 {
		return errInvalid
	}
	return nil
}
