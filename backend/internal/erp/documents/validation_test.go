package documents

import "testing"

func TestValidateDocument_ValidSale(t *testing.T) {
	if err := validateDocument(saleRequest()); err != nil {
		t.Fatalf("valid sale rejected: %v", err)
	}
}

func TestValidateDocument_Rejects(t *testing.T) {
	cases := map[string]func(CreateRequest) CreateRequest{
		"unknown type":      func(r CreateRequest) CreateRequest { r.Type = "BOGUS"; return r },
		"bad date":          func(r CreateRequest) CreateRequest { r.Date = "2026-13-40"; return r },
		"sale no warehouse": func(r CreateRequest) CreateRequest { r.WarehouseID = ""; return r },
		"sale with source":  func(r CreateRequest) CreateRequest { r.SourceWarehouseID = uuidWarehouse2; return r },
		"shift on non-sale": func(r CreateRequest) CreateRequest { r.Type = "RECEIPT"; r.ShiftID = uuidShift; return r },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			if err := validateDocument(mutate(saleRequest())); err == nil {
				t.Fatalf("expected rejection for %q", name)
			}
		})
	}
}

func TestValidateDocument_Transfer(t *testing.T) {
	base := func() CreateRequest {
		return CreateRequest{
			Type: "TRANSFER", Date: "2026-06-22",
			SourceWarehouseID: uuidWarehouse, TargetWarehouseID: uuidWarehouse2,
			Items: []ItemInput{{VariantID: uuidVariant, Qty: dec("1"), UnitPrice: dec("0")}},
		}
	}
	if err := validateDocument(base()); err != nil {
		t.Fatalf("valid transfer rejected: %v", err)
	}
	withWarehouse := base()
	withWarehouse.WarehouseID = uuidWarehouse
	if err := validateDocument(withWarehouse); err == nil {
		t.Fatal("transfer must not carry warehouse_id")
	}
	sameWarehouse := base()
	sameWarehouse.TargetWarehouseID = uuidWarehouse
	if err := validateDocument(sameWarehouse); err == nil {
		t.Fatal("transfer source and target must differ")
	}
	missingTarget := base()
	missingTarget.TargetWarehouseID = ""
	if err := validateDocument(missingTarget); err == nil {
		t.Fatal("transfer requires both warehouses")
	}
}

func TestValidateInventoryItems(t *testing.T) {
	dup := []ItemInput{
		{VariantID: uuidVariant, Qty: dec("10"), UnitPrice: dec("0")},
		{VariantID: uuidVariant, Qty: dec("3"), UnitPrice: dec("0")},
	}
	if err := validateInventoryItems("INVENTORY", dup); err == nil {
		t.Fatal("inventory must reject duplicate variant rows")
	}
	if err := validateInventoryItems("SALE", dup); err != nil {
		t.Fatalf("sale may repeat a variant: %v", err)
	}
	unique := []ItemInput{
		{VariantID: uuidVariant, Qty: dec("10"), UnitPrice: dec("0")},
		{VariantID: uuidVariant2, Qty: dec("3"), UnitPrice: dec("0")},
	}
	if err := validateInventoryItems("INVENTORY", unique); err != nil {
		t.Fatalf("unique inventory rows rejected: %v", err)
	}
}

func TestInvalidItem_QtyRules(t *testing.T) {
	zero := ItemInput{VariantID: uuidVariant, Qty: dec("0"), UnitPrice: dec("10")}
	if invalidItem("INVENTORY", zero) {
		t.Fatal("inventory count may be zero")
	}
	if !invalidItem("SALE", zero) {
		t.Fatal("sale qty must be positive")
	}
	negPrice := ItemInput{VariantID: uuidVariant, Qty: dec("1"), UnitPrice: dec("-1")}
	if !invalidItem("SALE", negPrice) {
		t.Fatal("negative price must be rejected")
	}
	badVariant := ItemInput{VariantID: "not-a-uuid", Qty: dec("1"), UnitPrice: dec("1")}
	if !invalidItem("SALE", badVariant) {
		t.Fatal("non-uuid variant must be rejected")
	}
}

func TestIsDocumentType(t *testing.T) {
	for _, ok := range []string{"RECEIPT", "SALE", "WRITEOFF", "INVENTORY", "TRANSFER", "RETURN"} {
		if !isDocumentType(ok) {
			t.Fatalf("%q should be a valid type", ok)
		}
	}
	if isDocumentType("PAYMENT") {
		t.Fatal("PAYMENT is not a document type")
	}
}

func TestAllowsBundleDocument(t *testing.T) {
	if !allowsBundleDocument("SALE") || !allowsBundleDocument("RETURN") {
		t.Fatal("bundles allowed for sale and return")
	}
	if allowsBundleDocument("RECEIPT") {
		t.Fatal("bundles not allowed for receipt")
	}
}
