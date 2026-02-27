package documents

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
)

type Repository struct {
	store *db.Store
}

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) List(ctx context.Context, tenantID string, query httpx.ListQuery) ([]ListItem, int, error) {
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
	sqlQuery := `SELECT COUNT(*) FROM documents.documents d WHERE d.tenant_id=$1`
	args := []any{tenantID}
	sqlQuery, args = appendListFilters(sqlQuery, args, query)
	row := r.store.Pool.QueryRow(ctx, sqlQuery, args...)
	total := 0
	return total, row.Scan(&total)
}

func (r *Repository) load(ctx context.Context, tenantID string, query httpx.ListQuery) ([]ListItem, error) {
	sqlQuery := `SELECT d.id::text, d.type, d.date::text, COALESCE(d.number,''), d.status,
        COALESCE(d.customer_id::text,''), COALESCE(d.note,''),
        ` + totalCostSQL() + `::text AS total_cost
        FROM documents.documents d WHERE d.tenant_id=$1`
	args := []any{tenantID}
	sqlQuery, args = appendListFilters(sqlQuery, args, query)
	sqlQuery, args = appendSortAndPaging(sqlQuery, args, query)
	rows, err := r.store.Pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanList(rows)
}

func appendListFilters(query string, args []any, listQuery httpx.ListQuery) (string, []any) {
	query, args = appendFilter(query, args, "d.type", listQuery.Filters["type"])
	query, args = appendFilter(query, args, "d.status", listQuery.Filters["status"])
	query, args = appendFilter(query, args, "d.date::text", listQuery.Filters["date"])
	query, args = appendSearch(query, args, listQuery.Search)
	return query, args
}

func appendFilter(query string, args []any, field, value string) (string, []any) {
	if value == "" {
		return query, args
	}
	query += ` AND ` + field + `=$` + nextPosition(args)
	args = append(args, value)
	return query, args
}

func appendSearch(query string, args []any, search string) (string, []any) {
	if search == "" {
		return query, args
	}
	query += ` AND COALESCE(d.number,'') ILIKE '%' || $` + nextPosition(args) + ` || '%'`
	args = append(args, search)
	return query, args
}

func appendSortAndPaging(query string, args []any, listQuery httpx.ListQuery) (string, []any) {
	query += ` ORDER BY ` + sortField(listQuery.SortBy) + ` ` + sortDir(listQuery.SortDir)
	query += `, d.created_at DESC`
	query += ` LIMIT $` + nextPosition(args)
	query += ` OFFSET $` + strconv.Itoa(len(args)+2)
	offset := httpx.Offset(listQuery.Page, listQuery.PerPage)
	args = append(args, listQuery.PerPage, offset)
	return query, args
}

