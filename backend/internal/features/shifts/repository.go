package shifts

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/shared/db"
)

type Repository struct {
	store *db.Store
}

type shiftTotals struct {
	totalSales      decimal.Decimal
	totalReturns    decimal.Decimal
	salesCash       decimal.Decimal
	salesCard       decimal.Decimal
	salesTransfer   decimal.Decimal
	returnsCash     decimal.Decimal
	returnsCard     decimal.Decimal
	returnsTransfer decimal.Decimal
	totalCashIn     decimal.Decimal
	totalCashOut    decimal.Decimal
}

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) Insert(ctx context.Context, tx pgx.Tx, tenantID string, shift Shift) error {
	query := `INSERT INTO documents.shifts
        (id, tenant_id, user_id, warehouse_id, opening_cash, status)
        VALUES ($1,$2,$3,$4,$5,'open')`
	_, err := tx.Exec(ctx, query, shift.ID, tenantID, shift.UserID, shift.WarehouseID, shift.OpeningCash)
	return err
}

func (r *Repository) WarehouseExists(ctx context.Context, tenantID string, id string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.warehouses WHERE tenant_id=$1 AND id=$2
    )`
	row := r.store.Pool.QueryRow(ctx, query, tenantID, id)
	exists := false
	err := row.Scan(&exists)
	return exists, err
}

func (r *Repository) FindOpen(ctx context.Context, tenantID, userID string) (Shift, error) {
	query := `SELECT id::text, user_id::text, warehouse_id::text, opened_at::text,
        COALESCE(closed_at::text,''), opening_cash, COALESCE(closing_cash,0), status
        FROM documents.shifts
        WHERE tenant_id=$1 AND user_id=$2 AND status='open'
        ORDER BY opened_at DESC
        LIMIT 1`
	row := r.store.Pool.QueryRow(ctx, query, tenantID, userID)
	return scanShiftRow(row)
}

func (r *Repository) Close(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	shiftID string,
	closingCash decimal.Decimal,
) error {
	query := `UPDATE documents.shifts
        SET status='closed', closed_at=now(), closing_cash=$3, updated_at=now()
        WHERE tenant_id=$1 AND id=$2 AND status='open'`
	tag, err := tx.Exec(ctx, query, tenantID, shiftID, closingCash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrShiftAlreadyClosed
	}
	return nil
}

func (r *Repository) ByID(ctx context.Context, tenantID, shiftID string) (Shift, error) {
	query := `SELECT id::text, user_id::text, warehouse_id::text, opened_at::text,
        COALESCE(closed_at::text,''), opening_cash, COALESCE(closing_cash,0), status
        FROM documents.shifts
        WHERE tenant_id=$1 AND id=$2`
	row := r.store.Pool.QueryRow(ctx, query, tenantID, shiftID)
	return scanShiftRow(row)
}

func scanShiftRow(row pgx.Row) (Shift, error) {
	item := Shift{}
	err := row.Scan(&item.ID, &item.UserID, &item.WarehouseID, &item.OpenedAt,
		&item.ClosedAt, &item.OpeningCash, &item.ClosingCash, &item.Status)
	return item, err
}

func (r *Repository) InsertCashOp(ctx context.Context, tx pgx.Tx, shiftID string, op CashOpRequest) error {
	query := `INSERT INTO documents.shift_cash_ops (id, shift_id, type, amount, note)
        VALUES ($1,$2,$3,$4,$5)`
	_, err := tx.Exec(ctx, query, uuid.NewString(), shiftID, op.Type, op.Amount, op.Note)
	return err
}

func (r *Repository) CashOps(ctx context.Context, shiftID string) ([]CashOp, error) {
	query := `SELECT id::text, type, amount, COALESCE(note,''), created_at::text
        FROM documents.shift_cash_ops
        WHERE shift_id=$1
        ORDER BY created_at`
	rows, err := r.store.Pool.Query(ctx, query, shiftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCashOps(rows)
}

func scanCashOps(rows pgx.Rows) ([]CashOp, error) {
	items := make([]CashOp, 0)
	for rows.Next() {
		item := CashOp{}
		if err := rows.Scan(&item.ID, &item.Type, &item.Amount, &item.Note, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) ShiftTotals(ctx context.Context, q db.DBTX, tenantID, shiftID string) (shiftTotals, error) {
	totals := shiftTotals{}
	if err := loadDocumentTotals(ctx, q, tenantID, shiftID, &totals); err != nil {
		return shiftTotals{}, err
	}
	if err := loadPaymentTotals(ctx, q, tenantID, shiftID, &totals); err != nil {
		return shiftTotals{}, err
	}
	if err := loadCashTotals(ctx, q, shiftID, &totals); err != nil {
		return shiftTotals{}, err
	}
	return totals, nil
}

func (r *Repository) LockOpenShift(ctx context.Context, tx pgx.Tx, tenantID, shiftID string) error {
	query := `SELECT 1 FROM documents.shifts
        WHERE tenant_id=$1 AND id=$2 AND status='open'
        FOR UPDATE`
	value := 0
	err := tx.QueryRow(ctx, query, tenantID, shiftID).Scan(&value)
	if err == pgx.ErrNoRows {
		return ErrShiftAlreadyClosed
	}
	return err
}

func loadDocumentTotals(
	ctx context.Context,
	q db.DBTX,
	tenantID string,
	shiftID string,
	totals *shiftTotals,
) error {
	query := `WITH doc_sums AS (
        SELECT d.type, SUM(i.total_amount) AS amount
        FROM documents.documents d
        JOIN documents.document_items i ON i.document_id=d.id
        WHERE d.tenant_id=$1 AND d.shift_id=$2 AND d.status='posted'
        GROUP BY d.type
    )
    SELECT
        COALESCE(SUM(CASE WHEN type='SALE' THEN amount ELSE 0 END),0),
        COALESCE(SUM(CASE WHEN type='RETURN' THEN amount ELSE 0 END),0)
    FROM doc_sums`
	row := q.QueryRow(ctx, query, tenantID, shiftID)
	return row.Scan(&totals.totalSales, &totals.totalReturns)
}

func loadPaymentTotals(
	ctx context.Context,
	q db.DBTX,
	tenantID string,
	shiftID string,
	totals *shiftTotals,
) error {
	query := `SELECT
        COALESCE(SUM(CASE WHEN d.type='SALE' AND p.method='cash' THEN p.amount ELSE 0 END),0),
        COALESCE(SUM(CASE WHEN d.type='SALE' AND p.method='card' THEN p.amount ELSE 0 END),0),
        COALESCE(SUM(CASE WHEN d.type='SALE' AND p.method='transfer' THEN p.amount ELSE 0 END),0),
        COALESCE(SUM(CASE WHEN d.type='RETURN' AND p.method='cash' THEN p.amount ELSE 0 END),0),
        COALESCE(SUM(CASE WHEN d.type='RETURN' AND p.method='card' THEN p.amount ELSE 0 END),0),
        COALESCE(SUM(CASE WHEN d.type='RETURN' AND p.method='transfer' THEN p.amount ELSE 0 END),0)
        FROM documents.documents d
        JOIN documents.document_payments p ON p.document_id=d.id
        WHERE d.tenant_id=$1 AND d.shift_id=$2 AND d.status='posted'`
	row := q.QueryRow(ctx, query, tenantID, shiftID)
	return row.Scan(&totals.salesCash, &totals.salesCard, &totals.salesTransfer,
		&totals.returnsCash, &totals.returnsCard, &totals.returnsTransfer)
}

func loadCashTotals(ctx context.Context, q db.DBTX, shiftID string, totals *shiftTotals) error {
	query := `SELECT
        COALESCE(SUM(CASE WHEN type='cash_in' THEN amount ELSE 0 END),0),
        COALESCE(SUM(CASE WHEN type='cash_out' THEN amount ELSE 0 END),0)
        FROM documents.shift_cash_ops
        WHERE shift_id=$1`
	row := q.QueryRow(ctx, query, shiftID)
	return row.Scan(&totals.totalCashIn, &totals.totalCashOut)
}
