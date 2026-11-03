package storefront

import (
	"cmp"
	"net/http"

	"smarterp/backend/internal/platform/storefront/catalog"
)

func requestScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	return "https"
}

func absoluteURL(r *http.Request, path string) string {
	return requestScheme(r) + "://" + r.Host + path
}

func canonicalURL(r *http.Request) string {
	return absoluteURL(r, r.URL.Path)
}

func seoHome(r *http.Request, brand brandRow) SeoVM {
	return SeoVM{
		Title:     brand.TenantName,
		Canonical: canonicalURL(r),
		Robots:    "index,follow",
		OGType:    "website",
	}
}

func seoList(r *http.Request, brand brandRow) SeoVM {
	return SeoVM{
		Title:     "Products — " + brand.TenantName,
		Canonical: canonicalURL(r),
		Robots:    "index,follow",
		OGType:    "website",
	}
}

func seoCart(r *http.Request, brand brandRow) SeoVM {
	title := "Cart"
	if shopLanguage(brand) == "pl" {
		title = "Koszyk"
	}
	return SeoVM{
		Title:  title + " — " + brand.TenantName,
		Robots: "noindex,nofollow",
		OGType: "website",
	}
}

func seoNotFound(r *http.Request, brand brandRow) SeoVM {
	return SeoVM{
		Title:  "Not found — " + brand.TenantName,
		Robots: "noindex,follow",
		OGType: "website",
	}
}

func seoProduct(r *http.Request, detail catalog.Product, ogImage string, price catalog.Money, hasPrice bool) SeoVM {
	title := cmp.Or(detail.SEOTitle, detail.Name)
	seo := SeoVM{
		Title:       title,
		Description: detail.SEODescription,
		Canonical:   canonicalURL(r),
		Robots:      "index,follow",
		OGType:      "product",
		OGImage:     ogImage,
	}
	if hasPrice {
		seo.JSONLD = productLD(productLDInput{
			Name: detail.Name, Description: detail.SEODescription, Image: ogImage,
			Price: price.Amount, Currency: price.Code, URL: canonicalURL(r),
		})
	}
	return seo
}
