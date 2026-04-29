package search

func AppendProductSearch(query string, args []any, value string) (string, []any) {
	if value == "" {
		return query, args
	}
	query += ` AND ` + productPredicate(Position(args))
	return query, append(args, value)
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
        OR COALESCE(v.barcode,'') ILIKE '%' || $` + param + ` || '%'`
}
