import { Plus, Save } from "lucide-react";

import BundleComponentRow from "./BundleComponentRow";
import ComponentOptionsToolbar from "./ComponentOptionsToolbar";

function BundleComponentsEditor(props) {
  const { bundle, canLoadMoreOptions, componentSearch, error, loading, onAdd, onChange, onLoadMoreOptions,
    onRemove, onSave, onSearchComponents, options, optionsLoading, rows, saving } = props;
  if (!bundle) return <section className="bundles-editor bundles-editor--empty">Select a bundle to edit components.</section>;
  return (
    <section className="bundles-editor">
      <header className="bundles-section-head">
        <div>
          <h2>{bundle.product_name}</h2>
          <p>{bundle.variant_name} · {bundle.unit}</p>
        </div>
        <button className="bundles-save-btn" type="button" onClick={onSave} disabled={saving || loading}>
          <Save size={16} /> {saving ? "Saving" : "Save"}
        </button>
      </header>
      {error && <p className="bundles-error">{error}</p>}
      <ComponentOptionsToolbar
        canLoadMore={canLoadMoreOptions}
        loading={optionsLoading}
        onLoadMore={onLoadMoreOptions}
        onSearch={onSearchComponents}
        search={componentSearch}
      />
      <div className="bundle-components-head">
        <span>Component</span>
        <span>Qty per unit</span>
        <span />
      </div>
      <div className="bundle-components-list">
        {rows.map((row) => <BundleComponentRow key={row.id} row={row} options={options} onChange={onChange} onRemove={onRemove} canRemove={rows.length > 1} />)}
      </div>
      <button className="bundles-add-btn" type="button" onClick={onAdd} disabled={!options.length}>
        <Plus size={16} /> Add component
      </button>
    </section>
  );
}

export default BundleComponentsEditor;
