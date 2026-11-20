package orders

import "context"

type thumbnailSource interface {
	PrimaryThumbsByOwners(context.Context, string, string, []string) (map[string]string, error)
}

func (h Handler) thumbnails(ctx context.Context, tenant string, card *Card) {
	if h.Media == nil || len(card.Items) == 0 {
		return
	}
	var products, variants []string
	for _, line := range card.Items {
		products = append(products, line.ProductID)
		variants = append(variants, line.VariantID)
	}
	p, _ := h.Media.PrimaryThumbsByOwners(ctx, tenant, "product", products)
	v, _ := h.Media.PrimaryThumbsByOwners(ctx, tenant, "variant", variants)
	for i, line := range card.Items {
		card.Items[i].Thumbnail = v[line.VariantID]
		if card.Items[i].Thumbnail == "" {
			card.Items[i].Thumbnail = p[line.ProductID]
		}
	}
}
