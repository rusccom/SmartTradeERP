import { Plus, RefreshCw } from "lucide-react";

import { useI18n } from "../../../shared/i18n/useI18n";

function BundleList({ bundles, loading, onCreate, onReload, onSelect, selectedID }) {
  const { t } = useI18n();
  const showHeader = bundles.length > 0;
  return (
    <section className="bundles-list">
      <header className="bundles-section-head">
        <h2>{t("bundles.table.title")}</h2>
        <BundleListActions loading={loading} onCreate={onCreate} onReload={onReload} t={t} />
      </header>
      <div className="bundles-table-wrap">
        <BundleTable
          bundles={bundles}
          loading={loading}
          onSelect={onSelect}
          selectedID={selectedID}
          showHeader={showHeader}
          t={t}
        />
      </div>
    </section>
  );
}

function BundleTable(props) {
  return (
    <table className="bundles-table">
      {props.showHeader && <BundleTableHead t={props.t} />}
      <tbody>{renderRows(props)}</tbody>
    </table>
  );
}

function BundleTableHead({ t }) {
  return (
    <thead>
      <tr>
        <th>{t("bundles.columns.product")}</th>
        <th>{t("bundles.columns.variant")}</th>
        <th>{t("bundles.columns.sku")}</th>
        <th>{t("bundles.columns.price")}</th>
      </tr>
    </thead>
  );
}

function BundleListActions({ loading, onCreate, onReload, t }) {
  return (
    <div className="bundles-head-actions">
      <button className="bundles-create-btn" type="button" onClick={onCreate}>
        <Plus size={16} /> {t("bundles.addButton")}
      </button>
      <button className="bundles-icon-btn" type="button" onClick={onReload} disabled={loading} title={t("bundles.table.reload")}>
        <RefreshCw size={16} />
      </button>
    </div>
  );
}

function renderRows(props) {
  if (props.loading) return renderEmptyRow(props.t("bundles.table.loading"));
  if (!props.bundles.length) return renderEmptyRow(props.t("bundles.table.empty"));
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
