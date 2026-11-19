package storefrontadmin

type Settings struct {
	ThemeID      string `json:"theme_id"`
	DraftThemeID string `json:"draft_theme_id"`
	ShopURL      string `json:"shop_url"`
}

type DraftRequest struct {
	ThemeID string `json:"theme_id"`
}

type Preview struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}