func scanList(rows pgx.Rows) ([]ListItem, error) {
	items := make([]ListItem, 0)
	for rows.Next() {
		item := ListItem{}
		err := rows.Scan(&item.ID, &item.Type, &item.Date, &item.Number, &item.Status, &item.CustomerID, &item.Note, &item.TotalCost)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func nextPosition(args []any) string {
	return strconv.Itoa(len(args) + 1)
}

func sortField(sortBy string) string {
	if sortBy == "number" {
		return "COALESCE(d.number,'')"
	}
	if sortBy == "total_cost" {
		return totalCostSQL()
	}
	return "d.date"
}

func sortDir(direction string) string {
	if direction == "asc" {
		return "asc"
	}
	return "desc"
}

func totalCostSQL() string {
	return `COALESCE((SELECT SUM(di.qty * di.unit_price)
        FROM documents.document_items di
        WHERE di.document_id = d.id), 0)`
}

func (r *Repository) ByID(ctx context.Context, tenantID, id string) (Document, error) {
	query := `SELECT id::text, type, date::text, COALESCE(number,''), status,
        COALESCE(warehouse_id::text,''), COALESCE(source_warehouse_id::text,''),
        COALESCE(target_warehouse_id::text,''), COALESCE(shift_id::text,''),
        COALESCE(customer_id::text,''), COALESCE(note,'')
        FROM documents.documents
        WHERE tenant_id=$1 AND id=$2`
	row := r.store.Pool.QueryRow(ctx, query, tenantID, id)
	item := Document{}
	err := row.Scan(&item.ID, &item.Type, &item.Date, &item.Number, &item.Status,
		&item.WarehouseID, &item.SourceWarehouseID, &item.TargetWarehouseID,
		&item.ShiftID, &item.CustomerID, &item.Note)
	return item, err
}

func (r *Repository) LoadItemsWithProfit(ctx context.Context, tenantID, documentID string) ([]DocumentItem, decimal.Decimal, error) {
	query := `SELECT i.id::text, i.variant_id::text, i.qty, i.unit_price, i.total_amount,
        COALESCE(SUM(l.profit),0) AS profit
        FROM documents.document_items i
        JOIN documents.documents d ON d.id=i.document_id
        LEFT JOIN ledger.cost_ledger l ON l.document_item_id=i.id AND l.tenant_id=d.tenant_id
        WHERE d.tenant_id=$1 AND d.id=$2
        GROUP BY i.id, i.variant_id, i.qty, i.unit_price, i.total_amount
        ORDER BY i.id`
	rows, err := r.store.Pool.Query(ctx, query, tenantID, documentID)
	if err != nil {
		return nil, decimal.Zero, err
	}
	defer rows.Close()
	return scanItems(rows)
}

func scanItems(rows pgx.Rows) ([]DocumentItem, decimal.Decimal, error) {
	items := make([]DocumentItem, 0)
	total := decimal.Zero
	for rows.Next() {
		item := DocumentItem{}
		err := rows.Scan(&item.ID, &item.VariantID, &item.Qty, &item.UnitPrice, &item.TotalAmount, &item.Profit)
		if err != nil {
			return nil, decimal.Zero, err
		}
		total = total.Add(item.Profit)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *Repository) InsertDocument(ctx context.Context, tx pgx.Tx, tenantID, documentID string, req CreateRequest) error {
	query := `INSERT INTO documents.documents
        (id, tenant_id, type, date, number, status, warehouse_id, source_warehouse_id, target_warehouse_id, shift_id, customer_id, note)
        VALUES ($1,$2,$3,$4,$5,'draft',NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),NULLIF($9,''),NULLIF($10,''),$11)`
	_, err := tx.Exec(ctx, query, documentID, tenantID, req.Type, req.Date,
		req.Number, req.WarehouseID, req.SourceWarehouseID, req.TargetWarehouseID,
		req.ShiftID, req.CustomerID, req.Note)
	return err
}

func (r *Repository) UpdateDocument(ctx context.Context, tx pgx.Tx, tenantID, documentID string, req UpdateRequest) error {
	query := `UPDATE documents.documents
        SET type=$3, date=$4, number=$5, warehouse_id=NULLIF($6,''),
            source_warehouse_id=NULLIF($7,''), target_warehouse_id=NULLIF($8,''),
            shift_id=NULLIF($9,''), customer_id=NULLIF($10,''), note=$11, updated_at=now()
        WHERE tenant_id=$1 AND id=$2`
	_, err := tx.Exec(ctx, query, tenantID, documentID, req.Type, req.Date, req.Number,
		req.WarehouseID, req.SourceWarehouseID, req.TargetWarehouseID,
		req.ShiftID, req.CustomerID, req.Note)
	return err
}

func (r *Repository) ReplaceItems(ctx context.Context, tx pgx.Tx, documentID string, items []ItemInput) error {
	if err := r.deleteItems(ctx, tx, documentID); err != nil {
		return err
	}
	for _, item := range items {
		if err := r.insertItem(ctx, tx, documentID, item); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) deleteItems(ctx context.Context, tx pgx.Tx, documentID string) error {
	_, err := tx.Exec(ctx, `DELETE FROM documents.document_items WHERE document_id=$1`, documentID)
	return err
}

func (r *Repository) insertItem(ctx context.Context, tx pgx.Tx, documentID string, item ItemInput) error {
	id := uuid.NewString()
	total := item.Qty.Mul(item.UnitPrice).Round(4)
	query := `INSERT INTO documents.document_items
        (id, document_id, variant_id, qty, unit_price, total_amount)
        VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := tx.Exec(ctx, query, id, documentID, item.VariantID, item.Qty, item.UnitPrice, total)
	return err
}

func (r *Repository) Status(ctx context.Context, tenantID, documentID string) (string, error) {
	row := r.store.Pool.QueryRow(ctx,
		`SELECT status FROM documents.documents WHERE tenant_id=$1 AND id=$2`, tenantID, documentID)
	status := ""
	return status, row.Scan(&status)
}

func (r *Repository) SetStatus(ctx context.Context, tx pgx.Tx, tenantID, documentID, status string) error {
	query := `UPDATE documents.documents SET status=$3, updated_at=now() WHERE tenant_id=$1 AND id=$2`
	_, err := tx.Exec(ctx, query, tenantID, documentID, status)
	return err
}

func (r *Repository) DeleteDraft(ctx context.Context, tenantID, id string) error {
	query := `DELETE FROM documents.documents WHERE tenant_id=$1 AND id=$2 AND status='draft'`
	_, err := r.store.Pool.Exec(ctx, query, tenantID, id)
	return err
}
