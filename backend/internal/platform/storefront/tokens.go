package storefront

import (
	"html/template"
	"sort"
	"strings"
)

// buildTokenCSS merges theme default tokens with per-tenant overrides and emits
// CSS custom properties for a :root block. Keys and values are sanitized so a
// malicious override cannot break out of the declaration or inject markup.
func buildTokenCSS(defaults, overrides map[string]string) template.CSS {
	merged := make(map[string]string, len(defaults)+len(overrides))
	for key, value := range defaults {
		merged[key] = value
	}
	for key, value := range overrides {
		merged[key] = value
	}
	return template.CSS(renderTokens(merged))
}

// SanitizeOverrides keeps only override entries whose key is a known token of
// the given theme and whose value passes the CSS-value whitelist. The admin
// uses it before persisting tenant draft/published tokens, so the same rules
// guard both render time and write time.
func (rg *Registry) SanitizeOverrides(themeID string, input map[string]string) map[string]string {
	defaults := rg.DefaultTokens(themeID)
	out := make(map[string]string)
	for key, value := range input {
		name := sanitizeTokenKey(key)
		clean := sanitizeTokenValue(value)
		if name == "" || clean == "" {
			continue
		}
		if _, ok := defaults[name]; ok {
			out[name] = clean
		}
	}
	return out
}

func renderTokens(merged map[string]string) string {
	var builder strings.Builder
	for _, key := range sortedKeys(merged) {
		name := sanitizeTokenKey(key)
		value := sanitizeTokenValue(merged[key])
		if name == "" || value == "" {
			continue
		}
		builder.WriteString("--")
		builder.WriteString(name)
		builder.WriteString(":")
		builder.WriteString(value)
		builder.WriteString(";")
	}
	return builder.String()
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sanitizeTokenKey(raw string) string {
	cleaned := strings.ToLower(strings.TrimSpace(raw))
	if cleaned == "" || len(cleaned) > 40 {
		return ""
	}
	for _, char := range cleaned {
		if !isTokenKeyChar(char) {
			return ""
		}
	}
	return cleaned
}

func isTokenKeyChar(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-'
}

func sanitizeTokenValue(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" || len(cleaned) > 120 {
		return ""
	}
	for _, char := range cleaned {
		if !isTokenValueChar(char) {
			return ""
		}
	}
	return cleaned
}

func isTokenValueChar(char rune) bool {
	switch char {
	case ' ', '#', '.', ',', '%', '(', ')', '/', '_', '+', '-':
		return true
	}
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}
