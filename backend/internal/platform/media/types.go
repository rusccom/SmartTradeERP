package media

import (
	"strings"
	"unicode"
)

const MaxUploadBytes = 8 << 20

func validContentType(value string) bool {
	switch value {
	case "image/jpeg", "image/png", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}

func fileExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ""
	}
}

func validOwnerType(value string) bool {
	if value == "" || len(value) > 80 {
		return false
	}
	for _, char := range value {
		if !validOwnerTypeChar(char) {
			return false
		}
	}
	return true
}

func validOwnerTypeChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' || char == '-'
}

func cleanFileName(value string) string {
	cleaned := strings.ReplaceAll(strings.TrimSpace(value), "\\", "/")
	if cleaned == "" {
		return ""
	}
	parts := strings.Split(cleaned, "/")
	cleaned = strings.TrimSpace(parts[len(parts)-1])
	runes := []rune(cleaned)
	if len(runes) <= 255 {
		return cleaned
	}
	return string(runes[:255])
}
