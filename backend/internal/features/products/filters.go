package products

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/shopspring/decimal"

	"smarterp/backend/internal/shared/httpx"
)

var ErrInvalidProductFilter = errors.New("invalid product filter")

func productListQuery(r *http.Request) (ProductListQuery, error) {
	list := httpx.ParseListQuery(r, productSortConfig(), []string{"is_composite"})
	stock, err := parseProductStockFilter(r.URL.Query())
	if err != nil {
		return ProductListQuery{}, err
	}
	return ProductListQuery{List: list, Stock: stock}, nil
}

func productSortConfig() httpx.SortConfig {
	return httpx.SortConfig{
		Allowed:  []string{"name", "created_at"},
		Fallback: "created_at",
	}
}

func parseProductStockFilter(values url.Values) (ProductStockFilter, error) {
	filter := ProductStockFilter{
		WarehouseID: strings.TrimSpace(values.Get("warehouse_id")),
		QtyState:    strings.TrimSpace(values.Get("qty_state")),
	}
	if !validQtyState(filter.QtyState) {
		return filter, ErrInvalidProductFilter
	}
	return parseProductQtyBounds(values, filter)
}

func parseProductQtyBounds(
	values url.Values,
	filter ProductStockFilter,
) (ProductStockFilter, error) {
	minQty, err := parseProductQty(values, "min_qty")
	if err != nil {
		return filter, err
	}
	maxQty, err := parseProductQty(values, "max_qty")
	filter.MinQty = minQty
	filter.MaxQty = maxQty
	return filter, err
}

func parseProductQty(values url.Values, key string) (*decimal.Decimal, error) {
	raw := strings.TrimSpace(values.Get(key))
	if raw == "" {
		return nil, nil
	}
	value, err := decimal.NewFromString(raw)
	if err != nil {
		return nil, ErrInvalidProductFilter
	}
	return &value, nil
}

func validQtyState(value string) bool {
	switch value {
	case "", "positive", "negative", "zero", "nonzero":
		return true
	default:
		return false
	}
}
