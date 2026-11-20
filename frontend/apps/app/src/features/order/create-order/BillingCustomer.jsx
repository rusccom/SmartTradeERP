import { useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { apiPaths } from "@smarterp/api/publicApi";
import CustomerPicker from "./CustomerPicker";
import useApiResource from "@smarterp/api/useApiResource";

export default function BillingCustomer({ editor }) {
  const { t } = useI18n();
  const [mode, setMode] = useState("");
  const id = editor.form.billing_customer_id;
  const customer = useApiResource(id ? apiPaths.customerById(id) : "");
  const locked = Boolean(editor.card?.invoice) || editor.busy;
  function pick(row) {
    editor.change("billing_customer_id", row.id);
    if (!editor.form.details.billing.name && !editor.form.details.billing.company) editor.change("details.billing.name", row.name);
    setMode("");
  }
  return <div className="order-field"><span>{t("entry.customer")}</span><div className="order-actions">
    <button type="button" className="order-btn" disabled={locked} onClick={() => setMode("pick")}>{customer.data?.name || (id ? id : t("entry.chooseCustomer"))}</button>
    {!locked && <button type="button" className="order-btn" aria-label={t("entry.addCustomer")} onClick={() => setMode("create")}>+</button>}
    {id && !locked && <button type="button" className="order-link" onClick={() => editor.change("billing_customer_id", "")}>{t("entry.clear")}</button>}
  </div>{customer.error && <span role="alert">{customer.error}</span>}
    {mode && <CustomerPicker create={mode === "create"} onPick={pick} onClose={() => setMode("")} />}</div>;
}
