package products

import (
	"strings"

	"smarterp/backend/internal/shared/sanitize"
)

// descriptionMaxBytes caps the stored description so a single product can never
// balloon the row; 256 KiB is far above any real rich-text body.
const descriptionMaxBytes = 256 * 1024

// emptyDocs are the markers TipTap leaves behind for an otherwise blank editor;
// they are collapsed to "" so an "empty" description is stored as NULL-like "".
var emptyDocs = map[string]bool{"": true, "<p></p>": true, "<p><br></p>": true}

// sanitizeDescriptionHTML cleans the raw HTML and collapses an empty document
// to "" so blank editors do not persist stray markup.
func sanitizeDescriptionHTML(s *sanitize.Sanitizer, html string) string {
	if s == nil {
		return ""
	}
	clean := strings.TrimSpace(s.Sanitize(html))
	if emptyDocs[clean] {
		return ""
	}
	return clean
}

// validateDescription enforces the stored byte cap on the (already sanitized)
// HTML so an oversized body is rejected as invalid data.
func validateDescription(html string) bool {
	return len(html) <= descriptionMaxBytes
}
