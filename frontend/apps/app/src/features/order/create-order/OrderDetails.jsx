import { useI18n } from "@smarterp/i18n/useI18n";
import BoundField from "./BoundField";
import InvoiceActions from "./InvoiceActions";
import ReceiptAction from "./ReceiptAction";

export default function OrderDetails({ editor, onIssue }) {
  const { t } = useI18n();
  const field = (path, key, props = {}) => <BoundField key={path} editor={editor} path={path} label={t(`entry.${key}`)} {...props} />;
  return <div className="order-grid">
    <section className="order-card"><h2>{t("entry.buyer")}</h2>{field("details.nickname", "nickname")}
      {field("customer_name", "customerName", { required: true, maxLength: 160 })}{field("customer_phone", "phone", { type: "tel" })}
      {field("customer_email", "email", { type: "email" })}</section>
    <section className="order-card"><h2>{t("entry.delivery")}</h2>
      {field("details.delivery_kind", "deliveryKind", { options: [["address", t("entry.door")], ["pickup", t("entry.pickup")]] })}
      {field("details.delivery_method", "method")}{field("delivery_cost", "deliveryCost", { type: "number", step: "0.0001", min: "0" })}
      {field("details.delivery_vat", "deliveryVat", { type: "number", step: "0.01", min: "0", max: "100" })}</section>
    <section className="order-card"><h2>{t("entry.payment")}</h2>{field("details.payment_method", "method")}
      {checkbox(editor, "details.cod", editor.form.details.cod, t("entry.cod"))}
      {checkbox(editor, "invoice_requested", editor.form.invoice_requested, t("entry.invoiceRequested"))}
      {editor.form.invoice_requested && <InvoiceActions editor={editor} onIssue={onIssue} />}
      {editor.card && <ReceiptAction editor={editor} />}</section>
  </div>;
}

function checkbox(editor, path, checked, label) {
  return <label className="order-check"><input type="checkbox" checked={checked} disabled={editor.busy || Boolean(editor.card?.invoice)}
    onChange={(e) => editor.change(path, e.target.checked)} />{label}</label>;
}
