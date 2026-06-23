package media

import "strings"

const thumbTransform = "width=96,height=96,fit=cover,format=auto,quality=75"

// thumbURL turns a public R2 URL into a Cloudflare image-transform URL. The
// r2.dev development domain cannot transform, so it is returned unchanged;
// pointing R2_PUBLIC_BASE_URL at a custom domain on the Cloudflare zone
// activates real thumbnails with no code change.
func thumbURL(publicURL string) string {
	if publicURL == "" {
		return ""
	}
	scheme, rest, ok := strings.Cut(publicURL, "://")
	if !ok {
		return publicURL
	}
	host, path, ok := strings.Cut(rest, "/")
	if !ok || strings.HasSuffix(host, ".r2.dev") {
		return publicURL
	}
	return scheme + "://" + host + "/cdn-cgi/image/" + thumbTransform + "/" + path
}
