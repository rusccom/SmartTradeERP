import { useRef } from "react";
import { useStatusDraft } from "../model/useStatusDraft";
import StatusDraftForm from "../ui/StatusDraftForm";
import StatusDraftTable from "../ui/StatusDraftTable";

export default function OrderSettingsPage() {
  const draft = useStatusDraft();
  const input = useRef(null);
  function startEditing(status) {
    draft.edit(status);
    input.current?.focus();
  }
  return (
    <section className="order-settings" aria-labelledby="order-settings-title">
      <h1 id="order-settings-title">Налаштування замовлень</h1>
      <div className="order-status-panel">
        <h2>Статуси замовлень</h2>
        <p>Назви залишаються в чернетці до виходу з розділу або оновлення сторінки.</p>
        <StatusDraftForm draft={draft} inputRef={input} />
        <StatusDraftTable statuses={draft.statuses} onEdit={startEditing} />
      </div>
    </section>
  );
}
