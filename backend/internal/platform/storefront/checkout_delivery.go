package storefront

import (
	"regexp"
	"strings"

	"smarterp/backend/internal/erp/orders"
	"smarterp/backend/internal/shared/validation"
)

type checkoutPickup struct {
	orders.Pickup
	Carrier string `json:"carrier"`
}

var polishPostalCode = regexp.MustCompile(`^([0-9]{2})-?([0-9]{3})$`)

var pickupMethods = map[string]string{
	"inpost": "Base Courier · InPost",
	"orlen":  "Base Courier · ORLEN Paczka",
	"dhl":    "Base Courier · DHL",
	"dpd":    "Base Courier · DPD Pickup",
}

func postalCode(country, raw string) (string, bool) {
	value := strings.TrimSpace(raw)
	if country != "PL" {
		return value, deliveryText(value, 32)
	}
	parts := polishPostalCode.FindStringSubmatch(value)
	if parts == nil {
		return "", false
	}
	return parts[1] + "-" + parts[2], true
}

func deliveryText(value string, limit int) bool {
	return !strings.ContainsRune(value, 0) && validation.Max(value, limit)
}

func (req checkoutRequest) deliveryDetails() (orders.Details, error) {
	if err := validateCustomer(req.Customer); err != nil {
		return orders.Details{}, err
	}
	code, valid := postalCode(req.Customer.Country, req.Customer.PostalCode)
	if !valid {
		return orders.Details{}, ErrCheckoutInvalid
	}
	shipping := orders.Address{Country: req.Customer.Country, PostalCode: code}
	if req.Pickup != nil {
		return req.Pickup.details(shipping)
	}
	shipping.Street = strings.TrimSpace(req.Customer.Address)
	if shipping.Street == "" || !deliveryText(shipping.Street, 512) {
		return orders.Details{}, ErrCheckoutInvalid
	}
	return orders.Details{DeliveryKind: "address", Shipping: shipping}, nil
}

func (p checkoutPickup) details(shipping orders.Address) (orders.Details, error) {
	method, supported := pickupMethods[p.Carrier]
	if !supported || shipping.Country != "PL" || p.Country != shipping.Country {
		return orders.Details{}, ErrCheckoutInvalid
	}
	point := p.Pickup
	point.Name, point.PointID = strings.TrimSpace(point.Name), strings.TrimSpace(point.PointID)
	point.Street, point.City = strings.TrimSpace(point.Street), strings.TrimSpace(point.City)
	code, valid := postalCode(point.Country, point.PostalCode)
	if !valid || !validPickupPoint(point) {
		return orders.Details{}, ErrCheckoutInvalid
	}
	point.PostalCode = code
	return orders.Details{DeliveryKind: "pickup", DeliveryMethod: method, Pickup: &point, Shipping: shipping}, nil
}

func validPickupPoint(p orders.Pickup) bool {
	return p.PointID != "" && p.Name != "" && p.Street != "" && p.City != "" &&
		deliveryText(p.PointID, 255) && deliveryText(p.Name, 255) && deliveryText(p.Street, 512) && deliveryText(p.City, 255)
}
