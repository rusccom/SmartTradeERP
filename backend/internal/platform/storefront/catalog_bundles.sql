SELECT vc.variant_id::text, COALESCE(r.price_mode,'fixed'), COALESCE(r.discount,0),
    vc.component_variant_id::text, cp.name, vc.qty,
    COALESCE((config.slot->>'discount')::numeric,0),
    COALESCE((SELECT jsonb_agg(jsonb_build_object('ID',opt.id,'Name',COALESCE(opt.name,'Default'),
        'Price',opt.price,'Stock',COALESCE((SELECT sum(b.qty) FROM ledger.stock_balances b
            WHERE b.variant_id = opt.id),0)) ORDER BY (opt.id = cv.id) DESC,opt.created_at,opt.id)
        FROM catalog.product_variants opt
        WHERE opt.product_id = cv.product_id AND opt.is_published AND cp.is_published
          AND (opt.id = cv.id OR config.slot->'options' ? opt.id::text)), '[]')
FROM catalog.variant_components vc
JOIN catalog.product_variants cv ON cv.id = vc.component_variant_id
JOIN catalog.products cp ON cp.id = cv.product_id
LEFT JOIN storefront.smartbundle_rules r ON r.variant_id = vc.variant_id AND r.tenant_id = $1
CROSS JOIN LATERAL (SELECT r.slots->vc.component_variant_id::text AS slot) config
WHERE vc.variant_id = ANY($2::uuid[]) AND cv.tenant_id = $1
ORDER BY vc.variant_id, cp.name, vc.component_variant_id
