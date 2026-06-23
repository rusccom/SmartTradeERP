package storefront

import (
	"net/url"
	"strings"
)

const (
	cardTransform = "width=500,height=500,fit=cover,format=auto,quality=80"
	mainTransform = "width=900,format=auto,quality=82"
	logoTransform = "height=64,format=auto,quality=85"
)

// imageURL joins the R2 public base URL with an object key, escaping each path
// segment. Returns "" when either part is missing.
func imageURL(base, key string) string {
	if base == "" || key == "" {
		return ""
	}
	return strings.TrimRight(base, "/") + "/" + escapeKey(key)
}

// transformedImageURL wraps a public image URL in a Cloudflare image-transform
// path. The r2.dev development domain cannot transform, so it is returned
// unchanged (matches the media feature's behaviour).
func transformedImageURL(base, key, transform string) string {
	full := imageURL(base, key)
	if full == "" || transform == "" {
		return full
	}
	scheme, rest, ok := strings.Cut(full, "://")
	if !ok {
		return full
	}
	host, path, ok := strings.Cut(rest, "/")
	if !ok || strings.HasSuffix(host, ".r2.dev") {
		return full
	}
	return scheme + "://" + host + "/cdn-cgi/image/" + transform + "/" + path
}

func escapeKey(key string) string {
	parts := strings.Split(strings.TrimLeft(key, "/"), "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}
