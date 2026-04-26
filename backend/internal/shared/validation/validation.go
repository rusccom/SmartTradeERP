package validation

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var ErrInvalidData = errors.New("invalid data")

func Clean(value string) string {
	return strings.TrimSpace(value)
}

func Required(value string) bool {
	return Clean(value) != ""
}

func Max(value string, size int) bool {
	return len([]rune(Clean(value))) <= size
}

func Email(value string) bool {
	cleaned := Clean(value)
	if cleaned == "" {
		return true
	}
	_, err := mail.ParseAddress(cleaned)
	return err == nil
}

func UUID(value string) bool {
	_, err := uuid.Parse(Clean(value))
	return err == nil
}

func NonNegative(value decimal.Decimal) bool {
	return value.GreaterThanOrEqual(decimal.Zero)
}

func Positive(value decimal.Decimal) bool {
	return value.GreaterThan(decimal.Zero)
}
