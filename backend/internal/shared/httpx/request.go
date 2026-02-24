package httpx

import (
    "encoding/json"
    "net/http"
    "strconv"
)

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

func Offset(page, perPage int) int {
    if page <= 1 {
        return 0
    }
    return (page - 1) * perPage
}

func parseInt(value string, fallback int) int {
    parsed, err := strconv.Atoi(value)
    if err != nil || parsed <= 0 {
        return fallback
    }
    return parsed
}
