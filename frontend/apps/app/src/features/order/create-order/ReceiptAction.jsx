import { useRef, useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import Modal from "@smarterp/ui/modal/Modal";
import OrderField from "./OrderField";

export default function ReceiptAction({ editor }) {
  const { t } = useI18n();
  const [open, setOpen] = useState(false);
  const [amount, setAmount] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
  const request = useRef(null);
  const start = () => { request.current ??= crypto.randomUUID(); setOpen(true); setError(""); };
  const change = (v) => { setAmount(v); request.current = crypto.randomUUID(); };
  const finish = () => { setOpen(false); setAmount(""); request.current = null; };
  const save = () => record({ editor, amount, id: request.current, setPending, setError, finish });
  return <div className="order-receipt"><button type="button" className="order-link" disabled={editor.busy || editor.dirty} onClick={start}>{t("entry.receipt")}</button>
    {open && <Modal open onClose={() => { if (!pending) setOpen(false); }} closeLabel={t("entry.close")} title={t("entry.receipt")}>
      <OrderField label={`${t("entry.receivedAmount")} (${editor.card.currency_code})`} type="number" min="0" step="any" value={amount} disabled={pending} onChange={change} />
      {error && <p role="alert" className="order-error">{error}</p>}
      <button className="order-btn order-primary" type="button" disabled={pending || !(Number(amount) > 0)} onClick={save}>{t("entry.record")}</button></Modal>}
  </div>;
}

async function record(s) {
  s.setPending(true);
  s.setError("");
  try {
    const card = await postJSON(apiPaths.orderPayment(s.editor.card.id), { request_id: s.id, amount: s.amount });
    s.editor.setCard(card);
    s.finish();
  } catch (e) { s.setError(e.message); }
  finally { s.setPending(false); }
}
