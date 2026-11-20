package connection

import "time"

type Method string

const (
	Manual Method = "token"
	OAuth  Method = "oauth"
)

type Seller struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

type Account struct {
	Seller
	Method         Method     `json:"method"`
	Version        string     `json:"version"`
	VerifiedAt     time.Time  `json:"verified_at"`
	ExpiresAt      *time.Time `json:"expires_at"`
	NeedsReconnect bool       `json:"needs_reconnect"`
	AutoRenew      bool       `json:"auto_renew"`
	RefreshedAt    *time.Time `json:"refreshed_at"`
}

type Tokens struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token,omitempty"`
}

type Credential struct {
	Tokens
	ExpiresAt   *time.Time
	RefreshedAt *time.Time
}

type Grant struct {
	Seller
	Credential
	Method        Method
	Version       string
	PrivateClient string
}

type stored struct {
	Account
	tenant        string
	encrypted     []byte
	privateClient string
}

func binding(tenant string, method Method, id string) string {
	return "allegro:" + string(method) + ":" + tenant + ":" + id
}
