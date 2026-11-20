package config

import (
	"os"
	"strings"
)

type Allegro struct {
	Private         AllegroPrivate
	ClientID        string
	ClientSecret    string
	EncryptionKey   string
	APIOrigin       string
	DashboardOrigin string
}

type AllegroPrivate struct {
	TenantID     string
	SellerID     string
	ClientID     string
	ClientSecret string
}

func loadAllegro() Allegro {
	return Allegro{
		Private: AllegroPrivate{
			TenantID:     strings.TrimSpace(os.Getenv("ALLEGRO_PRIVATE_TENANT_ID")),
			SellerID:     strings.TrimSpace(os.Getenv("ALLEGRO_PRIVATE_SELLER_ID")),
			ClientID:     strings.TrimSpace(os.Getenv("ALLEGRO_PRIVATE_CLIENT_ID")),
			ClientSecret: strings.TrimSpace(os.Getenv("ALLEGRO_PRIVATE_CLIENT_SECRET")),
		},
		ClientID:        strings.TrimSpace(os.Getenv("ALLEGRO_OAUTH_CLIENT_ID")),
		ClientSecret:    strings.TrimSpace(os.Getenv("ALLEGRO_OAUTH_CLIENT_SECRET")),
		EncryptionKey:   strings.TrimSpace(os.Getenv("ALLEGRO_ENCRYPTION_KEY")),
		APIOrigin:       strings.TrimRight(os.Getenv("ALLEGRO_API_ORIGIN"), "/"),
		DashboardOrigin: strings.TrimRight(os.Getenv("ALLEGRO_DASHBOARD_ORIGIN"), "/"),
	}
}
