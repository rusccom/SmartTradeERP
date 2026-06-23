package media

type Item struct {
	ID          string `json:"id"`
	OwnerType   string `json:"owner_type"`
	OwnerID     string `json:"owner_id"`
	ObjectKey   string `json:"-"`
	URL         string `json:"url"`
	ThumbURL    string `json:"thumb_url"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`
	IsPrimary   bool   `json:"is_primary"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type PrimaryRef struct {
	OwnerID   string
	ObjectKey string
}

type DirectUploadInput struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`
}

type DirectUploadRequest struct {
	TenantID  string
	OwnerType string
	OwnerID   string
	Input     DirectUploadInput
}

type DirectUpload struct {
	ID        string            `json:"id"`
	UploadURL string            `json:"upload_url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ExpiresAt string            `json:"expires_at"`
}

type CompleteRequest struct {
	TenantID  string
	OwnerType string
	OwnerID   string
	MediaID   string
}
