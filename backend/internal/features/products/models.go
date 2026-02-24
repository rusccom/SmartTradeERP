package products

type Product struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    IsComposite bool   `json:"is_composite"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type CreateRequest struct {
    Name        string  `json:"name"`
    IsComposite bool    `json:"is_composite"`
    Unit        string  `json:"unit"`
    Price       float64 `json:"price"`
    SKUCode     string  `json:"sku_code"`
    Barcode     string  `json:"barcode"`
}

type UpdateRequest struct {
    Name        string `json:"name"`
    IsComposite bool   `json:"is_composite"`
}
