package reports

type ProfitReport struct {
    Profit float64 `json:"profit"`
}

type StockRow struct {
    VariantID string  `json:"variant_id"`
    Name      string  `json:"name"`
    Qty       float64 `json:"qty"`
    Avg       float64 `json:"avg"`
}

type TopProduct struct {
    ProductID string  `json:"product_id"`
    Name      string  `json:"name"`
    Profit    float64 `json:"profit"`
}
