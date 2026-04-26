package products

func (f ProductStockFilter) HasQtyRule() bool {
	return f.MinQty != nil || f.MaxQty != nil || f.QtyState != ""
}

func (f ProductStockFilter) HasWarehouse() bool {
	return f.WarehouseID != ""
}

func appendProductStockFilter(
	query string,
	args []any,
	filter ProductStockFilter,
) (string, []any) {
	if !filter.HasQtyRule() {
		return query, args
	}
	query += ` AND EXISTS (SELECT 1 FROM catalog.product_variants v
        WHERE v.tenant_id=$1 AND v.product_id=p.id`
	query, args = appendStockQtyPredicates(query, args, filter)
	return query + `)`, args
}

func appendVariantStockFilter(
	query string,
	args []any,
	filter ProductStockFilter,
) (string, []any) {
	if !filter.HasQtyRule() {
		return query, args
	}
	return appendStockQtyPredicates(query, args, filter)
}

func appendStockQtyPredicates(
	query string,
	args []any,
	filter ProductStockFilter,
) (string, []any) {
	qtySQL := stockQtySQL(filter, &args)
	query, args = appendMinQtyPredicate(query, args, qtySQL, filter)
	query, args = appendMaxQtyPredicate(query, args, qtySQL, filter)
	return appendQtyStatePredicate(query, qtySQL, filter.QtyState), args
}

func stockQtySQL(filter ProductStockFilter, args *[]any) string {
	if !filter.HasWarehouse() {
		return globalStockQtySQL()
	}
	*args = append(*args, filter.WarehouseID)
	return warehouseStockQtySQL(position(len(*args)))
}

func globalStockQtySQL() string {
	return `COALESCE((SELECT l.running_qty FROM ledger.cost_ledger l
        WHERE l.tenant_id=$1 AND l.variant_id=v.id
        ORDER BY l.sequence_num DESC LIMIT 1),0)`
}

func warehouseStockQtySQL(warehousePosition string) string {
	return `COALESCE((SELECT SUM(CASE WHEN l.type='IN' THEN l.qty
        WHEN l.type='OUT' THEN -l.qty ELSE 0 END) FROM ledger.cost_ledger l
        WHERE l.tenant_id=$1 AND l.variant_id=v.id
        AND l.warehouse_id::text=$` + warehousePosition + `),0)`
}

func appendMinQtyPredicate(
	query string,
	args []any,
	qtySQL string,
	filter ProductStockFilter,
) (string, []any) {
	if filter.MinQty == nil {
		return query, args
	}
	query += ` AND (` + qtySQL + `) >= $` + position(len(args)+1)
	return query, append(args, *filter.MinQty)
}

func appendMaxQtyPredicate(
	query string,
	args []any,
	qtySQL string,
	filter ProductStockFilter,
) (string, []any) {
	if filter.MaxQty == nil {
		return query, args
	}
	query += ` AND (` + qtySQL + `) <= $` + position(len(args)+1)
	return query, append(args, *filter.MaxQty)
}

func appendQtyStatePredicate(query string, qtySQL, state string) string {
	switch state {
	case "positive":
		return query + ` AND (` + qtySQL + `) > 0`
	case "negative":
		return query + ` AND (` + qtySQL + `) < 0`
	case "zero":
		return query + ` AND (` + qtySQL + `) = 0`
	case "nonzero":
		return query + ` AND (` + qtySQL + `) <> 0`
	default:
		return query
	}
}

func appendWarehouseScope(
	query string,
	args []any,
	filter ProductStockFilter,
) (string, []any) {
	if !filter.HasWarehouse() {
		return query, args
	}
	query += ` AND w.id::text=$` + position(len(args)+1)
	return query, append(args, filter.WarehouseID)
}
