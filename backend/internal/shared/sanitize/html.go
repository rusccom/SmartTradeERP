// Package sanitize cleans untrusted product description HTML against a
// safe-but-permissive allow-list before it is persisted, so the storefront can
// render pasted store HTML (Shopify body_html style) directly without any
// client-side scrubbing. The editor stores RAW innerHTML, so the policy must
// preserve structure (div/span/section, classes, a safe subset of inline
// style) while staying XSS-safe.
package sanitize

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

// embedHostsRe restricts iframe sources to the YouTube and Vimeo embed players
// only. UGCPolicy strips <iframe> by default, so the embed allow is re-added on
// top, anchored to these hosts, and nothing else can smuggle in a frame.
var embedHostsRe = regexp.MustCompile(
	`^https://(www\.)?(youtube(-nocookie)?\.com|player\.vimeo\.com)/`)

// colorRe constrains color / background-color values to hex / rgb / rgba so no
// url() or other CSS can ride along inside the style attribute.
var colorRe = regexp.MustCompile(
	`(?i)^(#[0-9a-f]{3}|#[0-9a-f]{6}|rgb\(\s*\d{1,3}(\s*,\s*\d{1,3}){2}\s*\)|` +
		`rgba\(\s*\d{1,3}(\s*,\s*\d{1,3}){2}\s*,\s*(0|1|0?\.\d+)\s*\))$`)

// sizeRe constrains width / max-width to a plain px or % length (no url(), no
// calc()), so layout hints survive a save while CSS injection cannot.
var sizeRe = regexp.MustCompile(`(?i)^\d{1,4}(px|%)$`)

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

// allowSafeStyles permits a constrained subset of inline style. Each property
// is matched to an enum or value regex, so no url()-bearing or positioning CSS
// can ride inside the style attribute.
func allowSafeStyles(p *bluemonday.Policy) {
	allowTextStyles(p)
	allowColorStyles(p)
	allowBoxStyles(p)
}

// allowTextStyles permits text alignment and font styling globally.
func allowTextStyles(p *bluemonday.Policy) {
	p.AllowStyles("text-align").
		MatchingEnum("left", "center", "right", "justify").Globally()
	p.AllowStyles("font-weight").MatchingEnum("normal", "bold",
		"100", "200", "300", "400", "500", "600", "700", "800", "900").
		Globally()
	p.AllowStyles("font-style").MatchingEnum("italic", "normal").Globally()
	p.AllowStyles("text-decoration").
		MatchingEnum("underline", "line-through", "none").Globally()
}

// allowColorStyles permits constrained text and background colors globally.
func allowColorStyles(p *bluemonday.Policy) {
	p.AllowStyles("color").Matching(colorRe).Globally()
	p.AllowStyles("background-color").Matching(colorRe).Globally()
}

// allowBoxStyles permits sizing and float hints globally. position / z-index
// are never allowed, so styles cannot break out of the content flow.
func allowBoxStyles(p *bluemonday.Policy) {
	p.AllowStyles("width", "max-width").Matching(sizeRe).Globally()
	p.AllowStyles("float").MatchingEnum("left", "right", "none").Globally()
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
