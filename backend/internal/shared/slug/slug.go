// Package slug turns human names into DNS-safe labels for shop subdomains.
package slug

import "strings"

const maxLen = 40

// translit maps characters that ToLower alone cannot make DNS-safe: Cyrillic
// (RU/UA) and Polish diacritics, matching the product's PL/EN/RU audience.
var translit = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "i", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "c", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	'і': "i", 'ї': "yi", 'є': "ye", 'ґ': "g",
	'ą': "a", 'ć': "c", 'ę': "e", 'ł': "l", 'ń': "n", 'ó': "o",
	'ś': "s", 'ź': "z", 'ż': "z",
}

// Make converts a display name into a lowercase DNS label: transliterated,
// [a-z0-9-] only, no leading/trailing/double dashes, capped, never empty.
func Make(name string) string {
	cleaned := sanitize(transliterate(strings.ToLower(name)))
	if cleaned == "" {
		return "shop"
	}
	return cleaned
}

func transliterate(value string) string {
	out := strings.Builder{}
	for _, r := range value {
		if mapped, ok := translit[r]; ok {
			out.WriteString(mapped)
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}

func sanitize(value string) string {
	out := make([]rune, 0, len(value))
	for _, r := range value {
		if isLabelRune(r) {
			out = append(out, r)
		} else if len(out) > 0 && out[len(out)-1] != '-' {
			out = append(out, '-')
		}
	}
	return trimLabel(string(out))
}

func isLabelRune(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= '0' && r <= '9'
}

func trimLabel(value string) string {
	value = strings.Trim(value, "-")
	if len(value) > maxLen {
		value = strings.Trim(value[:maxLen], "-")
	}
	return value
}
