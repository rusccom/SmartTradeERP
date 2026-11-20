import { useI18n } from "@smarterp/i18n/useI18n";

import BundleComponentList from "./BundleComponentList";
import BundleEditorHeader from "./BundleEditorHeader";
import ComponentOptionsToolbar from "./ComponentOptionsToolbar";

function BundleComponentsEditor(props) {
  const { t } = useI18n();
  const { bundle, error, saving, loading, onSave, rows, options, onAdd, onChange, onRemove } = props;
  if (!bundle) return <section className="bundles-editor bundles-editor--empty">{t("bundles.editor.empty")}</section>;
  return (
    <section className="bundles-editor">
      <BundleEditorHeader bundle={bundle} saving={saving} loading={loading} onSave={onSave} />
      {error && <p className="bundles-error">{error}</p>}
      <ComponentOptionsToolbar canLoadMore={props.canLoadMoreOptions} loading={props.optionsLoading}
        onLoadMore={props.onLoadMoreOptions} onSearch={props.onSearchComponents} search={props.componentSearch} />
      <BundleComponentList rows={rows} options={options} onAdd={onAdd} onChange={onChange} onRemove={onRemove} />
    </section>
  );
}

export default BundleComponentsEditor;
