package currencies

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"smarterp/backend/internal/shared/validation"
)

var ErrCurrencyExists = errors.New("currency exists")

func mapCurrencyWriteError(err error) error {
	if isPgError(err, "23505") {
		return ErrCurrencyExists
	}
	if isPgError(err, "23503") || isPgError(err, "23514") {
		return validation.ErrInvalidData
	}
	return err
}

func isPgError(err error, code string) bool {
	pgErr := &pgconn.PgError{}
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == code
}
