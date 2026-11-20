import { useState } from "react";
import { Link } from "react-router-dom";
import { FileText } from "lucide-react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { apiPaths } from "@smarterp/api/publicApi";
import { invoiceRegions } from "../../regions/orderInvoices";
import OrderField from "./OrderField";
import useOrderResource from "./useOrderResource";

export default function InvoiceActions({ editor, onIssue }) {
  const { t } = useI18n();
  const invoice = editor.card?.invoice;
  const formats = useOrderResource(invoice ? "" : apiPaths.invoiceFormats);
  const [selection, setSelection] = useState("");
  const region = invoiceRegions[invoice?.format || selection];
  const Delivery = region?.Delivery;
  if (invoice) return <div className="order-invoice"><span className="order-invoice-number"><FileText size={16}/>{invoice.number}</span>
    {Delivery && <Delivery invoice={invoice} />}</div>;
  const options = [["", t("entry.chooseFormat")], ...(formats.data ?? []).map((f) => [f.id, f.label])];
  return <div className="order-issue"><OrderField label={t("entry.invoiceFormat")} value={selection} onChange={setSelection} options={options} disabled={editor.busy} />
    <button type="button" className="order-btn order-primary" disabled={editor.busy || !selection} onClick={() => onIssue(selection)}>{t("entry.create")}</button>
    {formats.error && <button type="button" className="order-link" onClick={formats.retry}>{formats.error} · {t("entry.retry")}</button>}
    {region && <Link className="order-link" to={region.settings} target="_blank" rel="noreferrer">{t(region.label)}</Link>}
  </div>;
}
