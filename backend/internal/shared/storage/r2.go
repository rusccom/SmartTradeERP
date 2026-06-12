package storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const r2Region = "auto"
const r2Service = "s3"

type R2Store struct {
	client    *http.Client
	options   R2Options
	endpoint  string
	publicURL string
}

func NewR2Store(options R2Options) (ObjectStore, error) {
	if !options.Enabled() {
		if options.PartiallyConfigured() {
			return nil, errors.New("R2 storage is partially configured")
		}
		return nil, nil
	}
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", options.AccountID)
	return &R2Store{
		client:    http.DefaultClient,
		options:   options,
		endpoint:  endpoint,
		publicURL: strings.TrimRight(options.PublicBaseURL, "/"),
	}, nil
}

func (s *R2Store) Delete(ctx context.Context, key string) error {
	req, err := s.newRequest(ctx, http.MethodDelete, key)
	if err != nil {
		return err
	}
	signR2Request(req, s.options, nil)
	return s.do(req)
}

func (s *R2Store) Head(ctx context.Context, key string) (ObjectInfo, error) {
	req, err := s.newRequest(ctx, http.MethodHead, key)
	if err != nil {
		return ObjectInfo{}, err
	}
	signR2Request(req, s.options, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return ObjectInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ObjectInfo{}, ErrObjectNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ObjectInfo{}, fmt.Errorf("r2 request failed: %s", resp.Status)
	}
	return readObjectInfo(key, resp), nil
}

func (s *R2Store) PublicURL(key string) string {
	if s.publicURL == "" || key == "" {
		return ""
	}
	return s.publicURL + "/" + escapePublicPath(key)
}

func (s *R2Store) newRequest(
	ctx context.Context,
	method string,
	key string,
) (*http.Request, error) {
	objectURL := s.objectURL(key)
	return http.NewRequestWithContext(ctx, method, objectURL, nil)
}

func (s *R2Store) objectURL(key string) string {
	escaped := escapeObjectKey(key)
	return s.endpoint + "/" + url.PathEscape(s.options.Bucket) + "/" + escaped
}

func (s *R2Store) do(req *http.Request) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("r2 request failed: %s", resp.Status)
}

func readObjectInfo(key string, resp *http.Response) ObjectInfo {
	return ObjectInfo{
		Key:         key,
		ContentType: resp.Header.Get("Content-Type"),
		SizeBytes:   resp.ContentLength,
	}
}

func escapeObjectKey(key string) string {
	parts := strings.Split(key, "/")
	for index, part := range parts {
		parts[index] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

func escapePublicPath(key string) string {
	return escapeObjectKey(strings.TrimLeft(key, "/"))
}
