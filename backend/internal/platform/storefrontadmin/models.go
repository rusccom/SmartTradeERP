package storefrontadmin

// SectionInput is one toggleable home section in display order.
type SectionInput struct {
	Key     string `json:"key"`
	Enabled bool   `json:"enabled"`
}

// Settings is the tenant's storefront configuration returned to the dashboard:
// the live look plus the unpublished draft (theme, tokens, sections).
type Settings struct {
	ThemeID         string            `json:"theme_id"`
	DraftThemeID    string            `json:"draft_theme_id"`
	PublishedTokens map[string]string `json:"published_tokens"`
	DraftTokens     map[string]string `json:"draft_tokens"`
	Sections        []SectionInput    `json:"sections"`
	DraftSections   []SectionInput    `json:"draft_sections"`
}

// DraftRequest saves the in-progress theme choice, token overrides and sections.
type DraftRequest struct {
	ThemeID  string            `json:"theme_id"`
	Tokens   map[string]string `json:"tokens"`
	Sections []SectionInput    `json:"sections"`
}

// Preview is a short-lived token and the storefront URL that renders the draft.
type Preview struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}
