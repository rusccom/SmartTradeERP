package storefront

import "context"

type Service struct {
	repo     *Repository
	registry *Registry
	opts     Options
}

func NewService(repo *Repository, registry *Registry, opts Options) *Service {
	return &Service{repo: repo, registry: registry, opts: opts}
}

// ResolveHost maps a normalized hostname to an active tenant. Suspended
// tenants are reported as not-found so a disabled account serves no shop.
func (s *Service) ResolveHost(ctx context.Context, host string) (ResolvedTenant, error) {
	if host == "" {
		return ResolvedTenant{}, ErrTenantNotFound
	}
	resolved, err := s.repo.FindTenantByHost(ctx, host)
	if err != nil {
		return ResolvedTenant{}, err
	}
	if resolved.Status == "suspended" {
		return ResolvedTenant{}, ErrTenantNotFound
	}
	return resolved, nil
}
