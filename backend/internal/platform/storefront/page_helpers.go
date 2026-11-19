package storefront

import "html/template"

// Theme templates can collect ordered groups and index rows without changing catalog data.
var pageHelpers = template.FuncMap{
	"dict": func() map[string]any { return map[string]any{} },
	"set":  func(m map[string]any, key string, value any) string { m[key] = value; return "" },
	"push": func(m map[string]any, key string, value any) string {
		list, _ := m[key].([]any)
		m[key] = append(list, value)
		return ""
	},
}
