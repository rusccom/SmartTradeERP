package documents

type ItemInput struct {
    VariantID  string  `json:"variant_id"`
    Qty        float64 `json:"qty"`
    UnitPrice  float64 `json:"unit_price"`
    Direction  string  `json:"direction"`
}

type CreateRequest struct {
    Type              string      `json:"type"`
    Date              string      `json:"date"`
    Number            string      `json:"number"`
    WarehouseID       string      `json:"warehouse_id"`
    SourceWarehouseID string      `json:"source_warehouse_id"`
    TargetWarehouseID string      `json:"target_warehouse_id"`
    Note              string      `json:"note"`
    Items             []ItemInput `json:"items"`
}

type UpdateRequest = CreateRequest

type Document struct {
    ID                string         `json:"id"`
    Type              string         `json:"type"`
    Date              string         `json:"date"`
    Number            string         `json:"number"`
    Status            string         `json:"status"`
    WarehouseID       string         `json:"warehouse_id"`
    SourceWarehouseID string         `json:"source_warehouse_id"`
    TargetWarehouseID string         `json:"target_warehouse_id"`
    Note              string         `json:"note"`
    Items             []DocumentItem `json:"items"`
    TotalProfit       float64        `json:"total_profit"`
}

type DocumentItem struct {
    ID          string  `json:"id"`
    VariantID   string  `json:"variant_id"`
    Qty         float64 `json:"qty"`
    UnitPrice   float64 `json:"unit_price"`
    TotalAmount float64 `json:"total_amount"`
    Profit      float64 `json:"profit"`
}

type ListItem struct {
    ID      string `json:"id"`
    Type    string `json:"type"`
    Date    string `json:"date"`
    Number  string `json:"number"`
    Status  string `json:"status"`
    Note    string `json:"note"`
}

type Filters struct {
    Type   string
    Status string
    Date   string
}

type postingItem struct {
    ID          string
    VariantID   string
    Qty         float64
    UnitPrice   float64
    TotalAmount float64
    IsComposite bool
}

type variantComponent struct {
    ComponentVariantID string
    QtyPerUnit         float64
}
