package customers

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
)

type Repository struct {
	store *db.Store
}

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) List(ctx context.Context, tenantID string, query httpx.ListQuery) ([]Customer, int, error) {
	total, err := r.count(ctx, tenantID, query)
	if err != nil {
		return nil, 0, err
	}
	data, err := r.load(ctx, tenantID, query)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (r *Repository) count(ctx context.Context, tenantID string, query httpx.ListQuery) (int, error) {
	sql := `SELECT COUNT(*) FROM catalog.customers WHERE tenant_id=$1`
	args := []any{tenantID}
	sql, args = addSearchFilter(sql, args, query.Search)
	row := r.store.Pool.QueryRow(ctx, sql, args...)
	total := 0
	return total, row.Scan(&total)
}

func (r *Repository) load(ctx context.Context, tenantID string, query httpx.ListQuery) ([]Customer, error) {
	sql := `SELECT id::text, name, COALESCE(phone,''), COALESCE(email,''),
		is_default, created_at::text, updated_at::text
		FROM catalog.customers WHERE tenant_id=$1`
	args := []any{tenantID}
	sql, args = addSearchFilter(sql, args, query.Search)
	sql, args = appendSortAndPaging(sql, args, query)
	rows, err := r.store.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCustomers(rows)
}

func addSearchFilter(sql string, args []any, search string) (string, []any) {
	if search == "" {
		return sql, args
	}
	p := pos(len(args) + 1)
	sql += ` AND (name ILIKE '%' || $` + p + ` || '%' OR phone ILIKE '%' || $` + p + ` || '%')`
	args = append(args, search)
	return sql, args
}

func appendSortAndPaging(sql string, args []any, q httpx.ListQuery) (string, []any) {
	sortBy, sortDir := q.SortBy, q.SortDir
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}
	sql += ` ORDER BY ` + sortBy + ` ` + sortDir
	sql += ` LIMIT $` + pos(len(args)+1)
	sql += ` OFFSET $` + pos(len(args)+2)
	args = append(args, q.PerPage, httpx.Offset(q.Page, q.PerPage))
	return sql, args
}

func pos(v int) string { return strconv.Itoa(v) }

func scanCustomers(rows pgx.Rows) ([]Customer, error) {
	items := make([]Customer, 0)
	for rows.Next() {
		c := Customer{}
		err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (r *Repository) Create(ctx context.Context, tenantID, id string, req CreateRequest) error {
	sql := `INSERT INTO catalog.customers (id, tenant_id, name, phone, email)
		VALUES ($1,$2,$3,$4,$5)`
	_, err := r.store.Pool.Exec(ctx, sql, id, tenantID, req.Name, req.Phone, req.Email)
	return err
}

func (r *Repository) GetByID(ctx context.Context, tenantID, id string) (Customer, error) {
	sql := `SELECT id::text, name, COALESCE(phone,''), COALESCE(email,''),
		is_default, created_at::text, updated_at::text
		FROM catalog.customers WHERE tenant_id=$1 AND id=$2`
	row := r.store.Pool.QueryRow(ctx, sql, tenantID, id)
	c := Customer{}
	err := row.Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repository) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
	sql := `UPDATE catalog.customers SET name=$3, phone=$4, email=$5, updated_at=now()
		WHERE tenant_id=$1 AND id=$2`
	_, err := r.store.Pool.Exec(ctx, sql, tenantID, id, req.Name, req.Phone, req.Email)
	return err
}

func (r *Repository) Delete(ctx context.Context, tenantID, id string) error {
	sql := `DELETE FROM catalog.customers WHERE tenant_id=$1 AND id=$2 AND is_default=false`
	_, err := r.store.Pool.Exec(ctx, sql, tenantID, id)
	return err
}

func (r *Repository) IsDefault(ctx context.Context, tenantID, id string) (bool, error) {
	sql := `SELECT is_default FROM catalog.customers WHERE tenant_id=$1 AND id=$2`
	row := r.store.Pool.QueryRow(ctx, sql, tenantID, id)
	var isDefault bool
	return isDefault, row.Scan(&isDefault)
}

func (r *Repository) HasDocuments(ctx context.Context, tenantID, customerID string) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM documents.documents WHERE tenant_id=$1 AND customer_id=$2)`
	row := r.store.Pool.QueryRow(ctx, sql, tenantID, customerID)
	var exists bool
	return exists, row.Scan(&exists)
}
