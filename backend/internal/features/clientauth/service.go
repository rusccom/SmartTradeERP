package clientauth

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

type Service struct {
	store  *db.Store
	repo   *Repository
	tokens *auth.TokenService
}

func NewService(store *db.Store, repo *Repository, tokens *auth.TokenService) *Service {
	return &Service{store: store, repo: repo, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (auth.TokenResponse, error) {
	user, err := s.repo.FindByEmail(ctx, s.store.Pool, req.Email)
	if err != nil {
		return auth.TokenResponse{}, err
	}
	if !auth.VerifyPassword(req.Password, user.PasswordHash) {
		return auth.TokenResponse{}, ErrInvalidCredentials
	}
	return s.tokens.Issue(user.ID, user.TenantID, user.Role, "client")
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (auth.TokenResponse, error) {
	ids := createIDs()
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return auth.TokenResponse{}, err
	}
	err = s.store.WithTx(ctx, func(tx pgx.Tx) error {
		return s.createTenantGraph(ctx, tx, ids, req, hash)
	})
	if err != nil {
		return auth.TokenResponse{}, err
	}
	return s.tokens.Issue(ids.userID, ids.tenantID, "owner", "client")
}

type registerIDs struct {
	tenantID    string
	userID      string
	warehouseID string
}

func createIDs() registerIDs {
	ids := registerIDs{}
	ids.tenantID = uuid.NewString()
	ids.userID = uuid.NewString()
	ids.warehouseID = uuid.NewString()
	return ids
}

func (s *Service) createTenantGraph(
	ctx context.Context,
	tx pgx.Tx,
	ids registerIDs,
	req RegisterRequest,
	hash string,
) error {
	if err := insertTenant(ctx, tx, ids.tenantID, req.TenantName); err != nil {
		return err
	}
	if err := insertOwner(ctx, tx, ids, req.Email, hash); err != nil {
		return err
	}
	return insertDefaultWarehouse(ctx, tx, ids)
}

func insertTenant(ctx context.Context, tx pgx.Tx, tenantID, name string) error {
	query := `INSERT INTO platform.tenants (id, name, status, plan)
        VALUES ($1,$2,'trial','free')`
	_, err := tx.Exec(ctx, query, tenantID, name)
	return err
}

func insertOwner(ctx context.Context, tx pgx.Tx, ids registerIDs, email, hash string) error {
	query := `INSERT INTO platform.tenant_users (id, tenant_id, email, password_hash, role)
        VALUES ($1,$2,$3,$4,'owner')`
	_, err := tx.Exec(ctx, query, ids.userID, ids.tenantID, email, hash)
	return err
}

func insertDefaultWarehouse(ctx context.Context, tx pgx.Tx, ids registerIDs) error {
	query := `INSERT INTO catalog.warehouses (id, tenant_id, name, is_default)
        VALUES ($1,$2,'Main Warehouse',true)`
	_, err := tx.Exec(ctx, query, ids.warehouseID, ids.tenantID)
	return err
}
