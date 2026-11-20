import { Trash2 } from "lucide-react";
import { useI18n } from "@smarterp/i18n/useI18n";

function BundleComponentRow({ canRemove, onChange, onRemove, options, row }) {
  const { t } = useI18n();
  const choices = componentChoices(options, row);
  return (
    <div className="bundle-component-row">
      <select name="componentVariantID" aria-label={t("bundles.editor.component")}
        value={row.componentVariantID} onChange={(event) => onChange(row.id, event)}>
        <option value="">{t("bundles.editor.select")}</option>
        {choices.map((option) => <option key={option.id} value={option.id}>{option.label}</option>)}
      </select>
      <input name="qty" type="number" min="0.001" step="0.001" aria-label={t("bundles.editor.quantity")}
        value={row.qty} onChange={(event) => onChange(row.id, event)} />
      <button className="bundles-icon-btn" type="button" onClick={() => onRemove(row.id)} disabled={!canRemove} title={t("bundles.editor.remove")}>
        <Trash2 size={16} />
      </button>
    </div>
  );
}

function componentChoices(options, row) {
  if (!row.componentVariantID || options.some((item) => item.id === row.componentVariantID)) {
    return options;
  }
  const label = row.label || row.componentVariantID;
  return [{ id: row.componentVariantID, label, unit: row.unit || "" }, ...options];
}

export default BundleComponentRow;
