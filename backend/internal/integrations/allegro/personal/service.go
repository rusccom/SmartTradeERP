package personal

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/integrations/allegro/connection"
	"smarterp/backend/internal/shared/config"
)

type Service struct {
	*connection.Service
	cfg config.AllegroPrivate
	app connection.Application
}

func New(s *connection.Service, cfg config.AllegroPrivate) *Service {
	return &Service{s, cfg, connection.Application{Client: s.Client, ID: cfg.ClientID, Secret: cfg.ClientSecret}}
}

func (s *Service) enabled() bool {
	return uuid.Validate(s.cfg.TenantID) == nil && s.cfg.SellerID != "" && s.cfg.ClientID != "" && s.cfg.ClientSecret != "" && s.Vault.AEAD != nil
}

func (s *Service) Renewal() connection.PrivateRenewal {
	if !s.enabled() {
		return connection.PrivateRenewal{}
	}
	return connection.PrivateRenewal{TenantID: s.cfg.TenantID, SellerID: s.cfg.SellerID, ClientID: s.cfg.ClientID, Renew: s.refresh}
}

func (s *Service) refresh(ctx context.Context, token string) (connection.Credential, error) {
	pair, err := s.app.Refresh(ctx, token)
	if err == nil {
		err = validateScopes(pair.Access, s.cfg.ClientID)
	}
	return pair, err
}

type actor struct{ tenant, owner string }

func (s *Service) checkTarget(ctx context.Context, tx pgx.Tx, who actor, version string) error {
	const query = `SELECT COALESCE(c.version::text,'') AS version,COALESCE(c.seller_id,'') AS seller,COALESCE(c.method,'') AS method
		FROM platform.tenant_users u LEFT JOIN platform.allegro_connections c ON c.tenant_id=u.tenant_id
		WHERE u.id=$1 AND u.tenant_id=$2 AND u.role='owner' AND u.is_active`
	var current, seller, method string
	err := tx.QueryRow(ctx, query, who.owner, who.tenant).Scan(&current, &seller, &method)
	if err == pgx.ErrNoRows {
		return forbidden
	}
	if err != nil {
		return err
	}
	if current != version {
		return connection.Changed
	}
	if current != "" && (seller != s.cfg.SellerID || method != string(connection.Manual)) {
		return connection.Different
	}
	return nil
}

func (s *Service) savePair(ctx context.Context, tx pgx.Tx, who actor, pair connection.Credential, version string) error {
	if err := validateScopes(pair.Access, s.cfg.ClientID); err != nil {
		return err
	}
	seller, err := s.Client.Verify(ctx, pair.Access)
	if err != nil {
		return err
	}
	if seller.ID != s.cfg.SellerID {
		return connection.Different
	}
	g := connection.Grant{Seller: seller, Credential: pair, Method: connection.Manual, Version: version, PrivateClient: s.cfg.ClientID}
	return s.SaveTx(ctx, tx, who.tenant, g)
}

func (s *Service) importRefresh(ctx context.Context, who actor, version, token string) error {
	return s.DB.WithTenantTx(ctx, who.tenant, func(tx pgx.Tx) error {
		if err := s.checkTarget(ctx, tx, who, version); err != nil {
			return err
		}
		pair, err := s.refresh(ctx, token)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		pair.RefreshedAt = &now
		return s.savePair(ctx, tx, who, pair, version)
	})
}

var forbidden = &connection.Problem{Code: "private_unavailable", Status: 403, Message: "Private application renewal is unavailable for this organization."}
