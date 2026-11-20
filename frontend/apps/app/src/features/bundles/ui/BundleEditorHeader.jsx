import { Save } from "lucide-react";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function BundleEditorHeader({ bundle, saving, loading, onSave }) {
  const { t } = useI18n();
  return (
    <header className="bundles-section-head">
      <div>
        <h2>{bundle.product_name}</h2>
        <p>{bundle.variant_name} · {bundle.unit}</p>
      </div>
      <button className="bundles-save-btn" type="button" onClick={onSave} disabled={saving || loading}>
        <Save size={16} /> {t(saving ? "bundles.editor.saving" : "bundles.editor.save")}
      </button>
    </header>
  );
}
