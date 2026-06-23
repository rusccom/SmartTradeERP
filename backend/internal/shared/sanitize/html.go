// Package sanitize cleans untrusted product description HTML against a
// safe-but-permissive allow-list before it is persisted, so the storefront can
// render pasted store HTML (Shopify body_html style) directly without any
// client-side scrubbing. The editor stores RAW innerHTML, so the policy must
// preserve structure (div/span/section, classes, a safe subset of inline
// style) while staying XSS-safe.
package sanitize

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// embedHostsRe restricts iframe sources to the YouTube and Vimeo embed players
// only. UGCPolicy strips <iframe> by default, so the embed allow is re-added on
// top, anchored to these hosts, and nothing else can smuggle in a frame.
var embedHostsRe = regexp.MustCompile(
	`^https://(www\.)?(youtube(-nocookie)?\.com|player\.vimeo\.com)/`)

// safeStyleProps is the inline CSS allow-list for pasted store HTML. It covers
// layout, box model, flexbox/grid, typography, background and visual effects so
// self-contained styled HTML (gradients, rounded cards, pill badges, grids)
// renders verbatim — but deliberately OMITS position / z-index / top / right /
// bottom / left, so styles can never break out of the description flow or
// overlay the page. Values are further constrained by safeStyleValue.
var safeStyleProps = []string{
	"display", "box-sizing", "visibility", "opacity", "overflow", "overflow-x", "overflow-y",
	"width", "min-width", "max-width", "height", "min-height", "max-height",
	"margin", "margin-top", "margin-right", "margin-bottom", "margin-left",
	"padding", "padding-top", "padding-right", "padding-bottom", "padding-left",
	"border", "border-top", "border-right", "border-bottom", "border-left",
	"border-width", "border-style", "border-color", "border-radius",
	"border-top-left-radius", "border-top-right-radius",
	"border-bottom-left-radius", "border-bottom-right-radius",
	"outline", "box-shadow", "float", "clear", "vertical-align",
	"flex", "flex-direction", "flex-wrap", "flex-grow", "flex-shrink", "flex-basis",
	"align-items", "align-content", "align-self",
	"justify-content", "justify-items", "justify-self",
	"gap", "row-gap", "column-gap", "order",
	"grid", "grid-template-columns", "grid-template-rows", "grid-template-areas",
	"grid-column", "grid-row", "grid-auto-flow", "grid-auto-columns", "grid-auto-rows",
	"font", "font-family", "font-size", "font-weight", "font-style", "font-variant",
	"line-height", "letter-spacing", "word-spacing", "white-space",
	"text-align", "text-decoration", "text-transform", "text-indent", "text-shadow",
	"color", "background", "background-color", "background-image",
	"background-position", "background-size", "background-repeat",
	"list-style", "list-style-type", "list-style-position", "transform", "transition",
}

// unsafeStyleTokens are substrings that disqualify a style value: external/JS
// resource loaders, old IE expression(), CSS comment tricks and escapes.
var unsafeStyleTokens = []string{
	"url(", "expression", "javascript:", "/*", "*/", "<", ">", "\\", "@",
}

// classElements lists the block/inline tags that may carry a class attribute.
// Classes are inert without our CSS but harmless, and they preserve the
// structure of pasted store HTML.
var classElements = []string{
	"div", "span", "p", "h1", "h2", "h3", "h4", "h5", "h6",
	"ul", "ol", "li", "table", "td", "th", "figure", "img", "a",
	"blockquote", "section",
}

var intRe = regexp.MustCompile(`^[0-9]+$`)

// Sanitizer holds a single compiled bluemonday policy. The policy is safe for
// concurrent use, so one instance is shared for the life of the process.
type Sanitizer struct {
	policy *bluemonday.Policy
}

// NewSanitizer builds the description allow-list on top of bluemonday's
// UGCPolicy, which already provides a safe default (formatting, lists, tables,
// links, http/https images) while stripping script, style elements, on*
// handlers, iframe, object, embed and form. We broaden it for pasted store HTML
// by allowing structural containers, class attributes, a safe inline-style
// subset and YouTube/Vimeo embeds.
//
// r2Host is retained in the signature only to keep callers stable; the image
// rule is no longer R2-anchored (UGCPolicy already restricts images to
// http/https), so any https image from imported HTML renders.
func NewSanitizer(r2Host string) *Sanitizer {
	_ = r2Host // intentionally unused: images are no longer locked to R2.
	p := bluemonday.UGCPolicy()
	p.AllowURLSchemes("http", "https", "mailto")
	allowStructure(p)
	allowClasses(p)
	allowSafeStyles(p)
	allowLinks(p)
	allowEmbeds(p)
	return &Sanitizer{policy: p}
}

// Sanitize returns clean HTML for the given untrusted input.
func (s *Sanitizer) Sanitize(raw string) string {
	return s.policy.Sanitize(raw)
}

// allowStructure permits the structural containers UGCPolicy omits, so pasted
// layout blocks keep their shape.
func allowStructure(p *bluemonday.Policy) {
	p.AllowElements("div", "span", "section", "figure", "figcaption",
		"h4", "h5", "h6")
}

// allowClasses permits an inert class attribute on common block/inline tags so
// the structure of imported HTML is preserved.
func allowClasses(p *bluemonday.Policy) {
	p.AllowAttrs("class").OnElements(classElements...)
}

// allowSafeStyles permits the broad inline-style allow-list (safeStyleProps),
// validating every value through safeStyleValue so no url()/expression()/JS or
// CSS-comment trick can ride inside the style attribute.
func allowSafeStyles(p *bluemonday.Policy) {
	p.AllowStyles(safeStyleProps...).MatchingHandler(safeStyleValue).Globally()
}

// safeStyleValue accepts a CSS value unless it contains a disallowed token, so
// gradients / lengths / shadows survive while resource loaders and escapes do not.
func safeStyleValue(value string) bool {
	low := strings.ToLower(value)
	for _, bad := range unsafeStyleTokens {
		if strings.Contains(low, bad) {
			return false
		}
	}
	return true
}

// allowLinks hardens the links UGCPolicy already permits for UGC.
func allowLinks(p *bluemonday.Policy) {
	p.RequireNoFollowOnLinks(true)
	p.AddTargetBlankToFullyQualifiedLinks(true)
}

// allowEmbeds re-adds YouTube/Vimeo iframes on top of UGCPolicy (which strips
// iframe by default), with framing attributes.
func allowEmbeds(p *bluemonday.Policy) {
	p.AllowAttrs("src").Matching(embedHostsRe).OnElements("iframe")
	p.AllowAttrs("width", "height").Matching(intRe).OnElements("iframe")
	p.AllowAttrs("frameborder", "allowfullscreen", "sandbox").OnElements("iframe")
}
