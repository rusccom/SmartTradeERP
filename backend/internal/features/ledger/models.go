package ledger

import (
	"time"

	"github.com/shopspring/decimal"
)

type EntryInput struct {
	TenantID       string
	VariantID      string
	DocumentID     string
	DocumentItemID string
	WarehouseID    string
	Date           time.Time
	Type           string
	Reason         string
	Qty            decimal.Decimal
	UnitPrice      decimal.Decimal
	TotalAmount    decimal.Decimal
	Revenue        *decimal.Decimal
}

type Movement struct {
	SequenceNum int64           `json:"sequence_num"`
	Date        string          `json:"date"`
	Type        string          `json:"type"`
	Reason      string          `json:"reason"`
	Qty         decimal.Decimal `json:"qty"`
	RunningQty  decimal.Decimal `json:"running_qty"`
	RunningAvg  decimal.Decimal `json:"running_avg"`
	COGS        decimal.Decimal `json:"cogs"`
	Profit      decimal.Decimal `json:"profit"`
}
