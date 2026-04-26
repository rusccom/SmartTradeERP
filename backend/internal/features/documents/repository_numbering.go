package documents

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type documentNumberInput struct {
	TenantID string
	Type     string
	Year     int
}

func (r *Repository) NextDocumentNumber(
	ctx context.Context,
	tx pgx.Tx,
	input documentNumberInput,
) (string, error) {
	sequence, err := r.reserveDocumentSequence(ctx, tx, input)
	if err != nil {
		return "", err
	}
	return formatDocumentNumber(input, sequence), nil
}

func (r *Repository) EnsureDocumentSequenceAtLeast(
	ctx context.Context,
	tx pgx.Tx,
	input documentNumberInput,
	nextNumber int,
) error {
	query := `INSERT INTO documents.document_sequences
        (tenant_id, document_type, year, next_number)
        VALUES ($1,$2,$3,$4)
        ON CONFLICT (tenant_id, document_type, year)
        DO UPDATE SET next_number = GREATEST(documents.document_sequences.next_number, EXCLUDED.next_number),
            updated_at = now()`
	_, err := tx.Exec(ctx, query, input.TenantID, input.Type, input.Year, nextNumber)
	return err
}

func (r *Repository) reserveDocumentSequence(
	ctx context.Context,
	tx pgx.Tx,
	input documentNumberInput,
) (int, error) {
	query := `INSERT INTO documents.document_sequences
        (tenant_id, document_type, year, next_number)
        VALUES ($1,$2,$3,2)
        ON CONFLICT (tenant_id, document_type, year)
        DO UPDATE SET next_number = documents.document_sequences.next_number + 1,
            updated_at = now()
        RETURNING next_number - 1`
	row := tx.QueryRow(ctx, query, input.TenantID, input.Type, input.Year)
	sequence := 0
	err := row.Scan(&sequence)
	return sequence, err
}

func formatDocumentNumber(input documentNumberInput, sequence int) string {
	return fmt.Sprintf("%s-%04d-%06d", input.Type, input.Year, sequence)
}

func (r *Repository) DocumentNumber(ctx context.Context, tx pgx.Tx, tenantID, id string) (string, error) {
	query := `SELECT COALESCE(number,'')
        FROM documents.documents
        WHERE tenant_id=$1 AND id=$2`
	row := tx.QueryRow(ctx, query, tenantID, id)
	number := ""
	err := row.Scan(&number)
	return number, err
}
