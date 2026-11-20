import { useMemo, useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import { ServerListTable } from "@smarterp/ui/tables/list-table";
import Modal from "@smarterp/ui/modal/Modal";
import { createCustomersTablePreset } from "../../customers/customersTablePreset";
import OrderField from "./OrderField";

export default function CustomerPicker({ create, onPick, onClose }) {
  const { t } = useI18n();
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const preset = useMemo(() => pickerPreset(t), [t]);
  async function add() {
    setBusy(true);
    try { onPick(await postJSON(apiPaths.customers, { name: name.trim() })); }
    catch (e) { setError(e.message); setBusy(false); }
  }
  return <Modal open onClose={onClose} closeOnEscape closeLabel={t("entry.close")} title={t(create ? "entry.addCustomer" : "entry.chooseCustomer")}>
    {create ? <div className="order-card"><OrderField label={t("entry.customerName")} value={name} onChange={setName} maxLength={160} />
      <button className="order-btn" type="button" disabled={busy || !name.trim()} onClick={add}>{t("entry.addCustomer")}</button></div>
      : <ServerListTable preset={preset} onRowOpen={onPick} />}
    {error && <p role="alert" className="order-error">{error}</p>}
  </Modal>;
}

function pickerPreset(t) {
  const preset = createCustomersTablePreset(t);
  return { ...preset, columns: preset.columns.map((c) => c.accessorKey === "name" ? { ...c, openOnClick: true } : c) };
}
