import { useI18n } from "@smarterp/i18n/useI18n";
import { useTenantPage } from "../api/useTenantPage";

export default function TenantPicker({ selection, select, changePage }) {
  const { t } = useI18n();
  const { items, error, loading, total } = useTenantPage(selection.page);
  if (error) return <p role="alert">{error}</p>;
  if (loading) return <p role="status">{t("admin.organizations.loading")}</p>;
  if (!items.length) return <p>{t("admin.organizations.empty")}</p>;
  return (
    <div className="plugins-panel">
      <label htmlFor="plugin-organization">{t("admin.organizations.label")}</label>
      <select id="plugin-organization" value={selection.id} onChange={event => select(event.target.value)}>
        <option value="">{t("admin.organizations.select")}</option>
        {items.map(item => <option value={item.id} key={item.id}>{item.name}</option>)}
      </select>
      {selection.page > 1 && <button type="button" onClick={() => changePage(-1)}>{t("dataTable.previousPage")}</button>}
      {selection.page * 50 < total && <button type="button" onClick={() => changePage(1)}>{t("dataTable.nextPage")}</button>}
    </div>
  );
}
