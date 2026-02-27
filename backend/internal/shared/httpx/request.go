package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type SortConfig struct {
	Allowed  []string
	Fallback string
}

type ListQuery struct {
	Page    int
	PerPage int
	SortBy  string
	SortDir string
	Search  string
	Filters map[string]string
}

func DecodeJSON(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func ParsePagination(r *http.Request) (int, int) {
	page := parseInt(r.URL.Query().Get("page"), 1)
	perPage := parseInt(r.URL.Query().Get("per_page"), 20)
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func ParseSort(r *http.Request, allowedFields []string, fallbackField string) (string, string) {
	sortBy := r.URL.Query().Get("sort_by")
	sortDir := r.URL.Query().Get("sort_dir")
	if !contains(allowedFields, sortBy) {
		sortBy = fallbackField
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}
	return sortBy, sortDir
}

func ParseSearch(r *http.Request) string {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	if len(search) <= 200 {
		return search
	}
	return search[:200]
}

func ParseFilters(r *http.Request, allowedKeys []string) map[string]string {
	filters := make(map[string]string)
	for _, key := range allowedKeys {
		value := r.URL.Query().Get(key)
		if value != "" {
			filters[key] = value
		}
	}
	return filters
}

func ParseListQuery(r *http.Request, sortCfg SortConfig, filterKeys []string) ListQuery {
	page, perPage := ParsePagination(r)
	sortBy, sortDir := ParseSort(r, sortCfg.Allowed, sortCfg.Fallback)
	search := ParseSearch(r)
	filters := ParseFilters(r, filterKeys)
	return ListQuery{
		Page: page, PerPage: perPage, SortBy: sortBy,
		SortDir: sortDir, Search: search, Filters: filters,
	}
}

func Offset(page, perPage int) int {
	if page <= 1 {
		return 0
	}
	return (page - 1) * perPage
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
