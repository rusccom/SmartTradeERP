package admin

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type Tenant struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Status    string `json:"status"`
    Plan      string `json:"plan"`
    CreatedAt string `json:"created_at"`
}

type Stats struct {
    TotalTenants  int `json:"total_tenants"`
    ActiveTenants int `json:"active_tenants"`
    NewLast30Days int `json:"new_last_30_days"`
}

type adminUser struct {
    ID           string
    PasswordHash string
}
