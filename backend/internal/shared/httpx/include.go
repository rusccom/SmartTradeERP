package httpx

import (
	"net/http"
	"strings"
)

func ParseIncludes(r *http.Request) map[string]bool {
	result := make(map[string]bool)
	for _, item := range strings.Split(r.URL.Query().Get("include"), ",") {
		value := strings.TrimSpace(item)
		if value != "" {
			result[value] = true
		}
	}
	return result
}

func HasInclude(includes map[string]bool, key string) bool {
	return includes[key]
}
