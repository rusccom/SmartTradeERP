import { useRef, useState } from "react";
import { Link } from "react-router-dom";
import { ChevronLeft, Package } from "lucide-react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { useCurrencies } from "../../currencies/model/useCurrencies";
import { orderPath } from "../menu";
import useOrderEditor from "./useOrderEditor";
import OrderLines from "./OrderLines";
import OrderDetails from "./OrderDetails";
import OrderAddresses from "./OrderAddresses";
import { displayMoney } from "./entryModel";

export default function OrderEntryForm({ initial, onPersist, initialError }) {
  const { t, locale } = useI18n();
  const { defaultCurrency } = useCurrencies();
  const formRef = useRef(null);
  const [notice, setNotice] = useState(initialError || "");
  const editor = useOrderEditor(initial, (c, err) => { setNotice(err || t("entry.saved")); onPersist(c, err); });
  const code = editor.card?.currency_code || editor.quote.data?.currency_code || defaultCurrency?.code;
  const save = (format) => { if (formRef.current?.reportValidity()) { setNotice(""); editor.save(format); } };
  const money = (v) => displayMoney(v, code, locale);
  return <form ref={formRef} className="order-entry" onSubmit={(e) => { e.preventDefault(); save(); }}>
    <Link to={orderPath} className="order-back"><ChevronLeft size={16}/>{t("orders.list.title")}</Link>
    <header className="order-section-head"><h1>{editor.card?.number || t("entry.new")}</h1>
      {editor.card && <span className="order-paid">{t("entry.paidOf", { paid: money(editor.card.paid_amount), total: money(editor.quote.data?.total_with_delivery) })}</span>}
      {!editor.card?.invoice && <button type="submit" className="order-btn order-primary" disabled={editor.busy}>{t(editor.busy ? "entry.saving" : "entry.save")}</button>}</header>
    {(editor.error || notice) && <p className={editor.error || notice === initialError ? "order-error" : "order-notice"} role="status">{editor.error || notice}</p>}
    <OrderLines editor={editor} code={code} />{totals(editor, money, t)}
    <OrderDetails editor={editor} onIssue={save} /><OrderAddresses editor={editor} />
    <section className="order-card order-shipments"><Package size={19}/><div><h2>{t("entry.shipments")}</h2><p>{t("entry.inpost")}</p></div></section>
  </form>;
}

function totals(editor, money, t) {
  const totals = editor.quote.data;
  return <dl className="order-totals">
    <div><dt>{t("entry.goods")}</dt><dd>{money(totals?.total_amount)}</dd></div>
    <div><dt>{t("entry.deliveryCost")}</dt><dd>{money(editor.form.delivery_cost)}</dd></div>
    <div><dt>{t("entry.total")}</dt><dd>{money(totals?.total_with_delivery)}</dd></div>
    <div><dt>{t("entry.qty")}</dt><dd>{totals?.quantity ?? "—"}</dd></div>
  </dl>;
}
