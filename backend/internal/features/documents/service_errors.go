package documents

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func mapDocumentWriteError(err error) error {
	if isDocumentNumberConflict(err) {
		return ErrDocumentNumberConflict
	}
	return err
}

func isDocumentNumberConflict(err error) bool {
	pgErr := &pgconn.PgError{}
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "23505"
}
