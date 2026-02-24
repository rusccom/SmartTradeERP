package ledger

import "time"

type EntryInput struct {
    TenantID       string
    VariantID      string
    DocumentID     string
    DocumentItemID string
    WarehouseID    string
    Date           time.Time
    Type           string
    Reason         string
    Qty            float64
    UnitPrice      float64
    TotalAmount    float64
    Revenue        *float64
}

type Movement struct {
    SequenceNum int64   `json:"sequence_num"`
    Date        string  `json:"date"`
    Type        string  `json:"type"`
    Reason      string  `json:"reason"`
    Qty         float64 `json:"qty"`
    RunningQty  float64 `json:"running_qty"`
    RunningAvg  float64 `json:"running_avg"`
    COGS        float64 `json:"cogs"`
    Profit      float64 `json:"profit"`
}
