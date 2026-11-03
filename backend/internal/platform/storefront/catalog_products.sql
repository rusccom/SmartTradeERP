WITH published_categories AS (
    SELECT pc.product_id, c.slug, c.name, c.sort_order, parent.slug AS parent_slug
    FROM catalog.product_categories pc
    JOIN catalog.categories c ON c.id = pc.category_id AND c.tenant_id = pc.tenant_id AND (c.is_published OR $8)
    LEFT JOIN catalog.categories parent ON parent.id = c.parent_id AND parent.tenant_id = c.tenant_id AND (parent.is_published OR $8)
    WHERE pc.tenant_id = $1
), page AS (
    SELECT p.* FROM catalog.products p
    WHERE p.tenant_id = $1 AND (p.is_published OR $8)
      AND ($2 = '' OR p.slug = $2)
      AND ($3 <> 'builder' OR (NOT p.is_composite AND EXISTS (
          SELECT 1 FROM catalog.product_variants v WHERE v.tenant_id = p.tenant_id AND v.product_id = p.id AND (v.is_published OR $8)))) AND ($3 <> 'bags' OR p.is_composite)
      AND ($4 = '' OR EXISTS (SELECT 1 FROM published_categories c WHERE c.product_id = p.id
          AND (c.slug = $4 OR ($3 <> 'builder' AND c.parent_slug = $4))))
      AND ($5::uuid[] IS NULL OR EXISTS (
          SELECT 1 FROM catalog.product_variants v WHERE v.tenant_id = p.tenant_id AND v.product_id = p.id AND v.id = ANY($5)))
    ORDER BY p.created_at DESC, p.id LIMIT NULLIF($6::int, 0) OFFSET $7
)
SELECT p.name, p.slug, COALESCE(p.seo_title,''), COALESCE(p.seo_description,''), p.is_composite,
    COALESCE((SELECT jsonb_agg(jsonb_build_object('Slug',c.slug,'Name',c.name,'ParentSlug',c.parent_slug) ORDER BY c.sort_order,c.name,c.slug)
        FROM published_categories c WHERE c.product_id = p.id), '[]'),
    ARRAY(SELECT m.object_key FROM platform.media_objects m WHERE m.tenant_id = p.tenant_id AND m.owner_id = p.id
        AND m.owner_type = 'product' AND m.status = 'ready' ORDER BY m.is_primary DESC,m.sort_order,m.created_at),
    COALESCE((SELECT jsonb_agg(jsonb_build_object('ID',v.id,'Name',COALESCE(v.name,'Default'),
        'Images', ARRAY(SELECT m.object_key FROM platform.media_objects m WHERE m.tenant_id=v.tenant_id
            AND m.owner_type='variant' AND m.owner_id=v.id AND m.status='ready' ORDER BY m.is_primary DESC,m.sort_order,m.created_at),
        'Price',v.price,'SKU',COALESCE(v.sku_code,''),'Stock',COALESCE(
            (SELECT sum(sb.qty) FROM ledger.stock_balances sb WHERE sb.tenant_id = v.tenant_id AND sb.variant_id = v.id),0))
        ORDER BY v.created_at,v.id) FROM catalog.product_variants v
        WHERE v.tenant_id = p.tenant_id AND v.product_id = p.id AND (v.is_published OR $8) AND ($5::uuid[] IS NULL OR v.id = ANY($5))), '[]')
FROM page p ORDER BY p.created_at DESC,p.id
