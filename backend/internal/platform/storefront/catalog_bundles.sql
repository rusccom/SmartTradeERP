SELECT kit.id::text, cv.id::text, cp.name, COALESCE(cv.name,'Default'), vc.qty,
    COALESCE((SELECT sum(b.qty) FROM ledger.stock_balances b
        WHERE b.tenant_id=kit.tenant_id AND b.variant_id=cv.id),0),
    ($3 OR (cv.is_published AND cp.is_published))
FROM catalog.variant_components vc
JOIN catalog.product_variants kit ON kit.id=vc.variant_id AND kit.tenant_id=$1
JOIN catalog.product_variants cv ON cv.id=vc.component_variant_id AND cv.tenant_id=kit.tenant_id
JOIN catalog.products cp ON cp.id=cv.product_id AND cp.tenant_id=kit.tenant_id
WHERE kit.id=ANY($2::uuid[])
ORDER BY kit.id,cp.name,cv.id
