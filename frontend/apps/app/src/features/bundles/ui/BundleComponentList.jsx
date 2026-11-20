import { Plus } from "lucide-react";
import { useI18n } from "@smarterp/i18n/useI18n";
import BundleComponentRow from "./BundleComponentRow";

export default function BundleComponentList({ rows, options, onAdd, onChange, onRemove }) {
  const { t } = useI18n();
  return (
    <>
      <div className="bundle-components-head">
        <span>{t("bundles.editor.component")}</span>
        <span>{t("bundles.editor.quantity")}</span>
        <span />
      </div>
      <div className="bundle-components-list">
        {rows.map(row => <BundleComponentRow key={row.id} row={row} options={options}
          onChange={onChange} onRemove={onRemove} canRemove={rows.length > 1} />)}
      </div>
      <button className="bundles-add-btn" type="button" onClick={onAdd} disabled={!options.length}>
        <Plus size={16} /> {t("bundles.editor.add")}
      </button>
    </>
  );
}
