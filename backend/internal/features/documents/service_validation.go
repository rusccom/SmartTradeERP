package documents

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/validation"
)

func normalizeRequest(req CreateRequest) CreateRequest {
	req.Type = validation.Clean(req.Type)
	req.Date = validation.Clean(req.Date)
	req.Number = validation.Clean(req.Number)
	req.WarehouseID = validation.Clean(req.WarehouseID)
	req.SourceWarehouseID = validation.Clean(req.SourceWarehouseID)
	req.TargetWarehouseID = validation.Clean(req.TargetWarehouseID)
	req.ShiftID = validation.Clean(req.ShiftID)
	req.CustomerID = validation.Clean(req.CustomerID)
	req.Note = validation.Clean(req.Note)
	req.Items = normalizeItems(req.Items)
	req.Payments = normalizePayments(req.Payments)
	return req
}

func normalizeItems(items []ItemInput) []ItemInput {
	for index := range items {
		items[index].VariantID = validation.Clean(items[index].VariantID)
	}
	return items
}

func normalizePayments(payments []PaymentInput) []PaymentInput {
	for index := range payments {
		payments[index].Method = validation.Clean(payments[index].Method)
	}
	return payments
}

func (s *Service) validateRequest(req CreateRequest) error {
	if err := validateDocument(req); err != nil {
		return err
	}
	if err := validateItems(req.Items); err != nil {
		return err
	}
	return validatePayments(req.Type, req.Items, req.Payments)
}

func validateDocument(req CreateRequest) error {
	if !isDocumentType(req.Type) || !isDate(req.Date) {
		return validation.ErrInvalidData
	}
	if !validation.Max(req.Number, 64) || !validation.Max(req.Note, 1000) {
		return validation.ErrInvalidData
	}
	if err := validateDocumentIDs(req); err != nil {
		return err
	}
	return validateWarehouseRules(req)
}

func validateDocumentIDs(req CreateRequest) error {
	values := append(warehouseIDs(req), req.CustomerID, req.ShiftID)
	for _, value := range values {
		if value != "" && !validation.UUID(value) {
			return validation.ErrInvalidData
		}
	}
	return nil
}

func validateWarehouseRules(req CreateRequest) error {
	if req.Type == "TRANSFER" {
		if req.WarehouseID != "" {
			return validation.ErrInvalidData
		}
		return validateTransferWarehouses(req)
	}
	if req.SourceWarehouseID != "" || req.TargetWarehouseID != "" {
		return validation.ErrInvalidData
	}
	if !requiresWarehouse(req.Type) {
		return nil
	}
	if !validation.Required(req.WarehouseID) {
		return validation.ErrInvalidData
	}
	return nil
}

func validateTransferWarehouses(req CreateRequest) error {
	source := req.SourceWarehouseID
	target := req.TargetWarehouseID
	if !validation.Required(source) || !validation.Required(target) {
		return validation.ErrInvalidData
	}
	if source == target {
		return validation.ErrInvalidData
	}
	return nil
}

func validateItems(items []ItemInput) error {
	if len(items) == 0 {
		return validation.ErrInvalidData
	}
	for _, item := range items {
		if invalidItem(item) {
			return validation.ErrInvalidData
		}
	}
	return nil
}

func invalidItem(item ItemInput) bool {
	if !validation.Required(item.VariantID) || !validation.UUID(item.VariantID) {
		return true
	}
	return !validation.Positive(item.Qty) || !validation.NonNegative(item.UnitPrice)
}

func (s *Service) validateReferences(ctx context.Context, tx pgx.Tx, tenantID string, req CreateRequest) error {
	if err := s.validateWarehouses(ctx, tx, tenantID, req); err != nil {
		return err
	}
	if err := s.validateCustomer(ctx, tx, tenantID, req.CustomerID); err != nil {
		return err
	}
	if err := s.validateShift(ctx, tx, tenantID, req.ShiftID, req.Type); err != nil {
		return err
	}
	return s.validateVariants(ctx, tx, tenantID, req.Type, req.Items)
}

func (s *Service) validateWarehouses(ctx context.Context, tx pgx.Tx, tenantID string, req CreateRequest) error {
	for _, id := range warehouseIDs(req) {
		exists, err := s.repo.WarehouseExists(ctx, tx, tenantID, id)
		if err != nil || !exists {
			return referenceError(err)
		}
	}
	return nil
}

func (s *Service) validateCustomer(ctx context.Context, tx pgx.Tx, tenantID string, id string) error {
	if id == "" {
		return nil
	}
	exists, err := s.repo.CustomerExists(ctx, tx, tenantID, id)
	if err != nil || !exists {
		return referenceError(err)
	}
	return nil
}

func (s *Service) validateShift(ctx context.Context, tx pgx.Tx, tenantID string, id string, documentType string) error {
	if id == "" {
		return nil
	}
	exists, err := s.shiftExists(ctx, tx, tenantID, id, documentType)
	if err != nil || !exists {
		return referenceError(err)
	}
	return nil
}

func (s *Service) shiftExists(ctx context.Context, tx pgx.Tx, tenantID, id, documentType string) (bool, error) {
	if requiresOpenShift(documentType) {
		return s.repo.OpenShiftExists(ctx, tx, tenantID, id)
	}
	return s.repo.ShiftExists(ctx, tx, tenantID, id)
}

func requiresOpenShift(documentType string) bool {
	return documentType == "SALE" || documentType == "RETURN"
}

func (s *Service) validateVariants(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	documentType string,
	items []ItemInput,
) error {
	for _, id := range itemVariantIDs(items) {
		if err := s.validateVariantReference(ctx, tx, tenantID, documentType, id); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) validateVariantReference(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	documentType string,
	id string,
) error {
	composite, exists, err := s.repo.VariantComposite(ctx, tx, tenantID, id)
	if err != nil || !exists {
		return referenceError(err)
	}
	if composite && !allowsBundleDocument(documentType) {
		return ErrBundleDocumentType
	}
	return nil
}

func allowsBundleDocument(documentType string) bool {
	return documentType == "SALE" || documentType == "RETURN"
}

func referenceError(err error) error {
	if err != nil {
		return err
	}
	return ErrInvalidDocumentReference
}

func warehouseIDs(req CreateRequest) []string {
	return uniqueNonEmpty([]string{req.WarehouseID, req.SourceWarehouseID, req.TargetWarehouseID})
}

func itemVariantIDs(items []ItemInput) []string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.VariantID)
	}
	return uniqueNonEmpty(values)
}

func uniqueNonEmpty(values []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func isDocumentType(value string) bool {
	switch value {
	case "RECEIPT", "SALE", "WRITEOFF", "INVENTORY", "TRANSFER", "RETURN":
		return true
	default:
		return false
	}
}

func requiresWarehouse(value string) bool {
	return value != "" && value != "TRANSFER"
}

func isDate(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}
