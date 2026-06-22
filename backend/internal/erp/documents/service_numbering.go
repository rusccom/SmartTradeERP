package documents

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *Service) withDocumentNumber(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	req CreateRequest,
) (CreateRequest, error) {
	if req.Number != "" {
		return req, s.syncManualNumber(ctx, tx, tenantID, req)
	}
	number, err := s.repo.NextDocumentNumber(ctx, tx, documentNumberInput{
		TenantID: tenantID,
		Type:     req.Type,
		Year:     documentYear(req.Date),
	})
	if err != nil {
		return req, err
	}
	req.Number = number
	return req, nil
}

func (s *Service) syncManualNumber(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	req CreateRequest,
) error {
	year, sequence, ok := parseGeneratedNumber(req.Type, req.Number)
	if !ok {
		return nil
	}
	input := documentNumberInput{TenantID: tenantID, Type: req.Type, Year: year}
	return s.repo.EnsureDocumentSequenceAtLeast(ctx, tx, input, sequence+1)
}

func (s *Service) withUpdateNumber(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	id string,
	req UpdateRequest,
) (UpdateRequest, error) {
	if req.Number != "" {
		return req, nil
	}
	number, err := s.repo.DocumentNumber(ctx, tx, tenantID, id)
	if err != nil || number != "" {
		req.Number = number
		return req, err
	}
	return s.withDocumentNumber(ctx, tx, tenantID, req)
}

func documentYear(raw string) int {
	value, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Now().Year()
	}
	return value.Year()
}

func parseGeneratedNumber(documentType, number string) (int, int, bool) {
	parts := strings.Split(number, "-")
	if !hasGeneratedNumberShape(documentType, parts) {
		return 0, 0, false
	}
	year, yearErr := strconv.Atoi(parts[1])
	sequence, sequenceErr := strconv.Atoi(parts[2])
	if yearErr != nil || sequenceErr != nil || year < 2000 || sequence <= 0 {
		return 0, 0, false
	}
	return year, sequence, true
}

func hasGeneratedNumberShape(documentType string, parts []string) bool {
	return len(parts) == 3 && parts[0] == documentType &&
		len(parts[1]) == 4 && len(parts[2]) == 6
}
