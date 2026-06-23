package products

import (
	"context"
	"strconv"
	"strings"
)

var cyrillicMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "i", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
}

// resolveSlug normalizes the handle, falling back to the product name so every
// product gets a products/<handle> URL even when the user leaves it blank.
func resolveSlug(slug, name string) string {
	normalized := normalizeSlug(slug)
	if normalized != "" {
		return normalized
	}
	return normalizeSlug(name)
}

// normalizeSlug lowercases the handle, transliterates Cyrillic to Latin, and
// keeps only [a-z0-9], collapsing every other run into a single hyphen so even
// a Cyrillic product name yields a URL-safe handle.
func normalizeSlug(value string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if latin, ok := cyrillicMap[r]; ok {
			b.WriteString(latin)
			dash = false
			continue
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			dash = false
			continue
		}
		if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// ensureUniqueSlug keeps the handle unique per tenant by appending -2, -3, ...
// when the base is already taken, so two products with the same name still get
// distinct products/<handle> URLs instead of failing.
func (s *Service) ensureUniqueSlug(ctx context.Context, tenantID, base, excludeID string) (string, error) {
	if base == "" {
		return "", nil
	}
	candidate := base
	for i := 2; i < 1000; i++ {
		taken, err := s.repo.SlugExists(ctx, tenantID, candidate, excludeID)
		if err != nil {
			return "", err
		}
		if !taken {
			return candidate, nil
		}
		candidate = base + "-" + strconv.Itoa(i)
	}
	return candidate, nil
}
