package storefront

import (
	"encoding/json"
	"html/template"
)

type productLDInput struct {
	Name        string
	Description string
	Image       string
	Price       string
	Currency    string
	URL         string
}

// jsonLDScript marshals a schema.org payload into a complete <script> element.
// encoding/json escapes <, > and & to \u-sequences, so the result is safe to
// emit verbatim in the document head.
func jsonLDScript(payload map[string]any) template.HTML {
	raw, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return template.HTML(`<script type="application/ld+json">` + string(raw) + `</script>`)
}

func productLD(in productLDInput) template.HTML {
	offer := map[string]any{"@type": "Offer", "price": in.Price, "url": in.URL, "availability": "https://schema.org/InStock"}
	if in.Currency != "" {
		offer["priceCurrency"] = in.Currency
	}
	payload := map[string]any{
		"@context": "https://schema.org",
		"@type":    "Product",
		"name":     in.Name,
		"offers":   offer,
	}
	if in.Description != "" {
		payload["description"] = in.Description
	}
	if in.Image != "" {
		payload["image"] = in.Image
	}
	return jsonLDScript(payload)
}
