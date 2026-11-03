package smartbundle

import (
	"context"
	"encoding/json"
	"slices"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/validation"
)

const readRules = `SELECT COALESCE(r.price_mode,'fixed'),COALESCE(r.discount,0),COALESCE(r.slots,'{}')
	FROM catalog.product_variants v JOIN catalog.products p ON p.id=v.product_id
	LEFT JOIN storefront.smartbundle_rules r ON r.variant_id=v.id AND r.tenant_id=v.tenant_id
	WHERE v.id=$1 AND p.is_composite`

const readOptions = `SELECT vc.component_variant_id::text,p.name,
	jsonb_agg(jsonb_build_object('id',a.id,'name',COALESCE(a.name,'Default')) ORDER BY a.created_at,a.id)
	FROM catalog.variant_components vc JOIN catalog.product_variants v ON v.id=vc.component_variant_id
	JOIN catalog.products p ON p.id=v.product_id JOIN catalog.product_variants a ON a.product_id=p.id
	WHERE vc.variant_id=$1 GROUP BY vc.component_variant_id,p.name ORDER BY p.name,vc.component_variant_id`

const saveRules = `INSERT INTO storefront.smartbundle_rules(tenant_id,variant_id,price_mode,discount,slots)
	VALUES ($1,$2,$3,$4,$5) ON CONFLICT (tenant_id,variant_id) DO UPDATE SET
	price_mode=EXCLUDED.price_mode,discount=EXCLUDED.discount,slots=EXCLUDED.slots`

func readEditor(ctx context.Context, tx pgx.Tx, id string) (Editor, error) {
	var result Editor
	s := &result.Settings
	err := tx.QueryRow(ctx, readRules, id).Scan(&s.Mode, &s.Discount, &s.Slots)
	if err != nil {
		return result, err
	}
	rows, err := tx.Query(ctx, readOptions, id)
	if err != nil {
		return result, err
	}
	result.Components, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (EditableSlot, error) {
		var slot EditableSlot
		err := row.Scan(&slot.Key, &slot.Name, &slot.Options)
		return slot, err
	})
	s.Slots = currentRules(s.Slots, result.Components)
	return result, err
}

func currentRules(stored map[string]Rule, components []EditableSlot) map[string]Rule {
	out := make(map[string]Rule, len(components))
	for _, c := range components {
		rule := stored[c.Key]
		allowed := make([]string, 0, len(c.Options))
		for _, v := range c.Options {
			if slices.Contains(rule.Options, v.ID) {
				allowed = append(allowed, v.ID)
			}
		}
		rule.Options = allowed
		out[c.Key] = rule
	}
	return out
}

func saveEditor(ctx context.Context, tx db.DBTX, tenantID, id string, s Settings) error {
	rules, err := json.Marshal(s.Slots)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, saveRules, tenantID, id, s.Mode, s.Discount, rules)
	return err
}

func (h *SettingsHandler) editor(ctx context.Context, tenantID, id string, update *Settings) (Editor, error) {
	var result Editor
	if !validation.UUID(id) {
		return result, validation.ErrInvalidData
	}
	err := h.store.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		var err error
		result, err = readEditor(ctx, tx, id)
		if err != nil || update == nil {
			return err
		}
		if !update.valid(result.Components) {
			return validation.ErrInvalidData
		}
		result.Settings = *update
		return saveEditor(ctx, tx, tenantID, id, *update)
	})
	return result, err
}
