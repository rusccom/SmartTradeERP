WITH page AS (
    SELECT p.* FROM catalog.products p
    WHERE p.tenant_id = $1 AND p.is_published
      AND ($2 = '' OR p.slug = $2)
      AND ($3 <> 'builder' OR (NOT p.is_composite AND EXISTS (
          SELECT 1 FROM catalog.product_variants v WHERE v.product_id = p.id AND v.is_published))) AND ($3 <> 'bags' OR p.is_composite)
      AND ($4 = '' OR p.slug = ANY($8::text[]) OR EXISTS (
          SELECT 1 FROM catalog.product_categories pc JOIN catalog.categories c ON c.id = pc.category_id
          WHERE pc.product_id = p.id AND c.is_published AND c.slug = $4))
      AND ($5::uuid[] IS NULL OR EXISTS (
          SELECT 1 FROM catalog.product_variants v WHERE v.product_id = p.id AND v.id = ANY($5)))
    ORDER BY p.created_at DESC, p.id LIMIT NULLIF($6::int, 0) OFFSET $7
)
SELECT p.name, p.slug, COALESCE(p.seo_title,''), COALESCE(p.seo_description,''), p.is_composite,
    COALESCE((SELECT jsonb_agg(jsonb_build_object('Slug',c.slug,'Name',c.name) ORDER BY c.sort_order,c.name)
        FROM catalog.product_categories pc JOIN catalog.categories c ON c.id = pc.category_id
        WHERE pc.product_id = p.id AND c.is_published), '[]'),
    ARRAY(SELECT m.object_key FROM platform.media_objects m WHERE m.owner_id = p.id
        AND m.owner_type = 'product' AND m.status = 'ready' ORDER BY m.is_primary DESC,m.sort_order,m.created_at),
    COALESCE((SELECT jsonb_agg(jsonb_build_object('ID',v.id,'Name',COALESCE(v.name,'Default'),
        'Price',v.price,'SKU',COALESCE(v.sku_code,''),'Stock',COALESCE(
            (SELECT sum(sb.qty) FROM ledger.stock_balances sb WHERE sb.variant_id = v.id),0))
        ORDER BY v.created_at,v.id) FROM catalog.product_variants v
        WHERE v.product_id = p.id AND v.is_published AND ($5::uuid[] IS NULL OR v.id = ANY($5))), '[]')
FROM page p ORDER BY p.created_at DESC,p.id
