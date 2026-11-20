package config

import (
	"os"
	"strings"
)

type Allegro struct {
	ClientID        string
	ClientSecret    string
	EncryptionKey   string
	APIOrigin       string
	DashboardOrigin string
}

func loadAllegro() Allegro {
	return Allegro{
		ClientID:        strings.TrimSpace(os.Getenv("ALLEGRO_OAUTH_CLIENT_ID")),
		ClientSecret:    strings.TrimSpace(os.Getenv("ALLEGRO_OAUTH_CLIENT_SECRET")),
		EncryptionKey:   strings.TrimSpace(os.Getenv("ALLEGRO_ENCRYPTION_KEY")),
		APIOrigin:       strings.TrimRight(os.Getenv("ALLEGRO_API_ORIGIN"), "/"),
		DashboardOrigin: strings.TrimRight(os.Getenv("ALLEGRO_DASHBOARD_ORIGIN"), "/"),
	}
}
