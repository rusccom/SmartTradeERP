import { RefreshCw } from "lucide-react";

function BundleList({ bundles, loading, onReload, onSelect, selectedID }) {
  return (
    <section className="bundles-list">
      <header className="bundles-section-head">
        <h2>Bundles</h2>
        <button className="bundles-icon-btn" type="button" onClick={onReload} disabled={loading} title="Reload bundles">
          <RefreshCw size={16} />
        </button>
      </header>
      <div className="bundles-table-wrap">
        <table className="bundles-table">
          <thead>
            <tr>
              <th>Product</th>
              <th>Variant</th>
              <th>SKU</th>
              <th>Price</th>
            </tr>
          </thead>
          <tbody>{renderRows({ bundles, loading, onSelect, selectedID })}</tbody>
        </table>
      </div>
    </section>
  );
}

function renderRows(props) {
  if (props.loading) return renderEmptyRow("Loading bundles...");
  if (!props.bundles.length) return renderEmptyRow("No composite products yet.");
  return props.bundles.map((bundle) => renderBundleRow(bundle, props));
}

function renderBundleRow(bundle, props) {
  const active = bundle.variant_id === props.selectedID ? "bundles-row bundles-row--active" : "bundles-row";
  return (
    <tr key={bundle.variant_id} className={active} onClick={() => props.onSelect(bundle.variant_id)}>
      <td>{bundle.product_name}</td>
      <td>{bundle.variant_name}</td>
      <td>{bundle.sku_code || "-"}</td>
      <td>{bundle.price}</td>
    </tr>
  );
}

function renderEmptyRow(text) {
  return (
    <tr>
      <td className="bundles-empty" colSpan="4">{text}</td>
    </tr>
  );
}

export default BundleList;
