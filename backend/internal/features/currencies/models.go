package currencies

type Currency struct {
	ID             string `json:"id"`
	CurrencyID     string `json:"currency_id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	DisplaySymbol  string `json:"display_symbol"`
	DecimalPlaces  int    `json:"decimal_places"`
	IsBase         bool   `json:"is_base"`
	IsEnabled      bool   `json:"is_enabled"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type CurrencyOption struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	DecimalPlaces int    `json:"decimal_places"`
}

type CreateRequest struct {
	CurrencyID    string `json:"currency_id"`
	DisplaySymbol string `json:"display_symbol"`
	IsBase        bool   `json:"is_base"`
}
