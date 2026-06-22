package clientauth

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    TenantName string `json:"tenant_name"`
    Email      string `json:"email"`
    Password   string `json:"password"`
}

type userRecord struct {
    ID           string
    TenantID     string
    Role         string
    PasswordHash string
}
