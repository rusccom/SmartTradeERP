package storefrontadmin

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"smarterp/backend/internal/shared/slug"
)

const maxHostAttempts = 10

var errNoFreeHost = errors.New("storefrontadmin: no free shop host")

// reservedLabels are subdomains of the shops base domain used by the platform
// itself (app, API, mail, ...) — a tenant slug must never claim them.
var reservedLabels = map[string]bool{
	"www": true, "api": true, "erp": true, "app": true, "admin": true,
	"mail": true, "smtp": true, "imap": true, "ftp": true, "cdn": true,
	"assets": true, "static": true, "status": true, "docs": true,
	"help": true, "support": true, "blog": true, "shop": true,
}

// ensureHost returns the tenant's shop host, provisioning a {slug}.{base}
// subdomain on first use. With no base domain configured it is a no-op, so
// local development works without DNS.
func (s *Service) ensureHost(ctx context.Context, tenantID string) (string, error) {
	host, err := s.repo.PrimaryHost(ctx, tenantID)
	if err != nil || host != "" || s.shopsBase == "" {
		return host, err
	}
	return s.provisionHost(ctx, tenantID)
}

func (s *Service) provisionHost(ctx context.Context, tenantID string) (string, error) {
	name, err := s.repo.TenantName(ctx, tenantID)
	if err != nil {
		return "", err
	}
	base := slug.Make(name)
	if reservedLabels[base] {
		base = base + "-store"
	}
	for attempt := 0; attempt < maxHostAttempts; attempt++ {
		host := s.candidateHost(base, attempt)
		err := s.repo.InsertDomain(ctx, uuid.NewString(), tenantID, host)
		if err == nil {
			return host, nil
		}
		if !isUniqueViolation(err) {
			return "", err
		}
	}
	return "", errNoFreeHost
}

// candidateHost tries the clean slug first, then numbered variants, and as a
// last resort a uuid fragment that cannot realistically collide.
func (s *Service) candidateHost(base string, attempt int) string {
	label := base
	if attempt > 0 && attempt < maxHostAttempts-1 {
		label = base + "-" + strconv.Itoa(attempt+1)
	}
	if attempt == maxHostAttempts-1 {
		label = base + "-" + uuid.NewString()[:8]
	}
	return label + "." + s.shopsBase
}

func isUniqueViolation(err error) bool {
	pgErr := &pgconn.PgError{}
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func shopURL(host string) string {
	if host == "" {
		return ""
	}
	return "https://" + host
}
