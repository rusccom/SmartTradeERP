package orders

import (
	"time"

	"github.com/shopspring/decimal"
)

type Address struct {
	Name       string `json:"name"`
	Company    string `json:"company"`
	Street     string `json:"street"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	Region     string `json:"region"`
	Country    string `json:"country"`
	TaxID      string `json:"tax_id,omitempty"`
}

type Pickup struct {
	Country    string `json:"country"`
	Name       string `json:"name"`
	PointID    string `json:"point_id"`
	Street     string `json:"street"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
}

type Details struct {
	Nickname       string  `json:"nickname"`
	DeliveryKind   string  `json:"delivery_kind"`
	DeliveryMethod string  `json:"delivery_method"`
	DeliveryVAT    string  `json:"delivery_vat"`
	PaymentMethod  string  `json:"payment_method"`
	COD            bool    `json:"cod"`
	Shipping       Address `json:"shipping"`
	Billing        Address `json:"billing"`
	Pickup         *Pickup `json:"pickup"`
}

type LineInput struct {
	ID         string          `json:"id,omitempty"`
	VariantID  string          `json:"variant_id"`
	Qty        decimal.Decimal `json:"qty"`
	VAT        string          `json:"vat_rate"`
	VATMode    string          `json:"vat_mode,omitempty"`
	VATCountry string          `json:"vat_country,omitempty"`
}

type Line struct {
	LineInput
	ProductID   string          `json:"-"`
	Thumbnail   string          `json:"image_thumb_url,omitempty"`
	ProductName string          `json:"product_name"`
	VariantName string          `json:"variant_name"`
	SKU         string          `json:"sku_code"`
	Unit        string          `json:"unit"`
	Price       decimal.Decimal `json:"unit_price"`
	Total       decimal.Decimal `json:"total_amount"`
}

type Entry struct {
	RequestID         string           `json:"request_id"`
	Revision          int64            `json:"revision"`
	Name              string           `json:"customer_name"`
	Email             string           `json:"customer_email"`
	Phone             string           `json:"customer_phone"`
	BillingCustomerID string           `json:"billing_customer_id"`
	InvoiceRequested  bool             `json:"invoice_requested"`
	DeliveryCost      *decimal.Decimal `json:"delivery_cost"`
	Details           Details          `json:"details"`
	Items             []LineInput      `json:"items"`
}

type Summary struct {
	ID       string           `json:"id"`
	Number   string           `json:"number"`
	Source   string           `json:"source"`
	Name     string           `json:"customer_name"`
	Status   string           `json:"status"`
	Currency string           `json:"currency_code"`
	Total    decimal.Decimal  `json:"total_amount"`
	Grand    *decimal.Decimal `json:"total_with_delivery"`
	Created  time.Time        `json:"created_at"`
}

type Invoice struct {
	ID       string    `json:"id"`
	Number   string    `json:"number"`
	Format   string    `json:"format"`
	IssuedAt time.Time `json:"issued_at"`
}

type Card struct {
	Summary
	Quantity          decimal.Decimal  `json:"quantity"`
	Email             string           `json:"customer_email"`
	Phone             string           `json:"customer_phone"`
	Revision          int64            `json:"revision"`
	Details           Details          `json:"details"`
	DeliveryCost      *decimal.Decimal `json:"delivery_cost"`
	Paid              *decimal.Decimal `json:"paid_amount"`
	InvoiceRequested  bool             `json:"invoice_requested"`
	BillingCustomerID string           `json:"billing_customer_id"`
	Items             []Line           `json:"items"`
	Invoice           *Invoice         `json:"invoice"`
}

type Quote struct {
	Items    []Line           `json:"items"`
	Total    decimal.Decimal  `json:"total_amount"`
	Grand    *decimal.Decimal `json:"total_with_delivery"`
	Currency string           `json:"currency_code"`
	Quantity decimal.Decimal  `json:"quantity"`
}
