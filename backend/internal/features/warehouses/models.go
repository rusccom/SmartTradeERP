package warehouses

type Warehouse struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Address   string `json:"address"`
    IsDefault bool   `json:"is_default"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
}

type CreateRequest struct {
    Name      string `json:"name"`
    Address   string `json:"address"`
    IsDefault bool   `json:"is_default"`
    IsActive  bool   `json:"is_active"`
}

type UpdateRequest struct {
    Name      string `json:"name"`
    Address   string `json:"address"`
    IsDefault bool   `json:"is_default"`
    IsActive  bool   `json:"is_active"`
}
