package variants

type Variant struct {
    ID        string  `json:"id"`
    ProductID string  `json:"product_id"`
    Name      string  `json:"name"`
    SKUCode   string  `json:"sku_code"`
    Barcode   string  `json:"barcode"`
    Unit      string  `json:"unit"`
    Price     float64 `json:"price"`
    Option1   string  `json:"option1"`
    Option2   string  `json:"option2"`
    Option3   string  `json:"option3"`
}

type CreateRequest struct {
    ProductID string  `json:"product_id"`
    Name      string  `json:"name"`
    SKUCode   string  `json:"sku_code"`
    Barcode   string  `json:"barcode"`
    Unit      string  `json:"unit"`
    Price     float64 `json:"price"`
    Option1   string  `json:"option1"`
    Option2   string  `json:"option2"`
    Option3   string  `json:"option3"`
}

type UpdateRequest struct {
    Name    string  `json:"name"`
    SKUCode string  `json:"sku_code"`
    Barcode string  `json:"barcode"`
    Unit    string  `json:"unit"`
    Price   float64 `json:"price"`
    Option1 string  `json:"option1"`
    Option2 string  `json:"option2"`
    Option3 string  `json:"option3"`
}

type Component struct {
    ComponentVariantID string  `json:"component_variant_id"`
    Qty                float64 `json:"qty"`
}

type StockByWarehouse struct {
    WarehouseID string  `json:"warehouse_id"`
    Warehouse   string  `json:"warehouse"`
    Qty         float64 `json:"qty"`
}

type Stock struct {
    GlobalQty       float64            `json:"global_qty"`
    RunningAvg      float64            `json:"running_avg"`
    WarehousesStock []StockByWarehouse `json:"warehouses"`
}
