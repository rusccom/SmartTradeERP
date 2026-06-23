package storefront

import "net/http"

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
	return SeoVM{
		Title:  "Cart — " + brand.TenantName,
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

func seoProduct(r *http.Request, detail productDetailRow, ogImage string, price MoneyVM, hasPrice bool) SeoVM {
	title := detail.SEOTitle
	if title == "" {
		title = detail.Name
	}
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
