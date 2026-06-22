package documents

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) ReplacePayments(
	ctx context.Context,
	tx pgx.Tx,
	documentID string,
	payments []PaymentInput,
) error {
	if err := r.deletePayments(ctx, tx, documentID); err != nil {
		return err
	}
	for _, payment := range payments {
		if err := r.insertPayment(ctx, tx, documentID, payment); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) deletePayments(ctx context.Context, tx pgx.Tx, documentID string) error {
	_, err := tx.Exec(ctx, `DELETE FROM documents.document_payments WHERE document_id=$1`, documentID)
	return err
}

func (r *Repository) insertPayment(
	ctx context.Context,
	tx pgx.Tx,
	documentID string,
	payment PaymentInput,
) error {
	query := `INSERT INTO documents.document_payments (id, document_id, method, amount)
        VALUES ($1,$2,$3,$4)`
	_, err := tx.Exec(ctx, query, uuid.NewString(), documentID, payment.Method, payment.Amount)
	return err
}

func (r *Repository) LoadPayments(ctx context.Context, documentID string) ([]Payment, error) {
	query := `SELECT id::text, method, amount
        FROM documents.document_payments
        WHERE document_id=$1
        ORDER BY id`
	rows, err := r.store.Pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPayments(rows)
}

func scanPayments(rows pgx.Rows) ([]Payment, error) {
	items := make([]Payment, 0)
	for rows.Next() {
		item := Payment{}
		if err := rows.Scan(&item.ID, &item.Method, &item.Amount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
