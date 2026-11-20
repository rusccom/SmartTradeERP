package storefront

import (
	"slices"
	"smarterp/backend/internal/shared/countries"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type checkoutCountry struct{ Code, Name string }

func checkoutCountries(locale string) []checkoutCountry {
	tag := language.Make(locale)
	names := display.Regions(tag)
	if names == nil {
		names = display.Regions(language.English)
	}
	codes := countries.All()
	options := make([]checkoutCountry, 0, len(codes))
	for _, code := range codes {
		region := language.MustParseRegion(code)
		options = append(options, checkoutCountry{code, names.Name(region)})
	}
	sorter := collate.New(tag)
	slices.SortFunc(options, func(a, b checkoutCountry) int { return sorter.CompareString(a.Name, b.Name) })
	return options
}
