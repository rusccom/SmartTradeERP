// Package sanitize cleans untrusted product description HTML against a strict
// allow-list before it is persisted, so the storefront can render it directly
// without any client-side scrubbing.
package sanitize

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

// embedHostsRe restricts iframe sources to the YouTube and Vimeo embed players
// only, so the editor's "Video" button can never smuggle in arbitrary frames.
var embedHostsRe = regexp.MustCompile(
	`^https://(www\.)?(youtube(-nocookie)?\.com|player\.vimeo\.com)/`)

// colorRe constrains the text-color style value to hex / rgb / rgba so the
// editor's color button survives a save while no other CSS can ride along.
var colorRe = regexp.MustCompile(
	`(?i)^(#[0-9a-f]{3}|#[0-9a-f]{6}|rgb\(\s*\d{1,3}(\s*,\s*\d{1,3}){2}\s*\)|` +
		`rgba\(\s*\d{1,3}(\s*,\s*\d{1,3}){2}\s*,\s*(0|1|0?\.\d+)\s*\))$`)

var intRe = regexp.MustCompile(`^[0-9]+$`)

// Sanitizer holds a single compiled bluemonday policy. The policy is safe for
// concurrent use, so one instance is shared for the life of the process.
type Sanitizer struct {
	policy *bluemonday.Policy
}

// NewSanitizer builds the description allow-list from an EMPTY policy (not
// UGCPolicy), so the only image rule is the R2-anchored one below — otherwise
// UGCPolicy's permissive AllowImages would let any https image through.
// r2Host is the R2 public base URL (scheme+host).
func NewSanitizer(r2Host string) *Sanitizer {
	p := bluemonday.NewPolicy()
	p.AllowURLSchemes("http", "https", "mailto")
	allowFormatting(p)
	allowAlignment(p)
	allowColor(p)
	allowLinks(p)
	allowImages(p, imageSrcRegexp(r2Host))
	allowTables(p)
	allowEmbeds(p)
	return &Sanitizer{policy: p}
}

// Sanitize returns clean HTML for the given untrusted input.
func (s *Sanitizer) Sanitize(raw string) string {
	return s.policy.Sanitize(raw)
}

// allowFormatting permits the inline and block text tags the editor emits.
func allowFormatting(p *bluemonday.Policy) {
	p.AllowElements("p", "br", "h1", "h2", "h3", "strong", "em", "u", "s",
		"code", "pre", "span", "blockquote", "hr")
	p.AllowLists()
}

// allowAlignment permits text-align via the "style" attribute on the block
// elements TipTap applies alignment to. AllowStyles filters the style value, so
// only the listed property/values survive — no raw style attribute is allowed.
func allowAlignment(p *bluemonday.Policy) {
	p.AllowStyles("text-align").
		MatchingEnum("left", "center", "right", "justify").
		OnElements("p", "h1", "h2", "h3", "td", "th")
}

// allowColor permits a constrained text color on <span> (TipTap Color mark),
// matched to hex/rgb/rgba so nothing else can ride inside the style attribute.
func allowColor(p *bluemonday.Policy) {
	p.AllowStyles("color").Matching(colorRe).OnElements("span")
}

// allowLinks permits safe outbound links and hardens them for UGC.
func allowLinks(p *bluemonday.Policy) {
	p.AllowAttrs("href").
		Matching(regexp.MustCompile(`(?i)^(https?|mailto):`)).
		OnElements("a")
	p.RequireNoFollowOnLinks(true)
	p.AddTargetBlankToFullyQualifiedLinks(true)
}

// allowImages permits images whose src points at the configured R2 host only.
func allowImages(p *bluemonday.Policy, src *regexp.Regexp) {
	p.AllowAttrs("src").Matching(src).OnElements("img")
	p.AllowAttrs("alt").OnElements("img")
	p.AllowAttrs("width", "height").Matching(intRe).OnElements("img")
}

// allowTables permits table structure and cell spans.
func allowTables(p *bluemonday.Policy) {
	p.AllowElements("table", "thead", "tbody", "tr", "th", "td")
	p.AllowAttrs("colspan", "rowspan").Matching(intRe).OnElements("th", "td")
}

// allowEmbeds permits YouTube/Vimeo iframes only, with framing attributes.
func allowEmbeds(p *bluemonday.Policy) {
	p.AllowAttrs("src").Matching(embedHostsRe).OnElements("iframe")
	p.AllowAttrs("width", "height").Matching(intRe).OnElements("iframe")
	p.AllowAttrs("frameborder", "allowfullscreen", "sandbox").OnElements("iframe")
}

// imageSrcRegexp anchors image sources to the R2 host; a blank host falls back
// to a regex that matches nothing, so a misconfigured deploy strips every image
// rather than silently widening the allow-list to all hosts.
func imageSrcRegexp(r2Host string) *regexp.Regexp {
	if r2Host == "" {
		return regexp.MustCompile(`[^\s\S]`)
	}
	return regexp.MustCompile(`^` + regexp.QuoteMeta(r2Host))
}
