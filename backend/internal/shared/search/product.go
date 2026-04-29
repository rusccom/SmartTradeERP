package search

import "strings"

func AppendProductSearch(query string, args []any, value string) (string, []any) {
	value = CleanProductSearch(value)
	if value == "" {
		return query, args
	}
	query += ` AND ` + productPredicate(Position(args))
	return query, append(args, value)
}

func CleanProductSearch(value string) string {
	value = strings.TrimSpace(value)
	lower := strings.ToLower(value)
	for _, prefix := range productSearchPrefixes() {
		if strings.HasPrefix(lower, prefix) {
			return strings.TrimSpace(value[len(prefix):])
		}
	}
	return value
}

func productSearchPrefixes() []string {
	return []string{"sku:", "barcode:", "bar code:", "штрихкод:", "баркод:"}
}

func productPredicate(param string) string {
	return `(p.name ILIKE '%' || $` + param + ` || '%'
        OR EXISTS (SELECT 1 FROM catalog.product_variants v
            WHERE v.tenant_id=$1 AND v.product_id=p.id
            AND (` + productVariantPredicate(param) + `)))`
}

func productVariantPredicate(param string) string {
	return `COALESCE(v.name,'') ILIKE '%' || $` + param + ` || '%'
        OR COALESCE(v.sku_code,'') ILIKE '%' || $` + param + ` || '%'
        OR COALESCE(v.barcode,'') ILIKE '%' || $` + param + ` || '%'
        OR COALESCE(v.option1,'') ILIKE '%' || $` + param + ` || '%'
        OR COALESCE(v.option2,'') ILIKE '%' || $` + param + ` || '%'
        OR COALESCE(v.option3,'') ILIKE '%' || $` + param + ` || '%'`
}
