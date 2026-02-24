package documents

import (
    "context"
    "strconv"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/shared/db"
)

type Repository struct {
    store *db.Store
}

func NewRepository(store *db.Store) *Repository {
    return &Repository{store: store}
}

func (r *Repository) List(ctx context.Context, tenantID string, filters Filters, page, perPage int) ([]ListItem, int, error) {
    total, err := r.count(ctx, tenantID, filters)
    if err != nil {
        return nil, 0, err
    }
    data, err := r.load(ctx, tenantID, filters, page, perPage)
    if err != nil {
        return nil, 0, err
    }
    return data, total, nil
}

func (r *Repository) count(ctx context.Context, tenantID string, filters Filters) (int, error) {
    query := `SELECT COUNT(*) FROM documents.documents WHERE tenant_id=$1`
    args := []any{tenantID}
    query, args = appendFilters(query, args, filters)
    row := r.store.Pool.QueryRow(ctx, query, args...)
    total := 0
    return total, row.Scan(&total)
}

func (r *Repository) load(ctx context.Context, tenantID string, filters Filters, page, perPage int) ([]ListItem, error) {
    query := `SELECT id::text, type, date::text, COALESCE(number,''), status, COALESCE(note,'')
        FROM documents.documents WHERE tenant_id=$1`
    args := []any{tenantID}
    query, args = appendFilters(query, args, filters)
    query, args = appendPaging(query, args, page, perPage)
    rows, err := r.store.Pool.Query(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanList(rows)
}

func appendFilters(query string, args []any, filters Filters) (string, []any) {
    query, args = appendFilter(query, args, "type", filters.Type)
    query, args = appendFilter(query, args, "status", filters.Status)
    query, args = appendFilter(query, args, "date::text", filters.Date)
    return query, args
}

func appendFilter(query string, args []any, field, value string) (string, []any) {
    if value == "" {
        return query, args
    }
    query += ` AND ` + field + `=$` + strconv.Itoa(len(args)+1)
    args = append(args, value)
    return query, args
}

func appendPaging(query string, args []any, page, perPage int) (string, []any) {
    offset := (page - 1) * perPage
    query += ` ORDER BY date DESC, created_at DESC`
    query += ` LIMIT $` + strconv.Itoa(len(args)+1)
    query += ` OFFSET $` + strconv.Itoa(len(args)+2)
    args = append(args, perPage, offset)
    return query, args
}

func scanList(rows pgx.Rows) ([]ListItem, error) {
    items := make([]ListItem, 0)
    for rows.Next() {
        item := ListItem{}
        err := rows.Scan(&item.ID, &item.Type, &item.Date, &item.Number, &item.Status, &item.Note)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func (r *Repository) ByID(ctx context.Context, tenantID, id string) (Document, error) {
    query := `SELECT id::text, type, date::text, COALESCE(number,''), status,
        COALESCE(warehouse_id::text,''), COALESCE(source_warehouse_id::text,''),
        COALESCE(target_warehouse_id::text,''), COALESCE(note,'')
        FROM documents.documents
        WHERE tenant_id=$1 AND id=$2`
    row := r.store.Pool.QueryRow(ctx, query, tenantID, id)
    item := Document{}
    err := row.Scan(&item.ID, &item.Type, &item.Date, &item.Number, &item.Status,
        &item.WarehouseID, &item.SourceWarehouseID, &item.TargetWarehouseID, &item.Note)
    return item, err
}

func (r *Repository) LoadItemsWithProfit(ctx context.Context, tenantID, documentID string) ([]DocumentItem, float64, error) {
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
        return nil, 0, err
    }
    defer rows.Close()
    return scanItems(rows)
}

func scanItems(rows pgx.Rows) ([]DocumentItem, float64, error) {
    items := make([]DocumentItem, 0)
    total := 0.0
    for rows.Next() {
        item := DocumentItem{}
        err := rows.Scan(&item.ID, &item.VariantID, &item.Qty, &item.UnitPrice, &item.TotalAmount, &item.Profit)
        if err != nil {
            return nil, 0, err
        }
        total += item.Profit
        items = append(items, item)
    }
    return items, total, rows.Err()
}

func (r *Repository) InsertDocument(ctx context.Context, tx pgx.Tx, tenantID, documentID string, req CreateRequest) error {
    query := `INSERT INTO documents.documents
        (id, tenant_id, type, date, number, status, warehouse_id, source_warehouse_id, target_warehouse_id, note)
        VALUES ($1,$2,$3,$4,$5,'draft',NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),$9)`
    _, err := tx.Exec(ctx, query, documentID, tenantID, req.Type, req.Date,
        req.Number, req.WarehouseID, req.SourceWarehouseID, req.TargetWarehouseID, req.Note)
    return err
}

func (r *Repository) UpdateDocument(ctx context.Context, tx pgx.Tx, tenantID, documentID string, req UpdateRequest) error {
    query := `UPDATE documents.documents
        SET type=$3, date=$4, number=$5, warehouse_id=NULLIF($6,''),
            source_warehouse_id=NULLIF($7,''), target_warehouse_id=NULLIF($8,''),
            note=$9, updated_at=now()
        WHERE tenant_id=$1 AND id=$2`
    _, err := tx.Exec(ctx, query, tenantID, documentID, req.Type, req.Date, req.Number,
        req.WarehouseID, req.SourceWarehouseID, req.TargetWarehouseID, req.Note)
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
    total := item.Qty * item.UnitPrice
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
