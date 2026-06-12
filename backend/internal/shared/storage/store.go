package storage

import (
	"context"
	"errors"
	"time"
)

var ErrObjectNotFound = errors.New("object not found")

type ObjectInfo struct {
	Key         string
	ContentType string
	SizeBytes   int64
}

type PresignPutRequest struct {
	Key         string
	ContentType string
	Expires     time.Duration
}

type PresignedPut struct {
	URL       string
	Method    string
	Headers   map[string]string
	ExpiresAt time.Time
}

type ObjectStore interface {
	Delete(ctx context.Context, key string) error
	Head(ctx context.Context, key string) (ObjectInfo, error)
	PresignPut(request PresignPutRequest) (PresignedPut, error)
	PublicURL(key string) string
}

type R2Options struct {
	AccountID     string
	Bucket        string
	AccessKeyID   string
	SecretKey     string
	PublicBaseURL string
}

func (o R2Options) Enabled() bool {
	return o.AccountID != "" && o.Bucket != "" && o.AccessKeyID != "" &&
		o.SecretKey != "" && o.PublicBaseURL != ""
}

func (o R2Options) PartiallyConfigured() bool {
	if o.Enabled() {
		return false
	}
	return o.AccountID != "" || o.Bucket != "" || o.AccessKeyID != "" ||
		o.SecretKey != "" || o.PublicBaseURL != ""
}
